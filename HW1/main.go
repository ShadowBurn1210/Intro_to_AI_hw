package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func (c *Cleaner) readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(f)

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	csvReader.Comment = '#' // Ignore lines starting with #

	for i := 0; i < 5; i++ {
		initialState, err := csvReader.Read()
		if err != nil {
			fmt.Print("Error reading csv data", err)
		}

		value, err := strconv.Atoi(strings.TrimSpace(initialState[0]))
		if err != nil {
			fmt.Println("Error converting initial state value to int:", err)
			continue
		}
		fmt.Println("Value:", value)
		if i == 0 {
			c.locationX = value
		}
		if i == 1 {
			c.locationY = value
		}
		if i == 2 {
			c.battery = value
		}
		if i == 3 {
			c.movementEnergy = value
		}
		if i == 4 {
			c.vacuumEnergy = value
		}

	}

	// Read all the remaining records
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Error reading csv data", err)
	}
	return records
}

type Cleaner struct {
	name           string
	model          string
	locationX      int
	locationY      int
	battery        int
	movementEnergy int
	vacuumEnergy   int
	dirtVolume     int
	tilesCleaned   int
}

func (c *Cleaner) feedback(path []string) {
	fmt.Println("Battery:", c.battery)
	fmt.Println("Dirt volume:", c.dirtVolume)
	fmt.Println("Path:", path)
	fmt.Println("Tiles cleaned:", c.tilesCleaned)
}
func (c *Cleaner) moveLeft(room [][]string) {
	if c.locationX > 0 && c.battery >= c.movementEnergy {
		nextTile := strings.TrimSpace(room[c.locationY][c.locationX-1])
		if nextTile == "9001" {
			fmt.Println("Cannot move left, there is a wall")
			return
		}
		c.locationX -= 1
		c.battery -= c.movementEnergy
	} else {
		fmt.Println("You can't move left or not enough battery")
	}
}

func (c *Cleaner) moveRight(room [][]string) {
	if c.locationX < len(room[0])-1 && c.battery >= c.movementEnergy {
		nextTile := strings.TrimSpace(room[c.locationY][c.locationX+1])
		if nextTile == "9001" {
			fmt.Println("Cannot move right, there is a wall")
			return
		}
		c.locationX += 1
		c.battery -= c.movementEnergy
	} else {
		fmt.Println("You can't move right or not enough battery")
	}
}

func (c *Cleaner) moveUp(room [][]string) {
	if c.locationY > 0 && c.battery >= c.movementEnergy {
		nextTile := strings.TrimSpace(room[c.locationY-1][c.locationX])
		if nextTile == "9001" {
			fmt.Println("Cannot move up, there is a wall")
			return
		}
		c.locationY -= 1
		c.battery -= c.movementEnergy
	} else {
		fmt.Println("You can't move up or not enough battery")
	}
}

func (c *Cleaner) moveDown(room [][]string) {
	if c.locationY < len(room)-1 && c.battery >= c.movementEnergy {
		nextTile := strings.TrimSpace(room[c.locationY+1][c.locationX])
		if nextTile == "9001" {
			fmt.Println("Cannot move down, there is a wall")
			return
		}
		c.locationY += 1
		c.battery -= c.movementEnergy
	} else {
		fmt.Println("You can't move down or not enough battery")
	}
}

func (c *Cleaner) clean(room [][]string) {
	if c.battery >= c.vacuumEnergy {
		c.battery -= c.vacuumEnergy
		tileValue := room[c.locationY][c.locationX]
		tileValue = strings.TrimSpace(tileValue)
		tileValueInt, err := strconv.Atoi(tileValue)
		if err != nil {
			fmt.Println("Error converting tile value to int:", err)
			return
		}
		c.dirtVolume += tileValueInt
		c.tilesCleaned += 1
		fmt.Println("Cleaning tile with value:", tileValue)
	} else {
		fmt.Println("You don't have enough battery")
	}
}

