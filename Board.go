package tic_tac_toe_minimax

import (
	"bytes"
	"fmt"
)

type BoardState struct {
	board [][]Player
}

type Player int

func (p Player) String() string {
	switch p {
	default:
		return "?"
	case NoPlayer:
		return " "
	case PlayerX:
		return "X"
	case PlayerO:
		return "O"
	}
}

const (
	NoPlayer Player = iota
	PlayerX  Player = iota
	PlayerO  Player = iota
)

type Pos struct {
	X, Y int
}

func (b *BoardState) Size() int {
	return len(b.board)
}

// board must be square
func NewBoard(size int) (b BoardState) {
	b.board = make([][]Player, size)
	for i := range b.board {
		b.board[i] = make([]Player, size)
	}
	return
}

func CopyBoard(bIn *BoardState) (bOut BoardState) {
	bOut.board = make([][]Player, bIn.Size())
	for i := range bIn.board {
		bOut.board[i] = make([]Player, bIn.Size())
		copy(bOut.board[i], bIn.board[i])
	}
	return
}

func (b *BoardState) GetPos(p Pos) Player {
	return b.board[p.X][p.Y]
}

func (b *BoardState) SetPos(p Pos, value Player) {
	b.board[p.X][p.Y] = value
}

func (b *BoardState) OpenCells() (positions []Pos) {
	for i, row := range b.board {
		for j, cell := range row {
			if cell == NoPlayer {
				positions = append(positions, Pos{i, j})
			}
		}
	}
	return
}

func (b *BoardState) OccupiedCells() (positions []Pos) {
	for i, row := range b.board {
		for j, cell := range row {
			if cell != NoPlayer {
				positions = append(positions, Pos{i, j})
			}
		}
	}
	return
}

func (b *BoardState) Full() bool {
	for _, row := range b.board {
		for _, cell := range row {
			if cell == NoPlayer {
				return false
			}
		}
	}
	return true
}

func (b *BoardState) CheckWin() Player {
	// here we use that board square assumption, because we only check the two main diagonals

	if p := b.checkDiagonal1(); p != NoPlayer {
		return p
	}
	if p := b.checkDiagonal2(); p != NoPlayer {
		return p
	}

	for i := 0; i < b.Size(); i++ {

		if p := b.checkX(i); p != NoPlayer {
			return p
		}

		if p := b.checkY(i); p != NoPlayer {
			return p
		}

	}
	return NoPlayer
}

func (b *BoardState) checkX(x int) Player {
	last := b.board[x][0]
	for _, cell := range b.board[x] {
		if cell != last {
			return NoPlayer
		} else {
			last = cell
		}
	}
	return last
}

func (b *BoardState) checkY(y int) Player {
	last := b.board[0][y]
	for _, row := range b.board {
		cell := row[y]
		if cell != last {
			return NoPlayer
		} else {
			last = cell
		}
	}
	return last
}

func (b *BoardState) checkDiagonal1() Player {
	last := b.board[0][0]
	for i := range b.board {
		if b.board[i][i] != last {
			return NoPlayer
		} else {
			last = b.board[i][i]
		}
	}
	return last
}

func (b *BoardState) checkDiagonal2() Player {
	n := b.Size()
	last := b.board[n-1][0]
	for i := range b.board {
		if b.board[n-1-i][i] != last {
			return NoPlayer
		} else {
			last = b.board[n-1-i][i]
		}
	}
	return last
}

func (b *BoardState) String() string {
	var buffer bytes.Buffer
	for i, row := range b.board {

		if i != 0 {
			for _, _ = range row {
				buffer.WriteString("---")
			}
			buffer.WriteByte('\n')
		}

		for j, cell := range row {

			if j != 0 {
				buffer.WriteByte('|')
			}
			buffer.WriteString(cell.String())
			buffer.WriteString(cell.String())
		}

		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func (b *BoardState) StringWithIndexes() string {
	// convert to string with indexes on available spaces
	first := 1

	var buffer bytes.Buffer
	for i, row := range b.board {

		if i != 0 {
			for _, _ = range row {
				buffer.WriteString("---")
			}
			buffer.WriteByte('\n')
		}

		for j, cell := range row {

			if j != 0 {
				buffer.WriteByte('|')
			}
			if cell == NoPlayer {
				buffer.WriteString(fmt.Sprint(first))
				first++
				buffer.WriteByte(' ')
			} else {
				buffer.WriteString(cell.String())
				buffer.WriteString(cell.String())
			}
		}

		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func (b *BoardState) AllPositions() (out []Pos) {
	for i, row := range b.board {
		for j, _ := range row {
			out = append(out, Pos{i, j})
		}
	}
	return
}

func BoardDiff(b0, b1 BoardState) Pos {
	// find which cell is different between two subsequent states
	for _, pos := range b0.AllPositions() {
		if b0.GetPos(pos) != b1.GetPos(pos) {
			return pos
		}
	}
	// TODO
	return Pos{-1, -1}
}
