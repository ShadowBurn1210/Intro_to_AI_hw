package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"flag"
)
var url string
func init() {
	// parse the command line argument for url
	flag.StringVar(&url, "url", "http://37.27.208.205:55555", "API URL")
	flag.Parse()
}
// var url = "https://gomoku.martinsp.org/"
// var url = "http://37.27.208.205:55555"
var student_id = "221RDB477"
var depth = 2
// Game represents the game state from the API response
type Game struct {
	Color         string  `json:"color"`
	GameID        int     `json:"game_id"`
	GameStatus    string  `json:"game_status"`
	TimeRemaining float64 `json:"time_remaining"`
	Turn          string  `json:"turn"`
	Gameboard     [][]int `json:"gameboard"`
	RequestStatus string  `json:"request_status"`
}

// sendStartRequest sends a GET request to the specified endpoint and loads response data into the Game struct
func (g *Game) sendStartRequest() error {
	endpoint := fmt.Sprintf("%s/%s/start",url  ,student_id)
	response, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err
	}

	fmt.Println("Response status:", response.Status)
	fmt.Println("Response body:", string(body))

	// Unmarshal the JSON response directly into the Game struct
	err = json.Unmarshal(body, g)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return err
	}

	return nil
}

func (g *Game) move(moveX int, moveY int) error {

	endpoint := fmt.Sprintf("%s/%s/%d/%d/%d", url,student_id, g.GameID, moveX, moveY)
	response, err := http.Get(endpoint)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer response.Body.Close()

	// update the game state with the new data
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err
	}

	fmt.Println("Response status:", response.Status)
	fmt.Println("Response body:", string(body))
	if response.Status != "200 OK" {
		return fmt.Errorf("Move request failed: %s", response.Status)
	}

	// Check the response status
	if response.StatusCode == 400 {
		fmt.Println("ERROR: Invalid move detected!")
		return fmt.Errorf("invalid move: %s", string(body))
	}

	// Unmarshal the JSON response directly into the Game struct
	err = json.Unmarshal(body, g)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return err
	}

	// Print game data to verify it was correctly loaded
	fmt.Println("Game data updated successfully:")
	fmt.Printf("Color: %s\n", g.Color)
	fmt.Printf("Game ID: %d\n", g.GameID)
	fmt.Printf("Game Status: %s\n", g.GameStatus)
	fmt.Printf("Time Remaining: %.1f\n", g.TimeRemaining)
	fmt.Printf("Turn: %s\n", g.Turn)

	return nil
}

var globalBoard = [][]int{
}

func give_correct_board() [][]int {
	return globalBoard
}


func main() {
	startNewGame()
}
var lastRequestTime time.Time
func startNewGame() {
	newGame := Game{}

	err := newGame.sendStartRequest()
	if err != nil {
		fmt.Println("Failed to initialize game:", err)
		return
	}
	fmt.Println("Game initialized successfully!")

	// create a while loop to keep the game running
	for {
		if newGame.GameStatus != "ONGOING" {
			fmt.Println("Game is finished with status:", newGame.GameStatus)
			// Start a new game when the current game completes
			if newGame.GameStatus == "LEAVE" {
				fmt.Println("Exiting game...")
				return
			} else if newGame.GameStatus == "BLACKWON" || newGame.GameStatus == "WHITEWON" || newGame.GameStatus == "DRAW" {
				fmt.Println("Starting new game...")
				startNewGame()
				return
			}
			break
		} else {
			// Ensure at least 50ms between requests
			timeSinceLastRequest := time.Since(lastRequestTime)
			if timeSinceLastRequest < 50*time.Millisecond {
				time.Sleep(50*time.Millisecond - timeSinceLastRequest)
			}
		
			// Construct the endpoint URL
			endpoint := fmt.Sprintf("%s/%s/%d", url, student_id, newGame.GameID)
		
			// Send the HTTP request
			response, err := http.Get(endpoint)
			if err != nil {
				fmt.Printf("Error making request: %v", err)
				return
			}
			defer response.Body.Close()
		
			// Record the time of the last request
			lastRequestTime = time.Now()
			// Update the game state with the new data
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}

			err = json.Unmarshal(body, &newGame)
			if err != nil {
				fmt.Println("Error unmarshalling response:", err)
				return
			}

			// Print board with proper indexing for better visualization
			board := convertGameboard(newGame.Gameboard)
			printBoardWithIndexing(board)
			globalBoard = board

			// Send a move request if it's our turn
			if newGame.Turn == newGame.Color {
				bestMove := findBestMove(board, newGame.Color, depth)

				fmt.Printf("Best move: %v\n", bestMove)
				// Send the move
				err = newGame.move(bestMove[0], bestMove[1])
				if err != nil {
					fmt.Println("Move failed, trying another position...")
					fmt.Print(board)
					fmt.Print(newGame.GameID)
					fmt.Print(newGame.Color)
				} else {
					fmt.Printf("Move [%d,%d] recorded successfully\n", bestMove[0], bestMove[1])
				}
			}
		}
	}
}

// Helper function to convert the gameboard from the API format to the algorithm format
// Ensuring proper 0,0 indexing
func convertGameboard(gameboard [][]int) [][]int {
	board := make([][]int, len(gameboard))
	for i := range gameboard {
		board[i] = make([]int, len(gameboard[i]))
		for j := range gameboard[i] {
			board[i][j] = gameboard[i][j]
		}
	}
	return board
}

// Print the board with proper indexing for debugging
func printBoardWithIndexing(board [][]int) {
	fmt.Println("Board with indexing (x,y):")
	fmt.Print("  ")
	for j := 0; j < len(board[0]); j++ {
		fmt.Printf("%2d ", j)
	}
	fmt.Println()

	for i := 0; i < len(board); i++ {
		fmt.Printf("%2d ", i)
		for j := 0; j < len(board[i]); j++ {
			fmt.Printf("%2d ", board[i][j])
		}
		fmt.Println()
	}
}
