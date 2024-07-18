//go:build windows

package steamworks

import (
	"golang.org/x/sys/windows"
)

func loadSteamLibrary(libraryPath string) (uintptr, error) {
	return windows.LoadLibrary(libraryPath)
}
