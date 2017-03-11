package main

import (
	tictac "../.."
	"bufio"
	"os"
	"fmt"
)

func main() {
	head := tictac.NewDefaultState(tictac.PlayerX)
	head.FindAllChildStates()

	cur := head
	cur = cur.Children[cur.BestChild]
	for {
		fmt.Println(cur.State.String())
		fmt.Println("fitness: ", cur.Fitness)
		if cur.State.CheckWin() == tictac.NoPlayer {
			// let the user move
			pos := PromptForMove(&cur.State)
			cur = cur.GetChildForMove(pos)

			if cur.State.CheckWin() != tictac.NoPlayer ||
					cur.State.Full() {
				fmt.Println(cur.State.String())
				fmt.Println("winner:", tictac.OppositePlayer(cur.NextPlayer).String())
				return
			}

			// make our move
			cur = cur.Children[cur.BestChild]
		} else {
			break
		}
	}

	fmt.Println("done")
	Pause()
}

func Pause() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func PromptForMove(b *tictac.BoardState) tictac.Pos {
	fmt.Println(b.StringWithIndexes())

	var i int
	var err error
	for {
		fmt.Println("Pick a cell")
		_, err = fmt.Scanf("%d", &i)
		if err == nil && i > 0 && i <= len(b.OpenCells()) {
			break
		} else {
			fmt.Println("bad number")
		}
	}

	return b.OpenCells()[i-1]
}
