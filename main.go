package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
	"strconv"
	"strings"	
)

func readCsvFile(filePath string) [][]string {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal("Unable to parse file as CSV for " + filePath, err)
    }

    return records
}

type Cleaner struct {
	name string
	model string
	location_x int
	location_y int
	battary int
	movement_energy int
	vacuum_energy int
	dirt_volume int
}

func (c *Cleaner) move_left(room [][]string) {
	if c.location_x > 0 && c.battary >= c.movement_energy {
		nextTile := strings.TrimSpace(room[c.location_y][c.location_x-1])
		if nextTile == "9001" {
			fmt.Println("Cannot move left, there is a wall")
			return
		}
		c.location_x -= 1
		c.battary -= c.movement_energy
	} else {
		fmt.Println("You can't move left or not enough battery")
	}
}

func (c *Cleaner) move_right(room [][]string) {
	if c.location_x < len(room[0])-1 && c.battary >= c.movement_energy {
		nextTile := strings.TrimSpace(room[c.location_y][c.location_x+1])
		if nextTile == "9001" {
			fmt.Println("Cannot move right, there is a wall")
			return
		}
		c.location_x += 1
		c.battary -= c.movement_energy
	} else {
		fmt.Println("You can't move right or not enough battery")
	}
}

func (c *Cleaner) move_up(room [][]string) {
	if c.location_y > 0 && c.battary >= c.movement_energy {
		nextTile := strings.TrimSpace(room[c.location_y-1][c.location_x])
		if nextTile == "9001" {
			fmt.Println("Cannot move up, there is a wall")
			return
		}
		c.location_y -= 1
		c.battary -= c.movement_energy
	} else {
		fmt.Println("You can't move up or not enough battery")
	}
}

func (c *Cleaner) move_down(room [][]string) {
	if c.location_y < len(room)-1 && c.battary >= c.movement_energy {
		nextTile := strings.TrimSpace(room[c.location_y+1][c.location_x])
		if nextTile == "9001" {
			fmt.Println("Cannot move down, there is a wall")
			return
		}
		c.location_y += 1
		c.battary -= c.movement_energy
	} else {
		fmt.Println("You can't move down or not enough battery")
	}
}

func (c *Cleaner) clean(room [][]string) {
	if c.battary >= c.vacuum_energy {
		c.battary -= c.vacuum_energy
		tileValue := room[c.location_y][c.location_x]
		tileValue = strings.TrimSpace(tileValue)
		tileValueInt, err := strconv.Atoi(tileValue)
		if err != nil {
			fmt.Println("Error converting tile value to int:", err)
			return
		}
		c.dirt_volume += tileValueInt
		fmt.Println("Cleaning tile with value:", tileValue)
	} else {
		fmt.Println("You don't have enough battery")
	}
}

func main() {
    room_data := readCsvFile("C:\\Users\\37129\\Intro_to_Ai\\Intro_to_AI_hw\\room.csv")
    fmt.Println(room_data)
	// get the room width and height
	roomWidth := len(room_data[0])
	roomHeight := len(room_data)
	// Create a new cleaner
	cleaner := Cleaner{name: "Rumba", model:"Anna", location_x: 0, location_y: 0, battary: 100, movement_energy: 3, vacuum_energy: 5}

	fmt.Println("Room Width: ", roomWidth)
	fmt.Println("Room Height: ", roomHeight)
	fmt.Println("Cleaner Name: ", cleaner.name)
	fmt.Println("Cleaner Model: ", cleaner.model)

	fmt.Println("Moving down")
	cleaner.move_down(room_data)

	cleaner.move_right(room_data)

	cleaner.clean(room_data)

	fmt.Println("Dirt value: ", cleaner.dirt_volume)
}