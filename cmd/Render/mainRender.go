package main

import (
	tictac "../.."
	"os"
)

func main() {
	f, err := os.OpenFile("out.svg", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	const S int = 16000

	tictac.RenderSVG(f, S)
}
