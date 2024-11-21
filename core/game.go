package core     // основная лрогика: игровое поле, проверка победы и ничья.

import (
	"errors"
	"fmt"
)

type Game struct {
	Board   [][]string
	Current string
	Size    int
}

func NewGame(size int) (*Game, error) {
	if size < 3 || size > 4 {
		return nil, errors.New("unsupported board size: use 3x3 or 4x4")
	}
	board := make([][]string, size)
	for i := range board {
		board[i] = make([]string, size)
	}
	return &Game{
		Board:   board,
		Current: "X",
		Size:    size,
	}, nil
}

func (g *Game) PrintBoard() {
	for _, row := range g.Board {
		fmt.Println(row)
	}
}

func (g *Game) MakeMove(x, y int) error {
	if x < 0 || y < 0 || x >= g.Size || y >= g.Size {
		return errors.New("move out of bounds")
	}
	if g.Board[x][y] != "" {
		return errors.New("cell already occupied")
	}
	g.Board[x][y] = g.Current
	if g.Current == "X" {
		g.Current = "O"
	} else {
		g.Current = "X"
	}
	return nil
}

func (g *Game) CheckWinner() (string, bool) {
	// Check rows, columns, and diagonals
	for i := 0; i < g.Size; i++ {
		if g.checkLine(g.Board[i]) || g.checkLine(g.getColumn(i)) {
			return g.Board[i][0], true
		}
	}
	if g.checkLine(g.getDiagonal1()) || g.checkLine(g.getDiagonal2()) {
		return g.Board[0][0], true
	}
	// Check draw
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == "" {
				return "", false
			}
		}
	}
	return "Draw", true
}

func (g *Game) checkLine(line []string) bool {
	first := line[0]
	if first == "" {
		return false
	}
	for _, cell := range line {
		if cell != first {
			return false
		}
	}
	return true
}

func (g *Game) getColumn(col int) []string {
	column := make([]string, g.Size)
	for i := 0; i < g.Size; i++ {
		column[i] = g.Board[i][col]
	}
	return column
}

func (g *Game) getDiagonal1() []string {
	diagonal := make([]string, g.Size)
	for i := 0; i < g.Size; i++ {
		diagonal[i] = g.Board[i][i]
	}
	return diagonal
}

func (g *Game) getDiagonal2() []string {
	diagonal := make([]string, g.Size)
	for i := 0; i < g.Size; i++ {
		diagonal[i] = g.Board[i][g.Size-i-1]
	}
	return diagonal
}
