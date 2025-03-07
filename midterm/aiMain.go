package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// make a construct that records all the moves made
type Move struct {
	x int
	y int
}

var student_id = "221RDB477"

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
	endpoint := fmt.Sprintf("http://37.27.208.205:55555/%s/start", student_id)
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

	endpoint := fmt.Sprintf("http://37.27.208.205:55555/%s/%d/%d/%d", student_id, g.GameID, moveX, moveY)
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

// Check if a move is valid on the board
func isValidMove(board [][]int, x, y int) bool {
	if x < 0 || y < 0 || x >= len(board) || y >= len(board) {
		return false
	}

	return board[x][y] == EMPTY
}

// Get best move ensuring it hasn't been played before
func getBestMoveWithValidation(board [][]int, color string, maxDepth int) [2]int {
	bestMove := findBestMove(board, color, maxDepth)

	// Check if the move has been played already
	for isMoveAlreadyPlayed(bestMove) {
		// If already played, mark the position as occupied and recalculate
		fmt.Printf("Move [%d,%d] already played, recalculating...\n", bestMove[0], bestMove[1])
		// Mark it as our piece to avoid it in the future
		if color == "BLACK" {
			board[bestMove[0]][bestMove[1]] = 1 // Mark as unavailable
		} else {
			board[bestMove[0]][bestMove[1]] = 2 // Mark as unavailable
		}
		bestMove = findBestMove(board, color, maxDepth)
	}

	return bestMove
}

// Check if a move has already been played
func isMoveAlreadyPlayed(move [2]int) bool {
	for _, playedMove := range movesPlayed {
		if playedMove.x == move[0] && playedMove.y == move[1] {
			
			return true
		}
	}

	return false
}

// Find a valid fallback move when the AI suggests an invalid one
func findSafeMove(board [][]int) ([2]int, error) {
	// First try to find an empty spot near existing pieces
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] != EMPTY {
				// Check neighbors
				for dx := -1; dx <= 1; dx++ {
					for dy := -1; dy <= 1; dy++ {
						nx, ny := i+dx, j+dy
						if nx >= 0 && nx < len(board) && ny >= 0 && ny < len(board) &&
							board[nx][ny] == EMPTY {
							return [2]int{nx, ny}, nil
						}
					}
				}
			}
		}
	}

	// If no empty spots near pieces, find any empty spot
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] == EMPTY {
				return [2]int{i, j}, nil
			}
		}
	}

	return [2]int{}, fmt.Errorf("no valid moves available")
}

var movesPlayed = []Move{}

// Reset the moves played list at the start of each game
func (g *Game) resetMovesPlayed() {
	movesPlayed = []Move{}
}

// Sync the moves played with the actual board state
func syncMovesWithBoard(board [][]int) {
	// Clear the moves list
	movesPlayed = []Move{}

	// Scan the board and add all non-empty positions to the moves list
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] != EMPTY {
				movesPlayed = append(movesPlayed, Move{x: i, y: j})
			}
		}
	}

	fmt.Printf("Synced %d moves with the board state\n", len(movesPlayed))
}

func main() {
	newGame := Game{}

	err := newGame.sendStartRequest()
	if err != nil {
		fmt.Println("Failed to initialize game:", err)
		return
	}
	fmt.Println("Game initialized successfully!")

	// Reset moves played at the start of a new game
	newGame.resetMovesPlayed()

	// create a while loop to keep the game running
	for {
		if newGame.GameStatus != "ONGOING" {
			fmt.Println("Game is finished!")
			break
		} else {
			// update the game state
			endpoint := fmt.Sprintf("http://37.27.208.205:55555/%s/%d", student_id, newGame.GameID)
			response, err := http.Get(endpoint)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer response.Body.Close()

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


			// Send a move request if it's our turn
			if newGame.Turn == newGame.Color {
				// Sync our moves tracking with the actual board state
				syncMovesWithBoard(board)

				// Convert the gameboard to the format expected by the algorithm
				board := convertGameboard(board)

				// Find the best move using the Minimax algorithm with validation
				bestMove := getBestMoveWithValidation(board, newGame.Color, 3) // Adjust depth as needed

				// Verify that the move is valid
				maxAttempts := 3
				attempts := 0

				for attempts < maxAttempts {
					fmt.Printf("Best move: %v\n", bestMove)

					// Validate the move before sending
					if !isValidMove(board, bestMove[0], bestMove[1]) {
						fmt.Println("AI tried to make an invalid move, finding fallback move...")
						fallbackMove, err := findSafeMove(board)
						if err != nil {
							fmt.Println("Error finding fallback move:", err)
							return
						}
						bestMove = fallbackMove
						fmt.Printf("Using fallback move instead: %v\n", bestMove)
					}

					// Send the move
					err = newGame.move(bestMove[0], bestMove[1])
					if err != nil {
						// If move failed, try to find another move
						if attempts < maxAttempts-1 {
							fmt.Println("Move failed, trying another position...")
							// Mark the failed position as non-empty to avoid trying it again
							board[bestMove[0]][bestMove[1]] = 999 // Using 999 to mark invalid positions

							// Get a new fallback move
							fallbackMove, err := findSafeMove(board)
							if err != nil {
								fmt.Println("Error finding fallback move:", err)
								return
							}
							bestMove = fallbackMove
							attempts++
						} else {
							fmt.Println("Max attempts reached, giving up on this turn")
							// make a random move to any empty spot
							for i := 0; i < len(board); i++ {
								for j := 0; j < len(board[i]); j++ {
									if board[i][j] == EMPTY {
										err = newGame.move(i, j)
										if err != nil {
											fmt.Println("Error making random move:", err)
											return
										}
									}
								}
							}

						}
					} else {
						// Move was successful
						// Add the move to the list of moves played
						movesPlayed = append(movesPlayed, Move{x: bestMove[0], y: bestMove[1]})
						fmt.Printf("Move [%d,%d] recorded successfully\n", bestMove[0], bestMove[1])
						break
					}
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
