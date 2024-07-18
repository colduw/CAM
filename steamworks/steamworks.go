package steamworks

import (
	"github.com/ebitengine/purego"
)

type (
	// const int k_cchMaxSteamErrMsg = 1024;
	// typedef char SteamErrMsg[ k_cchMaxSteamErrMsg ];
	SteamErrMsg [1024]byte
)

const (
	K_ESteamAPIInitResult_OK = 0
)

var (
	SteamAPI_Init       func(pOutErrMsg *SteamErrMsg) int
	SteamUserStats      func() uintptr
	RequestCurrentStats func(ISteamUserStats uintptr) bool
	GetNumAchievements  func(ISteamUserStats uintptr) uint32
	GetAchievementName  func(ISteamUserStats uintptr, iAchievement uint32) string
	GetAchievement      func(ISteamUserStats uintptr, pchName string, pbAchieved *bool) bool
	ClearAchievement    func(ISteamUserStats uintptr, pchName string) bool
	SetAchievement      func(ISteamUserStats uintptr, pchName string) bool
	StoreStats          func(ISteamUserStats uintptr) bool
	SteamAPI_Shutdown   func()

	SteamUser  func() uintptr
	GetSteamID func(ISteamUser uintptr) uint64
)

func RegisterLibraryFuncs(libraryPath string) {
	steamLib, openErr := loadSteamLibrary(libraryPath)
	if openErr != nil {
		// If we can't open the library, there is no point in continuing as RegisterLibFunc will panic either way.
		panic(openErr)
	}

	// Exported functions from the Steamworks SDK
	purego.RegisterLibFunc(&SteamAPI_Init, steamLib, "SteamAPI_InitFlat")
	purego.RegisterLibFunc(&SteamUserStats, steamLib, "SteamAPI_SteamUserStats_v012")
	purego.RegisterLibFunc(&RequestCurrentStats, steamLib, "SteamAPI_ISteamUserStats_RequestCurrentStats")
	purego.RegisterLibFunc(&GetNumAchievements, steamLib, "SteamAPI_ISteamUserStats_GetNumAchievements")
	purego.RegisterLibFunc(&GetAchievementName, steamLib, "SteamAPI_ISteamUserStats_GetAchievementName")
	purego.RegisterLibFunc(&GetAchievement, steamLib, "SteamAPI_ISteamUserStats_GetAchievement")
	purego.RegisterLibFunc(&ClearAchievement, steamLib, "SteamAPI_ISteamUserStats_ClearAchievement")
	purego.RegisterLibFunc(&SetAchievement, steamLib, "SteamAPI_ISteamUserStats_SetAchievement")
	purego.RegisterLibFunc(&StoreStats, steamLib, "SteamAPI_ISteamUserStats_StoreStats")
	purego.RegisterLibFunc(&SteamAPI_Shutdown, steamLib, "SteamAPI_Shutdown")

	purego.RegisterLibFunc(&SteamUser, steamLib, "SteamAPI_SteamUser_v023")
	purego.RegisterLibFunc(&GetSteamID, steamLib, "SteamAPI_ISteamUser_GetSteamID")
}
