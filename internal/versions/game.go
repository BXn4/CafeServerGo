package versions

var gameVersion = -1

func GetGameVersion() int {
	return gameVersion
}

func SetGameVersion(v int) {
	gameVersion = v
}
