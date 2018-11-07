package game

func init() {
	game := GetGame()
	go game.Run()
}
