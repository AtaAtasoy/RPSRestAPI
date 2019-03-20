//Author: Huseyin Ata Atasoy
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	s "strings"
	"time"
)

// Pair struct for storing scores
type Pair struct {
	Left  int
	Right int
}

// Global variable to keep track of rounds in a game
var roundsPerGame map[int]int
var scoresPerGame map[int]Pair

// AssignID to a game
func AssignID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()
}

// DisplayInstructions for the game
func DisplayInstructions(w http.ResponseWriter, r *http.Request, gameID int) {

	fmt.Fprintf(w, "<br><a href=\"/play?choose=rock&id=")
	fmt.Fprintf(w, strconv.Itoa(gameID))
	fmt.Fprintf(w, "\">Play Rock</a>")

	fmt.Fprintf(w, "<br><a href=\"/play?choose=paper&id=")
	fmt.Fprintf(w, strconv.Itoa(gameID))
	fmt.Fprintf(w, "\">Play Paper</a>")

	fmt.Fprintf(w, "<br><a href=\"/play?choose=scissors&id=")
	fmt.Fprintf(w, strconv.Itoa(gameID))
	fmt.Fprintf(w, "\">Play Scissors</a>")
}

// ProcessGame for the next move sequence
func ProcessGame(moveAI string, moveP string, gameID int) string {
	if _, ok := scoresPerGame[gameID]; !ok {
		scoresPerGame[gameID] = Pair{0, 0}
	}

	roundsPerGame[gameID]--

	if moveAI == "ROCK" && moveP == "SCISSORS" {
		scoresPerGame[gameID] = Pair{scoresPerGame[gameID].Left + 1, scoresPerGame[gameID].Right}
		return "YOU LOST THIS ROUND"
	} else if moveAI == "ROCK" && moveP == "PAPER" {
		scoresPerGame[gameID] = Pair{scoresPerGame[gameID].Left, scoresPerGame[gameID].Right + 1}
		return "YOU WON THIS ROUND"
	} else if moveAI == "PAPER" && moveP == "ROCK" {
		scoresPerGame[gameID] = Pair{scoresPerGame[gameID].Left + 1, scoresPerGame[gameID].Right}
		return "YOU LOST THIS ROUND"
	} else if moveAI == "PAPER" && moveP == "SCISSORS" {
		scoresPerGame[gameID] = Pair{scoresPerGame[gameID].Left, scoresPerGame[gameID].Right + 1}
		return "YOU WON THIS ROUND"
	} else if moveAI == "SCISSORS" && moveP == "ROCK" {
		scoresPerGame[gameID] = Pair{scoresPerGame[gameID].Left, scoresPerGame[gameID].Right + 1}
		return "YOU WON THIS ROUND"
	} else if moveAI == "SCISSORS" && moveP == "PAPER" {
		scoresPerGame[gameID] = Pair{scoresPerGame[gameID].Left + 1, scoresPerGame[gameID].Right}
		return "YOU LOST THIS ROUND"
	}

	roundsPerGame[gameID]++
	return "TIE"
}

// CreateNewGame for the user
func CreateNewGame(w http.ResponseWriter, r *http.Request) {

	id := AssignID()

	qs := r.URL.Query()
	if _, ok := qs["round"]; ok {
		rounds, _ := strconv.Atoi(qs.Get("round"))
		roundsPerGame[id] = rounds
	} else {
		roundsPerGame[id] = 1
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<!DOCTYPE html>\nNew Rock-Paper-Scissors game started<br>Session ID: %d", id)

	DisplayInstructions(w, r, id)
}

// PlayRound for the player
func PlayRound(w http.ResponseWriter, r *http.Request) {

	move := []string{"ROCK", "PAPER", "SCISSORS"}[rand.Int63n(3)]

	qs := r.URL.Query()
	gameID, _ := strconv.Atoi(qs.Get("id"))

	if qs["choose"] != nil {
		totalRounds := roundsPerGame[gameID]
		if roundsPerGame[gameID] > 0 {
			result := ProcessGame(move, s.ToUpper(qs.Get("choose")), gameID)

			fmt.Fprintf(w, "I played: %s\nYou played:%s\n", move, s.ToUpper(qs.Get("choose")))
			fmt.Fprintf(w, "%s\nScore\n%d - %d", result, scoresPerGame[gameID].Left, scoresPerGame[gameID].Right)

			if roundsPerGame[gameID] == 0 ||
				scoresPerGame[gameID].Left > totalRounds/2 ||
				scoresPerGame[gameID].Right > totalRounds/2 {
				if scoresPerGame[gameID].Left > scoresPerGame[gameID].Right {
					fmt.Fprintf(w, "\nYOU LOST THE GAME !")
				} else if scoresPerGame[gameID].Left < scoresPerGame[gameID].Right {
					fmt.Fprintf(w, "\nYOU WON THE GAME !")
				} else {
					fmt.Fprintf(w, "\nIT'S A TIE !")
				}
				roundsPerGame[gameID] = 0
			}

			// Debugging
			fmt.Println("YOU PLAYED:", s.ToUpper(qs.Get("choose")))
			fmt.Println(roundsPerGame[gameID], "rounds left in game :", gameID)
		} else {
			fmt.Fprintf(w, "Game finished !")
			fmt.Fprintf(w, "\nScore: %d - %d", scoresPerGame[gameID].Left, scoresPerGame[gameID].Right)
		}
	}
}

func hadleRequests() {
	http.HandleFunc("/newGame", CreateNewGame)
	http.HandleFunc("/play", PlayRound)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	roundsPerGame = make(map[int]int)
	scoresPerGame = make(map[int]Pair)
	hadleRequests()
}
