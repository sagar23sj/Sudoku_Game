package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Score struct {
	Name string `json:"Name"`
	Time string `json:"Time"`
}

var upgrader = websocket.Upgrader{}

var (
	DB_USER     string = os.Getenv("DATABASE_USERNAME")
	DB_PASSWORD string = os.Getenv("DATABASE_PASSWORD")
	DB_NAME     string = os.Getenv("DATABASE_NAME")
)

//Compare Puzzle with the Answer Grid
func (s *Sudoku) checkAnswer() bool {
	return (s.answerGrid == s.replicatedGrid)
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "game.html")
}

func saveScore(userTime time.Duration, name string) {
	hours := int(userTime / time.Hour)
	minutes := int(userTime / time.Minute)
	seconds := int(userTime / time.Second)
	seconds = seconds - minutes*60
	current := time.Now()
	date := current.Format("2006-01-02")
	usertime := strconv.Itoa(hours) + ":" + strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds)

	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	sql := "INSERT INTO Scores(Name, Time, Date) VALUES (?,?,?)"
	insert, err := db.Query(sql, name, usertime, date)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

func getTopScores() string {

	var top []Score
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT Name, Time FROM Scores Order by Time LIMIT 5")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var tag Score
		err = results.Scan(&tag.Name, &tag.Time)
		if err != nil {
			panic(err.Error())
		}
		top = append(top, tag)
	}

	jsonData, err := json.Marshal(top)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	return string(jsonData)
}

//Generate Stream for Sending Over Web Socket
func (s *Sudoku) generateStream() string {
	var puzzleDataStream string
	for i := 0; i < s.gridSize; i++ {
		for j := 0; j < s.gridSize; j++ {
			puzzleDataStream = puzzleDataStream + strconv.Itoa(s.sudokuGrid[i][j])
		}
	}
	return puzzleDataStream
}

func newGameHandler(rw http.ResponseWriter, req *http.Request) {
	// Start Timer for current game
	start := time.Now()

	type Score struct {
		name string
		time []int
	}
	c, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Print("Upgrade : ", err)
	}

	// To get difficuly level from UI
	_, level, err := c.ReadMessage()

	gameLevel := string(level)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Difficulty Level : ", gameLevel)
	c.WriteMessage(websocket.TextMessage, []byte(getTopScores()))

	Game := Sudoku{}
	Game.initializeGame(9, 3, gameLevel)
	Game.createPuzzle(gameLevel)
	fmt.Println("phase 2")
	Game.answerGrid = replicateOriginalGrid(Game.sudokuGrid)
	fmt.Println("phase 3")
	str := Game.generateStream()
	fmt.Println("phase 4")
	c.WriteMessage(websocket.TextMessage, []byte(str))
	fmt.Println("phase 5")

	for {
		// score := Score{}
		var userData map[string]int

		_, recvData, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		//Extracting data from UI
		_ = json.Unmarshal(recvData, &userData)
		value := userData["value"]
		row := userData["row"]
		col := userData["col"]

		Game.answerGrid[row][col] = value

		if Game.answerGrid[row][col] != Game.replicatedGrid[row][col] {
			c.WriteMessage(websocket.TextMessage, []byte("Violation"))
		} else {
			w := Game.checkAnswer()
			fmt.Println("status : ", w)
			if Game.checkAnswer() {
				c.WriteMessage(websocket.TextMessage, []byte("WIN"))
				userTiming := time.Since(start)

				//Getting Player Name
				_, nameData, _ := c.ReadMessage()
				name := string(nameData)
				saveScore(userTiming, name)
				break
			}
		}
	}

}

func InitRouter() (router *mux.Router) {

	router = mux.NewRouter()

	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./assets/"))))
	router.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	router.HandleFunc("/ws", newGameHandler).Methods(http.MethodGet)

	return
}

func serverStart() {

	router := InitRouter()
	server := negroni.Classic()
	server.UseHandler(router)
	server.Run(":3000")

}

func main() {
	serverStart()
}
