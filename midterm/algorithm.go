package main

import (
	"fmt"
	"math"
)

const (
	EMPTY = 0
	BLACK = 1
	WHITE = 2

	// Pattern scores
	FIVE_IN_A_ROW = 1000000
	OPEN_FOUR     = 50000
	CLOSED_FOUR   = 10000
	OPEN_THREE    = 5000
	CLOSED_THREE  = 1000
	OPEN_TWO      = 100
	CLOSED_TWO    = 10

)

// Direction vectors for checking lines
var DIRECTIONS = [][]int{
	{1, 0},  // Horizontal
	{0, 1},  // Vertical
	{1, 1},  // Diagonal
	{1, -1}, // Anti-diagonal
}

// For initial move on an empty board, use the center position
func firstMove(board [][]int) [2]int {
	size := len(board)
	return [2]int{size / 2, size / 2} // Center of board
}

func findBestMove(board [][]int, color string, maxDepth int) [2]int {
	// Check if this is the first move
	isEmpty := true
	for i := range board {
		for j := range board[i] {
			if board[i][j] != EMPTY {
				isEmpty = false
				break
			}
		}
		if !isEmpty {
			break
		}
	}

	if isEmpty {
		return firstMove(board)
	}

	player := toPlayer(color)
	opponent := getOpponent(player)
	bestScore := math.MinInt32
	var bestMove [2]int

	// Generate possible moves
	moves := generateMoves(board)

	// If no moves found (shouldn't happen)
	if len(moves) == 0 {
		return firstMove(board)
	}

	// Check for immediate winning move
	for _, move := range moves {
		row, col := move[0], move[1]
		newBoard := makeCopy(board)
		newBoard[row][col] = player
		if hasWon(newBoard, player) {
			printBoardClearly(newBoard)
			fmt.Println("Found winning move:", row, col)
			return [2]int{row, col}
		}
	}

	// Enhanced check for opponent's immediate winning moves with prioritization
	// First, collect all possible blocking moves
	blockingMoves := make([][2]int, 0)
	blockingScores := make([]int, 0) // Higher score means more critical to block

	for _, move := range moves {
		row, col := move[0], move[1]
		newBoard := makeCopy(board)
		newBoard[row][col] = opponent
		if hasWon(newBoard, opponent) {
			blockingMoves = append(blockingMoves, [2]int{row, col})
			// Calculate the "centrality" of this move in the winning line
			// This helps choose the most effective blocking position
			score := calculateBlockingScore(board, row, col, opponent)
			blockingScores = append(blockingScores, score)
		}
	}

	// If we found blocking moves, select the one with the highest priority
	if len(blockingMoves) > 0 {
		bestIndex := 0
		for i := 1; i < len(blockingScores); i++ {
			if blockingScores[i] > blockingScores[bestIndex] {
				bestIndex = i
			}
		}

		bestMove := blockingMoves[bestIndex]

		// Print the board state and the blocking move
		printBoardClearly(board)

		fmt.Println("Blocking opponent's winning move:", bestMove[0], bestMove[1])
		// revert to the original best move
		tempMove := bestMove[0]
		bestMove[0] = bestMove[1]
		bestMove[1] = tempMove

		return bestMove
	}

	// Use minimax with alpha-beta pruning for other moves
	for _, move := range moves {
		row, col := move[0], move[1]
		newBoard := makeCopy(board)
		newBoard[row][col] = player
		score := minimax(newBoard, maxDepth, math.MinInt32, math.MaxInt32, false, player, opponent)

		if score > bestScore {
			bestScore = score
			bestMove = [2]int{row, col}
		}
	}

	return bestMove
}

var once = true
var twice = true

func generateMoves(board [][]int) [][2]int {
	size := len(board)
	moves := make([][2]int, 0)
	visited := make([][]bool, size)
	visitedCoords := make([][2]int, 0)
	correct_board := give_correct_board()

	for i := range visited {
		visited[i] = make([]bool, size)
	}

	// Reset visited array
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			visited[x][y] = false
		}
	}

	// Gather visitedCoords using correct indexing
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if correct_board[y][x] != EMPTY {
				visited[x][y] = true
				visitedCoords = append(visitedCoords, [2]int{x, y})
			}
		}
	}

	if len(visitedCoords) > 10 && once {
		fmt.Println("Correct Board: qerqwrqewrewq ")
		printBoardWithIndexing(correct_board)
		once = false
		fmt.Println("Visited Coords:")
		fmt.Println(visitedCoords)
	}

	// Generate moves around existing pieces
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

	// If no moves found...
	if len(moves) == 0 {
		for x := 0; x < size; x++ {
			for y := 0; y < size; y++ {
				if board[x][y] == EMPTY && !visited[x][y] {
					moves = append(moves, [2]int{y, x})
				}
			}
		}
	}

	return moves
}

