# tic-tac-toe
* "Крестики-нолики"*

** Базовая реализация сетевой консольной игры "Крестики-нолики" на языке Go с использованием Docker Compose.
   Она включает серверную и клиентскую часть, поддержку игрового поля 3x3 или 4x4, обработку ввода, проверку победителя, а также юнит-тесты.**
   
*** Директория проекта

```
tic-tac-toe/
├── client/
│   ├── client.go
├── server/
│   ├── server.go
├── core/
│   ├── game.go
│   ├── game_test.go
├── docker-compose.yml
├── Dockerfile.server
├── Dockerfile.client
```

core/game.go

**** Основная игровая логика: игровое поле, проверка победы и ничьи.

```
package core


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
```

core/game_test.go

**** Юнит-тесты.

```
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
```

server/server.go

**** Серверная часть.

```
package main

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
```

client/client.go

**** Клиент для взаимодействия с сервером.

```
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	for {
		printBoard()
		var x, y int
		fmt.Print("Enter your move (x y): ")
		fmt.Scan(&x, &y)

		move := map[string]int{"x": x, "y": y}
		data, _ := json.Marshal(move)

		resp, err := http.Post("http://server:8080/move", "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}

func printBoard() {
	resp, err := http.Get("http://server:8080/board")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

docker-compose.yml

**** Docker Compose для запуска клиента и сервера.

```
version: "3.9"
services:
  server:
    build:
      context: .
      dockerfile: Dockerfile.server
    ports:
      - "8080:8080"

  client:
    build:
      context: .
      dockerfile: Dockerfile.client
    depends_on:
      - server
```

**** Dockerfiles

Dockerfile.server

```
FROM golang:1.20
WORKDIR /app
COPY ./core ./core
COPY ./server ./server
RUN go mod init server && go mod tidy
CMD ["go", "run", "server/server.go"]
```

Dockerfile.client

```
FROM golang:1.20
WORKDIR /app
COPY ./client ./client
RUN go mod init client && go mod tidy
CMD ["go", "run", "client/client.go"]
```

Соберите проект с помощью docker-compose up.
