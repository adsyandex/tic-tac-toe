package main     //Клиент для взаимодействия с сервером.

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
