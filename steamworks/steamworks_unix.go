//go:build darwin || freebsd || linux

package steamworks

import (
	"github.com/ebitengine/purego"
)

func loadSteamLibrary(libraryPath string) (uintptr, error) {
	return purego.Dlopen(libraryPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}
