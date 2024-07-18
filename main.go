// This main.go acts as a wrapper/helper for cmd/camagent
// It is almost identical to it; only that it (optionally) helps to read a file with appIDs (if the file exists), and
// creates a process for each appID, and passes the flags along
package main

import (
	"bufio"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"main/helpers"
)

var camLogger = helpers.CreateLogger("[CAM-HELPER]")

func main() {
	var idPath string
	flag.StringVar(&idPath, "idPath", "", "--idPath=/path/to/an/idfile.ext, for more information, refer to the README in cmd/idgen")
	configFlags := helpers.ReadFlags()

	camLogger.Debug("Runtime and parsed flags information", "OS", runtime.GOOS, "Architecture", runtime.GOARCH, "clearAchievements", configFlags.ClearAchievements, "zeroShouldFail", configFlags.ZeroShouldFail, "appID", configFlags.AppID, "libraryPath", configFlags.LibraryPath)

	// Only tested on linux (+ not tested anything x32), warn that it may not work
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		camLogger.Warnf("CAM was not tested with %s on the %s architecture, and it may not work", runtime.GOOS, runtime.GOARCH)
	}

	if configFlags.LibraryPath == "" {
		camLogger.Warn("libraryPath not set. Trying to guess based on your OS")
		if configFlags.LibraryPath = helpers.TryToGuessLibraryPath(); configFlags.LibraryPath == "" {
			camLogger.Fatal("Unsupported OS or Architecture (Wasn't able to guess it), exiting")
		}
	}

	var commandArgs []string
	if configFlags.ClearAchievements {
		commandArgs = append(commandArgs, "--clearAchievements")
	}

	if configFlags.ZeroShouldFail {
		commandArgs = append(commandArgs, "--zeroShouldFail")
	}

	commandArgs = append(commandArgs, "--libraryPath", configFlags.LibraryPath)

	if configFlags.AppID != "" {
		if agentErr := startAgent(append(commandArgs, "--appID", configFlags.AppID)); agentErr != nil {
			camLogger.Error("Agent encountered an error", "error", agentErr)
		}

		// We are done, we can exit
		os.Exit(0)
	}

	if idPath == "" {
		camLogger.Fatal("No appID or idPath was specified")
	}

	idFilePath := filepath.Clean(idPath)
	idFile, idFileErr := os.OpenFile(idFilePath, os.O_RDONLY, helpers.FilePermissions)
	if idFileErr != nil {
		camLogger.Fatal("Failed to read ID file", "error", idFileErr)
	}

	defer func() {
		if closeErr := idFile.Close(); closeErr != nil {
			camLogger.Fatal("Failed to close ID file", "error", closeErr)
		}
	}()

	camLogger.Debugf("ID file at path %s was found, now reading it", idFilePath)
	fileScanner := bufio.NewScanner(idFile)

	for fileScanner.Scan() {
		fappID := fileScanner.Text()

		if agentErr := startAgent(append(commandArgs, "--appID", fappID)); agentErr != nil {
			camLogger.Error("Agent encountered an error", "error", agentErr)
		}
	}

	if fileScannerErr := fileScanner.Err(); fileScannerErr != nil {
		camLogger.Fatal("Failed to read ID file", "error", fileScannerErr)
	}
}

func startAgent(commandArgs []string) error {
	camLogger.Info("Starting agent", "commandArgs", commandArgs)

	cmd := exec.Command("./camagent", commandArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
