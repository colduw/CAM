package main

import (
	"bytes"
	"os"
	"time"

	"main/helpers"
	"main/steamworks"
)

var (
	camLogger = helpers.CreateLogger("[CAM-AGENT]")
)

const (
	maxRetryAttempts = 128
)

func main() {
	configFlags := helpers.ReadFlags()

	if configFlags.LibraryPath == "" {
		camLogger.Warn("libraryPath not set. Trying to guess based on your OS")

		if configFlags.LibraryPath = helpers.TryToGuessLibraryPath(); configFlags.LibraryPath == "" {
			camLogger.Fatal("Unsupported OS or Architecture (Wasn't able to guess it), exiting.")
		}
	}

	if configFlags.AppID == "" {
		camLogger.Fatal("No appID was specified")
	}

	steamworks.RegisterLibraryFuncs(configFlags.LibraryPath)

	if writeErr := os.WriteFile("steam_appid.txt", []byte(configFlags.AppID), helpers.FilePermissions); writeErr != nil {
		camLogger.Fatal("Failed to write steam_appid.txt", "error", writeErr)
	}

	var errMsg steamworks.SteamErrMsg
	tresholdLimiter(func() bool {
		if steamworks.SteamAPI_Init(&errMsg) == steamworks.K_ESteamAPIInitResult_OK {
			return true
		}

		camLogger.Error("Failed to initialize Steam", "error", string(bytes.Trim(errMsg[:], "\x00")))
		return false
	})

	userStats := steamworks.SteamUserStats()
	tresholdLimiter(func() bool {
		if userStats = steamworks.SteamUserStats(); userStats != 0 {
			return true
		}

		camLogger.Error("SteamUserStats(): failed")
		return false
	})

	tresholdLimiter(func() bool {
		if steamworks.RequestCurrentStats(userStats) {
			return true
		}

		camLogger.Error("RequestCurrentStats(): failed")
		return false
	})

	numOfAchievements := steamworks.GetNumAchievements(userStats)
	if numOfAchievements == 0 {
		if !configFlags.ZeroShouldFail {
			os.Exit(0)
		}

		tresholdLimiter(func() bool {
			if numOfAchievements = steamworks.GetNumAchievements(userStats); numOfAchievements != 0 {
				return true
			}

			camLogger.Warnf("GetNumAchievements(): reported 0 achievements for appID %s", configFlags.AppID)
			return false
		})
	}

	camLogger.Infof("GetNumAchievements(): reported %d achievements for appID %s", numOfAchievements, configFlags.AppID)

	for i := uint32(0); i < numOfAchievements; i++ {
		achName := steamworks.GetAchievementName(userStats, i)
		var achieved bool

		if !steamworks.GetAchievement(userStats, achName, &achieved) {
			camLogger.Error("GetAchievement(): failed")
			continue
		}

		if configFlags.ClearAchievements {
			camLogger.Debugf("ClearAchievement(): %d->%s: Clearing achievement.....", i+1, achName)
			if !achieved {
				camLogger.Warn("Not achieved")
				continue
			}

			if !steamworks.ClearAchievement(userStats, achName) {
				camLogger.Error("Failed")
				continue
			}

			camLogger.Info("OK")
		} else {
			camLogger.Debugf("SetAchievement(): %d->%s: Setting achievement.....", i+1, achName)
			if achieved {
				camLogger.Warn("Already achieved")
				continue
			}

			if !steamworks.SetAchievement(userStats, achName) {
				camLogger.Error("Failed")
				continue
			}

			camLogger.Info("OK")
		}
	}

	camLogger.Debug("StoreStats(): Storing stats.....")
	if !steamworks.StoreStats(userStats) {
		camLogger.Error("Failed")
	} else {
		camLogger.Info("OK")
	}

	camLogger.Debug("SteamAPI_Shutdown(): Finished, shutting down")
	steamworks.SteamAPI_Shutdown()
}

func tresholdLimiter(callee func() bool) {
	treshold := 0

	for !callee() {
		if treshold > maxRetryAttempts {
			os.Exit(1)
		}

		time.Sleep(helpers.SleepDelay)
		treshold++
	}
}
