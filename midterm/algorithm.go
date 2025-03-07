package main

import (
	"fmt"
	"math"
)

const (
	EMPTY = 0
	BLACK = 1
	WHITE = 2
)

func findBestMove(board [][]int, color string, maxDepth int) [2]int {
	player := toPlayer(color)
	opponent := getOpponent(player)

	bestScore := math.MinInt32
	var bestMove [2]int

	moves, visited := generateMoves(board)

	for _, move := range moves {
		x, y := move[0], move[1]
		newBoard := makeCopy(board)
		newBoard[x][y] = player

		score := minimax(newBoard, maxDepth, math.MinInt32, math.MaxInt32, false, player, opponent)
		if score > bestScore {
			bestScore = score
			bestMove = [2]int{x, y}
		}
	}

	fmt.Printf("Visited %d moves\n", len(visited))
	fmt.Println(visited)

	return bestMove
}

// Print the board with proper indexing for debugging
func printBoardWithIndexing2(board [][]int) {
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

func minimax(board [][]int, depth int, alpha int, beta int, maximizingPlayer bool, player int, opponent int) int {
	if depth == 0 || isTerminal(board, player, opponent) {
		return evaluate(board, player, opponent)
	}

	if maximizingPlayer {
		maxVal := math.MinInt32
		moves, _ := generateMoves(board)

		for _, move := range moves {
			x, y := move[0], move[1]
			newBoard := makeCopy(board)
			newBoard[x][y] = player

			val := minimax(newBoard, depth-1, alpha, beta, false, player, opponent)
			if val > maxVal {
				maxVal = val
			}
			if val > alpha {
				alpha = val
			}
			if alpha >= beta {
				break
			}
		}
		return maxVal
	} else {
		minVal := math.MaxInt32
		moves, _ := generateMoves(board)
		for _, move := range moves {
			x, y := move[0], move[1]
			newBoard := makeCopy(board)
			newBoard[x][y] = opponent

			val := minimax(newBoard, depth-1, alpha, beta, true, player, opponent)
			if val < minVal {
				minVal = val
			}
			if val < beta {
				beta = val
			}
			if beta <= alpha {
				break
			}
		}
		return minVal
	}
}

func evaluate(board [][]int, player int, opponent int) int {
	if hasWon(board, player) {
		return 1000000
	}
	if hasWon(board, opponent) {
		return -1000000
	}
	if isBoardFull(board) {
		return 0
	}

	score := 0
	size := len(board)
	directions := [][]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if board[x][y] == player {
				for _, dir := range directions {
					dx, dy := dir[0], dir[1]
					count := 1
					openStart, openEnd := true, true

					prevX, prevY := x-dx, y-dy
					if prevX < 0 || prevX >= size || prevY < 0 || prevY >= size || board[prevX][prevY] != EMPTY {
						openStart = false
					}

					i := 1
					for ; i <= 4; i++ {
						nx, ny := x+dx*i, y+dy*i
						if nx < 0 || nx >= size || ny < 0 || ny >= size || board[nx][ny] != player {
							break
						}
						count++
					}

					endX, endY := x+dx*i, y+dy*i
					if endX < 0 || endX >= size || endY < 0 || endY >= size || board[endX][endY] != EMPTY {
						openEnd = false
					}

					score += evaluateConsecutive(count, openStart, openEnd)
				}
			} else if board[x][y] == opponent {
				for _, dir := range directions {
					dx, dy := dir[0], dir[1]
					count := 1
					openStart, openEnd := true, true

					prevX, prevY := x-dx, y-dy
					if prevX < 0 || prevX >= size || prevY < 0 || prevY >= size || board[prevX][prevY] != EMPTY {
						openStart = false
					}

					i := 1
					for ; i <= 4; i++ {
						nx, ny := x+dx*i, y+dy*i
						if nx < 0 || nx >= size || ny < 0 || ny >= size || board[nx][ny] != opponent {
							break
						}
						count++
					}

					endX, endY := x+dx*i, y+dy*i
					if endX < 0 || endX >= size || endY < 0 || endY >= size || board[endX][endY] != EMPTY {
						openEnd = false
					}

					score -= evaluateConsecutive(count, openStart, openEnd)
				}
			}
		}
	}

	return score
}

