package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
)

var (
	answerGrid [9][9]int
	result     = make(chan bool)
	upgrader   = websocket.Upgrader{}
)

//Compare Puzzle with the Answer Grid
func checkAnswer() bool {
	return (answerGrid == replicatedGrid)
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "sudoku.html")
}

//Generate Stream for Sending Over Web Socket
func generateStream() string {
	var puzzleDataStream string
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			puzzleDataStream = puzzleDataStream + strconv.Itoa(sudokuGrid[i][j])
		}
	}
	return puzzleDataStream
}

func newGameHandler(rw http.ResponseWriter, req *http.Request) {

	c, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Print("Upgrade : ", err)
	}

	str := generateStream()
	c.WriteMessage(websocket.TextMessage, []byte(str))

	for {
		//receive data from Web Browser
		_, recvData, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		//Extracting data from UI

		data := string(recvData)
		split := strings.Split(data, ",")
		value, _ := strconv.Atoi(split[0])
		row, _ := strconv.Atoi(split[1])
		col, _ := strconv.Atoi(split[2])
		fmt.Println("display --->", data)

		answerGrid[row][col] = value

		if answerGrid[row][col] != replicatedGrid[row][col] {
			c.WriteMessage(websocket.TextMessage, []byte("Violation"))
		} else {
			if checkAnswer() {
				c.WriteMessage(websocket.TextMessage, []byte("WIN"))
				result <- true
				return
			}
		}
	}
}

func InitRouter() (router *mux.Router) {

	router = mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	router.HandleFunc("/ws", newGameHandler).Methods(http.MethodGet)

	return
}

func serverStart() {

	answerGrid = replicateOriginalGrid(sudokuGrid)

	router := InitRouter()
	server := negroni.Classic()
	server.UseHandler(router)
	server.Run(":9009")

}

func userDisplay() {

	t := time.NewTimer(3 * time.Minute)

	go serverStart()
	select {
	case <-result:
		fmt.Println("Result of Game is : ", result)
	case <-t.C:
		fmt.Println("Game Timed Out")
	}
}
