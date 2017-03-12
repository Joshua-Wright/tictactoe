package tic_tac_toe_minimax

import (
	"fmt"
	"io"
)

type svgRender struct {
	barWidth  float64
	cellWidth float64
	buffer    float64
	f         io.Writer
}

func RenderSVG(f io.Writer, size int) {
	//headPlayerX := tictac.NewDefaultState(tictac.PlayerX)
	headPlayerX := NewState(3, PlayerX)
	headPlayerX.FindAllChildStates()

	headPlayerO := NewState(3, PlayerO)
	headPlayerO.FindAllChildStates()

	barWidth := 0.01
	buffer := 0.98

	render := svgRender{
		barWidth:  0.01,
		cellWidth: (1.0 - 2.0*barWidth) / 3,
		buffer:    buffer,
		f:         f,
	}

	fmt.Fprintf(f, `<svg viewBox="0 0 1 1" width="%v" height="%v">`, size, size)
	defer fmt.Fprint(f, `</svg>`)

	fmt.Fprint(f, `<g stroke-width="0.05" stroke-linecap="round" stroke="rgb(0,0,0)">`)
	defer fmt.Fprint(f, `</g>`)

	// rectangle to be background
	fmt.Fprint(f, `<rect width="1" height="1" style="fill:rgb(255,255,255);stroke:rgb(255,255,255)"/>`)

	render.renderState(&headPlayerO)
	//render.renderState(&headPlayerX)

}

func (r *svgRender) beginTranslation(pos Pos) {

	dx := float64(pos.X) * (r.cellWidth + r.barWidth)
	dy := float64(pos.Y) * (r.cellWidth + r.barWidth)
	scale := r.cellWidth
	fmt.Fprintf(r.f, `<g transform="translate(%v %v) scale(%v) translate(%v %v) scale(%v)">`,
		dx, dy, scale, (1-r.buffer)/2, (1-r.buffer)/2, r.buffer)
}

func (r *svgRender) endTranslation() {
	fmt.Fprint(r.f, `</g>`)
}

func (r *svgRender) renderState(s *StateTreeNode) {
	if p := s.Board.CheckWin(); p != NoPlayer && p != s.OurPlayer {
		return
	} else if p == s.OurPlayer {
		r.drawWins(s.Board)
		return
	}
	if (s.Board.Full()) {
		return
	}

	r.drawBars(&s.Board)

	// draw squares that are not occupied
	for _, pos := range s.Board.OccupiedCells() {
		r.beginTranslation(pos)
		// print an X
		if (s.Board.GetPos(pos) == PlayerX) {
			r.drawX()
		} else {
			r.drawO()
		}
		r.endTranslation()
	}

	if (s.OurPlayer == s.NextPlayer) {
		// draw our turn, then go again
		child := s.Children[s.BestChild]
		pos := BoardDiff(s.Board, child.Board)

		// print our X move in red
		r.beginTranslation(pos)
		if (s.OurPlayer == PlayerX) {
			r.drawRedX()
		} else {
			r.drawRedO()
		}
		r.endTranslation()

		s = &child
	}

	if p := s.Board.CheckWin(); p != NoPlayer {
		r.drawWins(s.Board)
		return
	}

	// draw each child inside the cell it represents
	for _, child := range s.Children {
		diff := BoardDiff(s.Board, child.Board)
		r.beginTranslation(diff)
		r.renderState(&child)
		r.endTranslation()
	}
}

func (r *svgRender) drawBars(b *BoardState) {
	/*
	format:
		X1	X2
		|	|
	Y1--|---|---
		|	|
	Y2--|---|---
		|	|
	*/
	fmt.Fprint(r.f, `<g>`)
	defer fmt.Fprint(r.f, `</g>`)

	for row := 1; row < b.Size(); row++ {
		row := float64(row)
		y := row*r.cellWidth + r.barWidth/2
		if row > 1 {
			y += r.barWidth * (row - 1)
		}
		fmt.Fprintf(r.f, `<line x1="0" x2="1" y1="%v" y2="%v" style="stroke:rgb(0,0,0);stroke-width:%v" />`, y, y, r.barWidth)
		fmt.Fprintf(r.f, `<line x1="%v" x2="%v" y1="0" y2="1" style="stroke:rgb(0,0,0);stroke-width:%v" />`, y, y, r.barWidth)
	}
}

func (r *svgRender) drawX() {
	fmt.Fprint(r.f, `<line x1="0.05" x2="0.95" y1="0.05" y2="0.95" />`)
	fmt.Fprint(r.f, `<line x1="0.95" x2="0.05" y1="0.05" y2="0.95" />`)
}

func (r *svgRender) drawRedX() {
	fmt.Fprint(r.f, `<g stroke="rgb(255,0,0)">`)
	r.drawX()
	fmt.Fprint(r.f, `</g>`)
}

func (r *svgRender) drawO() {
	fmt.Fprint(r.f, `<circle cx="0.5" cy="0.5" r="0.45" fill="none" />`)
}

func (r *svgRender) drawRedO() {
	fmt.Fprint(r.f, `<g stroke="rgb(255,0,0)">`)
	r.drawO()
	fmt.Fprint(r.f, `</g>`)
}

func (r *svgRender) drawRowWin(row int) {
	y := float64(row)*(r.cellWidth+r.barWidth) + r.cellWidth/2
	fmt.Fprintf(r.f, `<line x1="%v" x2="%v" y1="0" y2="1" style="stroke:rgb(255,0,0);stroke-width:%v" />`,
		y, y, 2*r.barWidth)
}

func (r *svgRender) drawColumnWin(row int) {
	y := float64(row)*(r.cellWidth+r.barWidth) + r.cellWidth/2
	fmt.Fprintf(r.f, `<line x1="0" x2="1" y1="%v" y2="%v" style="stroke:rgb(255,0,0);stroke-width:%v" />`,
		y, y, 2*r.barWidth)
}

func (r *svgRender) drawWinDiagonal1() {
	fmt.Fprintf(r.f, `<line x1="0" x2="1" y1="0" y2="1" style="stroke:rgb(255,0,0);stroke-width:%v" />`,
		2*r.barWidth)
}
func (r *svgRender) drawWinDiagonal2() {
	fmt.Fprintf(r.f, `<line x1="1" x2="0" y1="0" y2="1" style="stroke:rgb(255,0,0);stroke-width:%v" />`,
		2*r.barWidth)
}

func (r *svgRender) drawWins(b BoardState) {
	// check all rows and columns
	for i := 0; i < b.Size(); i++ {
		if b.CheckRow(i) != NoPlayer {
			r.drawRowWin(i)
		}
		if b.CheckColumn(i) != NoPlayer {
			r.drawColumnWin(i)
		}
	}
	// check diagonals
	if b.CheckDiagonal1() != NoPlayer {
		r.drawWinDiagonal1()
	}

	if b.CheckDiagonal2() != NoPlayer {
		r.drawWinDiagonal2()
	}
}
