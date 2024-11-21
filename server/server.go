package main      // серверная часть

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"tic-tac-toe/core"
)

var game *core.Game
var mu sync.Mutex

func main() {
	var err error
	game, err = core.NewGame(3)
	if err != nil {
		log.Fatalf("failed to initialize game: %v", err)
	}

	http.HandleFunc("/move", handleMove)
	http.HandleFunc("/board", handleBoard)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var move struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	json.NewDecoder(r.Body).Decode(&move)

	err := game.MakeMove(move.X, move.Y)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	winner, finished := game.CheckWinner()
	response := map[string]interface{}{
		"board":    game.Board,
		"finished": finished,
		"winner":   winner,
	}
	json.NewEncoder(w).Encode(response)
}

func handleBoard(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	json.NewEncoder(w).Encode(game.Board)
}
