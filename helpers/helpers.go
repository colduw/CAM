package helpers

import (
	"flag"
	"os"
	"runtime"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type configFlags struct {
	AppID             string
	LibraryPath       string
	ClearAchievements bool
	ZeroShouldFail    bool
}

const (
	// File permissions (owner reads and writes)
	FilePermissions = 0o600
	// Steam is pretty wonky sometimes, and can fail to either initialize, request the user stats,
	// or request the amount of achievements a game has. SleepDelay is used to  sleep in a
	// for loop, until the aforementioned operation(s) are completed successfully.
	SleepDelay = 300 * time.Millisecond
)

// Assumes that the libraries are in the current working directories "lib" folder
// This is only ran if --libraryPath is not set
func TryToGuessLibraryPath() string {
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "386":
			return "./lib/steam_api32.dll"
		case "amd64":
			return "./lib/steam_api64.dll"
		}
	case "linux", "freebsd":
		switch runtime.GOARCH {
		case "386":
			return "./lib/libsteam_api32.so"
		case "amd64":
			return "./lib/libsteam_api64.so"
		}
	case "darwin":
		return "./lib/libsteam_api.dylib"
	}

	return ""
}

func ReadFlags() configFlags {
	var cflags configFlags

	flag.BoolVar(&cflags.ClearAchievements, "clearAchievements", false, "--clearAchievements to reset (clear) achievements")
	flag.BoolVar(&cflags.ZeroShouldFail, "zeroShouldFail", false, "--zeroShouldFail to treat zero reported achievements as a fail and retry")
	flag.StringVar(&cflags.AppID, "appID", "", "--appID=... to set the appid to unlock/clear the achievements for")
	flag.StringVar(&cflags.LibraryPath, "libraryPath", "", "--libraryPath=/path/to/lib/steam_api to use a specific path pointing at the steam library. If not set, we'll try to automatically guess it (assuming the library is in the ./lib folder)")
	flag.Parse()

	return cflags
}

func CreateLogger(name string) *log.Logger {
	thisLogger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		Level:           log.DebugLevel,
		TimeFormat:      time.DateTime,
		Prefix:          name,
	})

	newStyle := log.DefaultStyles()
	newColors := map[log.Level]lipgloss.Style{
		log.ErrorLevel: lipgloss.NewStyle().
			SetString("ERR").
			Bold(true).
			Width(3).
			MaxWidth(3).
			Foreground(lipgloss.Color("1")),
		log.InfoLevel: lipgloss.NewStyle().
			SetString("INF").
			Bold(true).
			Width(3).
			MaxWidth(3).
			Foreground(lipgloss.Color("2")),
		log.WarnLevel: lipgloss.NewStyle().
			SetString("WRN").
			Bold(true).
			Width(3).
			MaxWidth(3).
			Foreground(lipgloss.Color("3")),
		log.DebugLevel: lipgloss.NewStyle().
			SetString("DBG").
			Bold(true).
			Width(3).
			MaxWidth(3).
			Foreground(lipgloss.Color("4")),
		log.FatalLevel: lipgloss.NewStyle().
			SetString("FTL").
			Bold(true).
			Width(3).
			MaxWidth(3).
			Foreground(lipgloss.Color("5")),
	}

	newStyle.Levels = newColors

	thisLogger.SetStyles(newStyle)

	return thisLogger
}
