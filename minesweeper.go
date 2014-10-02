package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var PlayerName string
var PlayerVersion int
var ServerHost string
var CurrentGameID string
var Stdin io.Reader
var Stdout io.Writer

// ReadIn a line of input
func ReadIn() string {
	reader := bufio.NewReader(Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

// PrintOut a line of output
func PrintOut(line string) {
	fmt.Fprintln(Stdout, line)
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
func PlayGame() bool {
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
			if result == "win" {
				return true
			} else {
				return false
			}
		}

		x, y = GetTarget()
	}
	return false
}

// Plays the specified number of games, returns the number of games won
func PlayGames(command []string, howMany int) int {
	wins := 0
	for i := 0; i < howMany; i++ {
		if len(command) >= 1 {
			LaunchAI(command)
		} else {
			Stdin = os.Stdin
			Stdout = os.Stdout
		}

		win := PlayGame()
		if win {
			wins++
		}
	}

	return wins
}

// Launches the AI program, capturing its stdout and stdin
func LaunchAI(command []string) {
	cmd := exec.Command(command[0], command[1:]...)

	Stdout, _ = cmd.StdinPipe()
	Stdin, _ = cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		log.Fatalf("error starting the command: %s", err)
	}
}

func main() {
	games := flag.Int("games", 1, "how many games to play")
	flag.StringVar(&PlayerName, "name", "", "your name")
	flag.IntVar(&PlayerVersion, "version", -1, "the version of your application")
	flag.StringVar(&ServerHost, "server", "minesweeper.nm.io", "server to connect to")
	flag.Parse()
	command := flag.Args()

	if PlayerName == "" {
		log.Fatalf("name is a required flag for this program.")
	}

	wins := PlayGames(command, *games)
	fmt.Printf("%v games played, %v games won\n", *games, wins)
}
