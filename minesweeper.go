package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var PlayerName string
var PlayerVersion int
var ServerHost string
var CurrentGameID string

// ReadIn a line of input
func ReadIn() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

// PrintOut a line of output
func PrintOut(line string) {
	fmt.Println(line)
}

// Query "info" endpoint and return its result
func GameInfo() string {
	endpoint := "http://" + ServerHost + "/info"
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Fatalf("failed to connect to server at " + endpoint)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

// Query "new" endpoint to create a new game and return its ID
func NewGame() {
	endpoint := "http://" + ServerHost + "/new"
	resp, err := http.PostForm(endpoint, url.Values{
		"name":    {PlayerName},
		"version": {fmt.Sprintf("%v", PlayerVersion)},
	})
	if err != nil {
		log.Fatalf("failed to connect to server at " + endpoint)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	CurrentGameID = string(body)
}

// Read in and parse x,y coordinates of the player's guess
func GetTarget() (int64, int64) {
	text := ReadIn()
	split := strings.Split(strings.TrimSpace(text), ",")
	if len(split) != 2 {
		log.Fatalf("Expecting guess coordinates in x,y form.")
	}

	x, err1 := strconv.ParseInt(split[0], 0, 64)
	y, err2 := strconv.ParseInt(split[1], 0, 64)

	if (err1 != nil) || (err2 != nil) {
		log.Fatalf("Non-integer input for guess coordinates")
	}

	return x, y
}

// Query the "open" endpoint with the player's guess and return its result
func OpenCell(x int64, y int64) string {
	xs := fmt.Sprintf("%v", x)
	ys := fmt.Sprintf("%v", y)
	endpoint := "http://" + ServerHost + "/open"
	resp, err := http.PostForm(endpoint, url.Values{
		"id": {CurrentGameID},
		"x":  {xs},
		"y":  {ys},
	})
	if err != nil {
		log.Fatalf("failed to connect to server at " + endpoint)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

// Play one instance of a game
func PlayGame() {
	// Initiate a game with the client by printing the board parameters.
	gameInfo := GameInfo()
	PrintOut(gameInfo)

	// Get first guess before creating a game to avoid lots of empty games if
	// the program exits early.
	x, y := GetTarget()

	// Start a game with the server.
	NewGame()

	for true {
		// Open the cell guessed by the player and print the result.
		result := OpenCell(x, y)
		PrintOut(result)

		// If the result isn't numeric (the number of neighbors),
		// then the game has ended for one reason or another.
		_, notAnInt := strconv.ParseInt(result, 0, 64)
		if notAnInt != nil {
			break
		}

		x, y = GetTarget()
	}
}

// Plays the specified number of games
func PlayGames(howMany int) {
	for i := 0; i < howMany; i++ {
		PlayGame()
	}
}

func main() {
	games := flag.Int("games", 1, "how many games to play")
	flag.StringVar(&PlayerName, "name", "", "your name")
	flag.IntVar(&PlayerVersion, "version", -1, "the version of your application")
	flag.StringVar(&ServerHost, "server", "minesweeper.nm.io", "server to connect to")
	flag.Parse()

	if PlayerName == "" {
		log.Fatalf("name is a required flag for this program.")
	}

	PlayGames(*games)
}
