package main

import (
	tictac "../.."
	"fmt"
	"os"
	"io"
	"text/template"
)

func main() {
	head := tictac.NewDefaultState(tictac.PlayerX)
	head.FindAllChildStates()

	f, err := os.OpenFile("out.svg", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	const S int = 4096

	fmt.Fprintf(f, "<svg viewBox=\"0 0 1 1\" width=\"%v\" height=\"%v\">\n", S, S)
	fmt.Fprint(f, `<g stroke-width="0.05" stroke-linecap="round" stroke="rgb(0,0,0)">`)
	fmt.Fprint(f, `<rect width="1" height="1" style="fill:rgb(255,255,255);stroke:rgb(255,255,255)"/>`)
	renderState(f, &head)
	fmt.Fprint(f, `</g>`)
	fmt.Fprint(f, "</svg>\n")
}

const (
	// TODO work for boards other than 3x3
	barWidth  float64 = 0.01
	cellWidth float64 = (1.0 - 2.0*barWidth) / 3
	buffer            = 0.98
)

func translation(pos tictac.Pos) string {

	dx := float64(pos.X) * (cellWidth + barWidth)
	dy := float64(pos.Y) * (cellWidth + barWidth)
	scale := cellWidth
	return fmt.Sprintf("<g transform=\"translate(%v %v) scale(%v) translate(%v %v) scale(%v)\">\n",
		dx, dy, scale, (1-buffer)/2, (1-buffer)/2, buffer)
}

func renderState(f io.Writer, s *tictac.StateTreeNode) {
	if p := s.Board.CheckWin(); p != tictac.NoPlayer {
		return
	}
	if (s.Board.Full()) {
		return
	}

	drawBars(f)

	// draw squares that are not occupied
	for _, pos := range s.Board.OccupiedCells() {
		fmt.Fprint(f, translation(pos))
		// print an X
		if (s.Board.GetPos(pos) == tictac.PlayerX) {
			drawX(f)
		} else {
			drawO(f)
		}
		fmt.Fprint(f, "</g>\n")
	}

	if (s.OurPlayer == s.NextPlayer) {
		// draw our thing, then go again
		child := s.Children[s.BestChild]
		pos := tictac.BoardDiff(s.Board, child.Board)
		fmt.Fprint(f, translation(pos))
		// print an X
		drawRedX(f)
		fmt.Fprint(f, "</g>\n")

		s = &child
	}

	// draw each child inside the cell it represents
	for _, child := range s.Children {
		diff := tictac.BoardDiff(s.Board, child.Board)
		fmt.Fprint(f, translation(diff))
		renderState(f, &child)
		fmt.Fprint(f, "</g>\n")
	}
}

func drawBars(f io.Writer) {
	/*
	format:
		X1	X2
		|	|
	Y1--|---|---
		|	|
	Y2--|---|---
		|	|
	*/
	fmt.Fprint(f, "<g>\n")
	defer fmt.Fprint(f, "</g>\n")

	t, err := template.New("BoardTemplate").Parse(`
	<line x1="0" x2="1" y1="{{.Y1}}" y2="{{.Y1}}" style="stroke:rgb(0,0,0);stroke-width:{{.BarWidth}}" />
	<line x1="0" x2="1" y1="{{.Y2}}" y2="{{.Y2}}" style="stroke:rgb(0,0,0);stroke-width:{{.BarWidth}}" />
	<line x1="{{.X1}}" x2="{{.X1}}" y1="0" y2="1" style="stroke:rgb(0,0,0);stroke-width:{{.BarWidth}}" />
	<line x1="{{.X2}}" x2="{{.X2}}" y1="0" y2="1" style="stroke:rgb(0,0,0);stroke-width:{{.BarWidth}}" />
	`)
	if err != nil {
		panic(err)
	}
	type BoardGrid struct {
		X1, X2, Y1, Y2 float64
		BarWidth       float64
	}
	bg := BoardGrid{
		X1:       cellWidth + barWidth/2,
		X2:       2*cellWidth + barWidth + barWidth/2,
		Y1:       cellWidth + barWidth/2,
		Y2:       2*cellWidth + barWidth + barWidth/2,
		BarWidth: barWidth,
	}
	t.Execute(f, bg)
}

func drawX(f io.Writer) {
	fmt.Fprint(f, `<line x1="0.05" x2="0.95" y1="0.05" y2="0.95" />`)
	fmt.Fprint(f, `<line x1="0.95" x2="0.05" y1="0.05" y2="0.95" />`)
}

func drawRedX(f io.Writer) {
	fmt.Fprint(f, `<g stroke="rgb(255,0,0)">`)
	drawX(f)
	fmt.Fprint(f, `</g>`)
}

func drawO(f io.Writer) {
	fmt.Fprint(f, `<circle cx="0.5" cy="0.5" r="0.45" stroke="black" fill="none" />`)
}
