package main

import (
	tictac "../.."
	"os"
)

func main() {
	//head := tictac.NewDefaultState(tictac.PlayerX)
	head := tictac.NewState(3, tictac.PlayerX)
	head.FindAllChildStates()

	f, err := os.OpenFile("out.svg", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	const S int = 8000

	tictac.RenderSVG(f, S)
}
