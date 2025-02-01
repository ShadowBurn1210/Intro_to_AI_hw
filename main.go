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

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
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
}

func (c *Cleaner) feedback() {
	fmt.Println("Location X:", c.locationX)
	fmt.Println("Location Y:", c.locationY)
	fmt.Println("Battery:", c.battery)
	fmt.Println("Dirt volume:", c.dirtVolume)
}
func (c *Cleaner) move_left(room [][]string) {
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

func (c *Cleaner) move_right(room [][]string) {
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

func (c *Cleaner) move_up(room [][]string) {
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

func (c *Cleaner) move_down(room [][]string) {
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
		fmt.Println("Cleaning tile with value:", tileValue)
	} else {
		fmt.Println("You don't have enough battery")
	}
}

// Whole a* algorithm was implemented by Deep Seek R1 model, which does not provide link
func AStar(startX, startY int, room [][]string) []string {
	type Node struct {
		x, y, g, h, f int
		parent        *Node
	}

	// Preprocess to find dirtiest nodes (max non-wall value)
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

		// Check if we've reached a dirtiest node
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
		c.move_right(roomData)
	}
	if c.locationX > x {
		c.move_left(roomData)
	}
	if c.locationY < y {
		c.move_down(roomData)
	}
	if c.locationY > y {
		c.move_up(roomData)
	}

}

func (c *Cleaner) decideToClean(roomData [][]string) {
	// Add cleaning logic here
	if tileValue, err := strconv.Atoi(strings.TrimSpace(roomData[c.locationY][c.locationX])); err == nil && tileValue > 0 {
		c.clean(roomData)
	}

}

func main() {
	roomData := readCsvFile(`C:\Users\Admin\GolandProjects\Intro_to_AI_hw\room.csv`)
	// print room data
	for i := 0; i < len(roomData); i++ {
		fmt.Println(roomData[i])
	}

	// Create a new cleaner
	cleaner := Cleaner{
		name:           "Rumba",
		model:          "Anna",
		locationX:      0,
		locationY:      0,
		battery:        50,
		movementEnergy: 1,
		vacuumEnergy:   5,
	}
	//path := []string{"(0,0)", "(0,1)", "(1,1)", "(2,1)", "(2,2)", "(2,3)", "(1,3)", "(1,4)", "(2,4)", "(3,4)", "(4,4)"}

	//fmt.Println("Path to the dirtiest tile:", path)
	//// Find the path to the dirtiest tile
	myPath := AStar(0, 0, roomData) // Start at (0,0)
	fmt.Println("Path to the dirtiest tile:", myPath)

	// Move the cleaner along the path
	for _, node := range myPath {
		if node == "(0,0)" {
			continue
		}
		cleaner.moveSomewhere(node, roomData)
		cleaner.decideToClean(roomData)
		cleaner.feedback()
	}

}