func minimax(board [][]int, depth int, alpha int, beta int, maximizingPlayer bool, player int, opponent int) int {
	// Terminal conditions
	if hasWon(board, player) {
		return FIVE_IN_A_ROW
	}
	if hasWon(board, opponent) {
		return -FIVE_IN_A_ROW
	}
	if isBoardFull(board) || depth == 0 {
		return evaluate(board, player, opponent)
	}

	moves := generateMoves(board)

	if maximizingPlayer {
		maxVal := math.MinInt32
		for _, move := range moves {
			x, y := move[0], move[1]
			if board[x][y] != EMPTY {
				continue
			}

			newBoard := makeCopy(board)
			newBoard[x][y] = player

			val := minimax(newBoard, depth-1, alpha, beta, false, player, opponent)
			maxVal = max(maxVal, val)
			alpha = max(alpha, val)

			if beta <= alpha {
				break // Beta cutoff
			}
		}
		return maxVal
	} else {
		minVal := math.MaxInt32
		for _, move := range moves {
			x, y := move[0], move[1]
			if board[x][y] != EMPTY {
				continue
			}

			newBoard := makeCopy(board)
			newBoard[x][y] = opponent

			val := minimax(newBoard, depth-1, alpha, beta, true, player, opponent)
			minVal = min(minVal, val)
			beta = min(beta, val)

			if beta <= alpha {
				break // Alpha cutoff
			}
		}
		return minVal
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func evaluate(board [][]int, player int, opponent int) int {
	// Evaluate board position for both players and return difference
	playerScore := evaluateForPlayer(board, player)
	opponentScore := evaluateForPlayer(board, opponent)
	return playerScore - opponentScore
}

func evaluateForPlayer(board [][]int, player int) int {
	size := len(board)
	score := 0

	// Check all possible lines
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if board[x][y] != player {
				continue
			}

			// Check all directions
			for _, dir := range DIRECTIONS {
				dx, dy := dir[0], dir[1]

				// Count consecutive stones and open ends
				count := 1
				openEnds := 0

				// Check backward for open end
				prevX, prevY := x-dx, y-dy
				if prevX >= 0 && prevX < size && prevY >= 0 && prevY < size {
					if board[prevX][prevY] == EMPTY {
						openEnds++
					}
				}

				// Count forward streak
				for i := 1; i < 5; i++ {
					nx, ny := x+dx*i, y+dy*i
					if nx < 0 || nx >= size || ny < 0 || ny >= size {
						break
					}
					if board[nx][ny] == player {
						count++
					} else if board[nx][ny] == EMPTY {
						openEnds++
						break
					} else {
						break
					}
				}

				// Score based on pattern
				score += getPatternScore(count, openEnds)
			}
		}
	}

	return score
}

func getPatternScore(count, openEnds int) int {
	if count >= 5 {
		return FIVE_IN_A_ROW // Win
	}

	switch count {
	case 4:
		if openEnds >= 1 {
			return OPEN_FOUR
		}
		return CLOSED_FOUR
	case 3:
		if openEnds == 2 {
			return OPEN_THREE
		}
		if openEnds == 1 {
			return CLOSED_THREE
		}
	case 2:
		if openEnds == 2 {
			return OPEN_TWO
		}
		if openEnds == 1 {
			return CLOSED_TWO
		}
	}
	return 0
}

func hasWon(board [][]int, player int) bool {
	size := len(board)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if board[x][y] == player {
				// Check all directions
				for _, dir := range DIRECTIONS {
					dx, dy := dir[0], dir[1]

					// Check if 5 in a row
					count := 1
					for i := 1; i < 5; i++ {
						nx, ny := x+dx*i, y+dy*i
						if nx < 0 || nx >= size || ny < 0 || ny >= size ||
							board[nx][ny] != player {
							break
						}
						count++
					}

					if count >= 5 {
						return true
					}
				}
			}
		}
	}

	return false
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

// Function to print the board in a more readable format
func printBoardClearly(board [][]int) {
	size := len(board)

	// Print column indices
	fmt.Print("   ")
	for i := 0; i < size; i++ {
		fmt.Printf("%2d", i)
	}
	fmt.Println()

	// Print top border
	fmt.Print("  +")
	for i := 0; i < size; i++ {
		fmt.Print("--")
	}
	fmt.Println("+")

	// Print board with row indices
	for i := 0; i < size; i++ {
		fmt.Printf("%2d |", i)
		for j := 0; j < size; j++ {
			switch board[i][j] {
			case EMPTY:
				fmt.Print(" .")
			case BLACK:
				fmt.Print(" B")
			case WHITE:
				fmt.Print(" W")
			}
		}
		fmt.Println(" |")
	}

	// Print bottom border
	fmt.Print("  +")
	for i := 0; i < size; i++ {
		fmt.Print("--")
	}
	fmt.Println("+")
}

// Add this new function to calculate the blocking priority score
func calculateBlockingScore(board [][]int, x, y, player int) int {
	size := len(board)
	score := 0

	// Check all directions for ongoing patterns
	for _, dir := range DIRECTIONS {
		dx, dy := dir[0], dir[1]

		// Check how central this position is to a potential winning line
		// Look in both directions
		totalPieces := 1 // Count the current position
		maxLength := 1

		// Check forward
		for i := 1; i < 5; i++ {
			nx, ny := x+dx*i, y+dy*i
			if nx < 0 || nx >= size || ny < 0 || ny >= size || board[nx][ny] != player {
				break
			}
			totalPieces++
			maxLength++
		}

		// Check backward
		for i := 1; i < 5; i++ {
			nx, ny := x-dx*i, y-dy*i
			if nx < 0 || nx >= size || ny < 0 || ny >= size || board[nx][ny] != player {
				break
			}
			totalPieces++
			maxLength++
		}

		// Score based on both the length and whether this move blocks the middle of a pattern
		// Higher score for more central blocking positions and longer patterns
		if maxLength >= 4 {
			score += maxLength * 100

			// If this position is in the middle of the pattern, it's a better block
			if totalPieces > 1 {
				middlePriority := (maxLength - abs(maxLength/2-totalPieces/2)) * 50
				score += middlePriority
			}
		}
	}

	return score
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