func evaluateConsecutive(count int, openStart bool, openEnd bool) int {
	switch count {
	case 4:
		if openStart || openEnd {
			return 100000
		}
		return 50000
	case 3:
		if openStart && openEnd {
			return 10000
		} else if openStart || openEnd {
			return 5000
		}
		return 1000
	case 2:
		if openStart && openEnd {
			return 100
		} else if openStart || openEnd {
			return 50
		}
		return 10
	case 1:
		if openStart && openEnd {
			return 10
		}
		return 0
	default:
		if count >= 5 {
			return 1000000
		}
		return 0
	}
}

var once bool = true
var twice bool = true

func generateMoves(board [][]int) ([][2]int, [][2]int) {
	size := len(board)
	moves := make([][2]int, 0)
	visited := make([][]bool, size)
	visitedCoords := make([][2]int, 0)
	
	for i := range visited {
		visited[i] = make([]bool, size)
	}

	// Mark positions of existing moves on the board as visited
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if board[x][y] != EMPTY {
				visited[x][y] = true
				visitedCoords = append(visitedCoords, [2]int{x, y})
			}
		}
	}

	if len(visitedCoords) >= 10 && once {
		once = false
		fmt.Println("Visited Coords qwertyuio:")
		fmt.Println(visitedCoords)
		fmt.Println("Board qwertyuio:")
		printBoardWithIndexing2(board)

	}


	if len(visitedCoords) >= 12 && twice {
		twice = false
		fmt.Println("Visited Coords qwertyuio:")
		fmt.Println(visitedCoords)
		fmt.Println("Board qwertyuio:")
		printBoardWithIndexing2(board)

	}
	// Generate possible moves (empty cells around existing pieces)
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if board[x][y] != EMPTY {
				for dx := -2; dx <= 2; dx++ {
					for dy := -2; dy <= 2; dy++ {
						nx, ny := x+dx, y+dy
							if nx >= 0 && nx < size && ny >= 0 && ny < size &&
								!visited[nx][ny] && board[nx][ny] == EMPTY {

								visited[nx][ny] = true
							moves = append(moves, [2]int{nx, ny})
						}
					}
				}
			}
		}
	}

	// If no moves found, consider any empty cell
	if len(moves) == 0 {
		for x := 0; x < size; x++ {
			for y := 0; y < size; y++ {
				if board[x][y] == EMPTY && !visited[x][y] {
					moves = append(moves, [2]int{y, x})  // Swap x and y herex][y] {
				}
			}}
	}

	return moves, visitedCoords
}

func isTerminal(board [][]int, player int, opponent int) bool {
	return hasWon(board, player) || hasWon(board, opponent) || isBoardFull(board)
}

func hasWon(board [][]int, player int) bool {
	size := len(board)
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if board[x][y] == player {
				if (y <= size-5 && checkLine(board, x, y, 0, 1, player)) ||
					(x <= size-5 && checkLine(board, x, y, 1, 0, player)) ||
					(x <= size-5 && y <= size-5 && checkLine(board, x, y, 1, 1, player)) ||
					(x <= size-5 && y >= 4 && checkLine(board, x, y, 1, -1, player)) {
					return true
				}
			}
		}
	}
	return false
}

func checkLine(board [][]int, x, y, dx, dy, player int) bool {
	for i := 0; i < 5; i++ {
		nx := x + dx*i
		ny := y + dy*i
		if nx < 0 || nx >= len(board) || ny < 0 || ny >= len(board) || board[nx][ny] != player {
			return false
		}
	}
	return true
}

func isBoardFull(board [][]int) bool {
	for _, row := range board {
		for _, cell := range row {
			if cell == EMPTY {
				return false
			}
		}
	}
	return true
}

func makeCopy(board [][]int) [][]int {
	newBoard := make([][]int, len(board))
	for i := range board {
		newBoard[i] = make([]int, len(board[i]))
		copy(newBoard[i], board[i])
	}

	return newBoard
}

func toPlayer(color string) int {
	if color == "black" {
		return BLACK
	}
	return WHITE
}

func getOpponent(player int) int {
	if player == BLACK {
		return WHITE
	}
	return BLACK
}
