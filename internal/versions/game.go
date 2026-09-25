/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package versions

var gameVersion = -1

func GetGameVersion() int {
	return gameVersion
}

func SetGameVersion(v int) {
	gameVersion = v
}
