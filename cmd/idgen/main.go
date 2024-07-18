package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"main/helpers"
	"main/steamworks"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GetOwnedGames struct {
	Response struct {
		Games []struct {
			AppID int64 `json:"appid"`
		} `json:"games"`
	} `json:"response"`
}

type GetPlayerAchievements struct {
	PlayerStats struct {
		Achievements []struct {
			Achieved int `json:"achieved"`
		} `json:"achievements"`
	} `json:"playerstats"`
}

var (
	genLogger     = helpers.CreateLogger("[ID-GEN]")
	timeoutClient = &http.Client{Timeout: time.Minute}
	gMap          = make(map[int64]struct{})
	gMutex        sync.Mutex
)

const (
	concurrentRequests = 4
)

func main() {
	var webapi, steamid string
	flag.StringVar(&webapi, "key", "", "--key=STEAM_WEBAPI_KEY")
	flag.StringVar(&steamid, "steamid", "", "--steamid=YOUR_STEAM_ID")
	configFlags := helpers.ReadFlags()

	if webapi == "" {
		genLogger.Fatal("No API Key provided")
	}

	if steamid == "" {
		if configFlags.LibraryPath == "" {
			genLogger.Warn("libraryPath not set. Trying to guess based on your OS")

			if configFlags.LibraryPath = helpers.TryToGuessLibraryPath(); configFlags.LibraryPath == "" {
				genLogger.Fatal("Unsupported OS or Architecture (Wasn't able to guess it), exiting.")
			}
		}

		steamworks.RegisterLibraryFuncs(configFlags.LibraryPath)

		var errMsg steamworks.SteamErrMsg
		for steamworks.SteamAPI_Init(&errMsg) != steamworks.K_ESteamAPIInitResult_OK {
			genLogger.Error("Failed to initialize Steam", "error", string(bytes.Trim(errMsg[:], "\x00")))
			time.Sleep(helpers.SleepDelay)
		}

		steamid = strconv.FormatUint(steamworks.GetSteamID(steamworks.SteamUser()), 10)
		genLogger.Debugf("SteamID is %s", steamid)
		steamworks.SteamAPI_Shutdown()
	}

	idsFile, idsFileErr := os.OpenFile("ids.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, helpers.FilePermissions)
	if idsFileErr != nil {
		genLogger.Fatal("Failed to create or open ids.txt", "error", idsFileErr)
	}

	defer func() {
		if closeErr := idsFile.Close(); closeErr != nil {
			genLogger.Fatal("Failed to close IDs file", "error", closeErr)
		}
	}()

	request, requestErr := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v1?key=%s&steamid=%s&include_played_free_games=true&include_free_sub=true&skip_unvetted_apps=false", webapi, steamid), http.NoBody)
	if requestErr != nil {
		genLogger.Fatal("Failed to create new request", "error", requestErr)
	}

	response, responseErr := timeoutClient.Do(request)
	if responseErr != nil {
		genLogger.Fatal("Failed to do the request (or client timed out)", "error", responseErr)
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			genLogger.Error("Failed to close response body", "error", closeErr)
		}
	}()

	var ownedGames GetOwnedGames
	if decodeErr := json.NewDecoder(response.Body).Decode(&ownedGames); decodeErr != nil {
		genLogger.Fatal("Failed to decode JSON", "error", decodeErr)
	}

	genLogger.Debugf("Steam reported %d games in the account", len(ownedGames.Response.Games))

	sem := make(chan struct{}, concurrentRequests)
	var waitGroup sync.WaitGroup
	for i := 0; i < len(ownedGames.Response.Games); i++ {
		waitGroup.Add(1)
		sem <- struct{}{}

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			defer func() { <-sem }()

			appID := ownedGames.Response.Games[i].AppID
			requestURL := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetPlayerAchievements/v1?key=%s&steamid=%s&appid=%d", webapi, steamid, appID)

			genLogger.Infof("Sent request for appID %d", appID)
			for retryErr := retry(requestURL, appID); retryErr != nil; retryErr = retry(requestURL, appID) {
				if strings.Contains(retryErr.Error(), "invalid character") {
					genLogger.Error("Likely timed out, waiting a minute", "error", retryErr)
					time.Sleep(time.Minute)
				} else {
					genLogger.Error("Request failed due to an unknown error, waiting for a bit", "error", retryErr)
					time.Sleep(10 * time.Second)
				}
			}
		}(&waitGroup)
	}

	waitGroup.Wait()
	genLogger.Debugf("%d games have achievements left to unlock, writing them to ids.txt", len(gMap))

	for v := range gMap {
		if _, writeErr := idsFile.WriteString(strconv.FormatInt(v, 10) + "\n"); writeErr != nil {
			genLogger.Error("Failed to append appID to file", "appID", v, "error", writeErr)
			continue
		}

		genLogger.Debug("Wrote appID to file", "appID", v)
	}
}

func retry(url string, appID int64) error {
	gRequest, gRequestErr := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if gRequestErr != nil {
		return gRequestErr
	}

	gResponse, gResponseErr := timeoutClient.Do(gRequest)
	if gResponseErr != nil {
		return gResponseErr
	}
	defer func() {
		if closeErr := gResponse.Body.Close(); closeErr != nil {
			genLogger.Error("Failed to close response body", "error", closeErr)
		}
	}()

	var pachievements GetPlayerAchievements
	if decodeErr := json.NewDecoder(gResponse.Body).Decode(&pachievements); decodeErr != nil {
		return decodeErr
	}

	for i := 0; i < len(pachievements.PlayerStats.Achievements); i++ {
		if pachievements.PlayerStats.Achievements[i].Achieved == 0 {
			gMutex.Lock()
			gMap[appID] = struct{}{}
			gMutex.Unlock()
		}
	}

	return nil
}