// AStar Whole a* algorithm was implemented with large help of Deep Seek R1 model, which does not provide link for chat reference
// This was my promt : "Can you edit this Astar algorithm, so that it finds shortest path to the dirtiest node, but if there is a node with value 9001 it knows it is a wall"
// Which used initialy a* algorithm from internet site which I can't find anymore
func AStar(startX, startY int, room [][]string) []string {
	type Node struct {
		x, y, g, h, f int
		parent        *Node
	}

	// Preprocess to find the dirtiest nodes (max non-wall value)
	maxDirt := -1
	var dirtiestNodes []struct{ x, y int }
	for y, row := range room {
		for x := range row {
			valStr := strings.TrimSpace(room[y][x])
			if valStr == "9001" {
				continue // Skip walls
			}
			val, err := strconv.Atoi(valStr)
			if err != nil || val <= 0 {
				continue // Skip invalid or clean nodes
			}
			if val > maxDirt {
				maxDirt = val
				dirtiestNodes = []struct{ x, y int }{{x, y}}
			} else if val == maxDirt {
				dirtiestNodes = append(dirtiestNodes, struct{ x, y int }{x, y})
			}
		}
	}

	if len(dirtiestNodes) == 0 {
		return nil // No dirty nodes to clean
	}

	// Create quick lookup map for dirtiest nodes
	dirtiestMap := make(map[string]bool)
	for _, node := range dirtiestNodes {
		dirtiestMap[fmt.Sprintf("%d,%d", node.x, node.y)] = true
	}

	// Heuristic: minimum Manhattan distance to any dirtiest node
	heuristic := func(x, y int) int {
		minDist := math.MaxInt32
		for _, dn := range dirtiestNodes {
			dist := abs(x-dn.x) + abs(y-dn.y)
			if dist < minDist {
				minDist = dist
			}
		}
		return minDist
	}

	// Get valid neighbors (non-wall nodes)
	neighbors := func(node *Node) []*Node {
		directions := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		var result []*Node
		for _, d := range directions {
			nx, ny := node.x+d[0], node.y+d[1]
			if nx >= 0 && ny >= 0 && nx < len(room[0]) && ny < len(room) {
				valStr := strings.TrimSpace(room[ny][nx])
				if valStr != "9001" { // Skip walls
					result = append(result, &Node{x: nx, y: ny})
				}
			}
		}
		return result
	}

	// Initialize open set with start node
	openSet := []*Node{{
		x: startX,
		y: startY,
		g: 0,
		h: heuristic(startX, startY),
		f: heuristic(startX, startY),
	}}
	closedSet := make(map[string]bool)
	nodeMap := make(map[string]*Node)

	for len(openSet) > 0 {
		// Find node with lowest f-cost
		currentIndex := 0
		current := openSet[0]
		for i, node := range openSet {
			if node.f < current.f {
				current = node
				currentIndex = i
			}
		}

		// Check if we've reached the dirtiest node
		if dirtiestMap[fmt.Sprintf("%d,%d", current.x, current.y)] {
			var path []string
			for current != nil {
				path = append([]string{fmt.Sprintf("(%d,%d)", current.x, current.y)}, path...)
				current = current.parent
			}
			return path
		}

		// Move current node to closed set
		openSet = append(openSet[:currentIndex], openSet[currentIndex+1:]...)
		closedSet[fmt.Sprintf("%d,%d", current.x, current.y)] = true

		// Process neighbors
		for _, neighbor := range neighbors(current) {
			neighborKey := fmt.Sprintf("%d,%d", neighbor.x, neighbor.y)
			if closedSet[neighborKey] {
				continue
			}

			tentativeG := current.g + 1
			existing, exists := nodeMap[neighborKey]

			if !exists || tentativeG < existing.g {
				neighbor.g = tentativeG
				neighbor.h = heuristic(neighbor.x, neighbor.y)
				neighbor.f = neighbor.g + neighbor.h
				neighbor.parent = current

				if !exists {
					openSet = append(openSet, neighbor)
					nodeMap[neighborKey] = neighbor
				} else {
					existing.g = neighbor.g
					existing.h = neighbor.h
					existing.f = neighbor.f
					existing.parent = neighbor.parent
				}
			}
		}
	}

	return nil // No path found
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
// moveSomewhere is a function that moves the cleaner to the next node, based on the path calculated by the A* algorithm
func (c *Cleaner) moveSomewhere(nextNode string, roomData [][]string) {
	coords := strings.Trim(nextNode, "()")
	parts := strings.Split(coords, ",")
	if len(parts) != 2 {
		fmt.Println("Invalid node format")
		return
	}
	x, err1 := strconv.Atoi(parts[0])
	y, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		fmt.Println("Error converting coordinates:", err1, err2)
		return
	}
	if c.locationX < x {
		c.moveRight(roomData)
	}
	if c.locationX > x {
		c.moveLeft(roomData)
	}
	if c.locationY < y {
		c.moveDown(roomData)
	}
	if c.locationY > y {
		c.moveUp(roomData)
	}

}


// decideToClean is a simple function that decides whether the cleaner should clean the current tile or not, based on if the tiles value is above 0
func (c *Cleaner) decideToClean(roomData [][]string) {
	// Add cleaning logic here
	if tileValue, err := strconv.Atoi(strings.TrimSpace(roomData[c.locationY][c.locationX])); err == nil && tileValue > 0 {
		c.clean(roomData)

		// Update the room data
		roomData[c.locationY][c.locationX] = "0"
	}

}

func main() {

	// Create a new cleaner, which data will be overwritten by the csv file
	cleaner := Cleaner{
		name:           "Rummba",
		model:          "Elizabete",
		locationX:      0,
		locationY:      0,
		battery:        50,
		movementEnergy: 1,
		vacuumEnergy:   5,
		tilesCleaned:   0,
	}

	// Read the csv file and get all the data
	roomData := cleaner.readCsvFile(`C:\Users\37129\Intro_to_Ai\Intro_to_AI_hw\HW1\room.csv`)
	totalPath := []string{}
	fmt.Println(roomData)
	for cleaner.battery > 0 {
		myPath := AStar(cleaner.locationX, cleaner.locationY, roomData) // Start at current location

		if len(myPath) == 0 {
			fmt.Println("No more paths to dirtiest tiles.")
			break
		}

		// Move the cleaner along the path
		for _, node := range myPath {
			if node == fmt.Sprintf("(%d,%d)", cleaner.locationX, cleaner.locationY) {
				continue
			}
			cleaner.moveSomewhere(node, roomData)
			cleaner.decideToClean(roomData)
		}
		// Add the path to the total path
		totalPath = append(totalPath, myPath...)
	}
	cleaner.feedback(totalPath)

}
