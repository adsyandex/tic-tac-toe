package core

import "testing"

func TestGameCreation(t *testing.T) {
	game, err := NewGame(3)
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}
	if len(game.Board) != 3 {
		t.Fatalf("invalid board size")
	}
}

func TestMakeMove(t *testing.T) {
	game, _ := NewGame(3)
	err := game.MakeMove(0, 0)
	if err != nil {
		t.Fatalf("failed to make move: %v", err)
	}
	if game.Board[0][0] != "X" {
		t.Fatalf("move not recorded")
	}
}

func TestWinner(t *testing.T) {
	game, _ := NewGame(3)
	game.MakeMove(0, 0)
	game.MakeMove(1, 0)
	game.MakeMove(0, 1)
	game.MakeMove(1, 1)
	game.MakeMove(0, 2)
	winner, finished := game.CheckWinner()
	if !finished || winner != "X" {
		t.Fatalf("winner not detected")
	}
}
