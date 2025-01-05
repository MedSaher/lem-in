package main

import (
	"fmt"
	"os"

	"github.com/Douirat/lem-in/logic"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Not enough arguments!")
		return
	}
	colony := logic.NewColony()
	err := colony.RockAndRoll(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	colony.JoinPaths(colony.FindAllPathsDFS(colony.Start, colony.End))
	colony.RoomMap()

	colony.CalcInterWeight()
	colony.SortByIntersectionWeight()
	uselessPaths := make(map[*logic.Path]bool)
	visitedRooms := make(map[*logic.Room]bool)
	cleanPaths := &logic.Paths{}
	for _, path := range colony.PathsSet.Paths {
		if uselessPaths[path] {

			continue
		}
		cleanPaths.Paths = append(cleanPaths.Paths, path)
		for _, room := range path.Rooms {
			if room.Name == colony.Start.Name || room.Name == colony.End.Name || visitedRooms[room] {
				continue
			}
			for _, uslesspaths := range colony.PathsSet.RoomMap[room] {
				uselessPaths[uslesspaths] = true
				visitedRooms[room] = true
			}
		}
	}
	colony.PathsSet.Paths = cleanPaths.Paths
	colony.SortByLength()
	// for _, path := range colony.PathsSet.Paths{
	// 	for _, room := range path.Rooms{

	// 		fmt.Printf("%v -> ", room.Name)
	// 	}
	// 	fmt.Println()
	// }
	colony.EditNextRoom()
	colony.ExtractAntNum()
	// for _, ant := range colony.Ants{
	// 	fmt.Printf("ant name : %v; current position : %v\n", ant.Name, ant.CurrentPosition.Name)
	// }
	colony.Start.Quantity = colony.AntNum
	// fmt.Println(colony.Start.Quantity)
	// Ensure all ants move only once per turn
// Process turns until all ants reach the end room
for colony.End.Quantity < colony.AntNum {
    // Store moves for the current turn
    firstSlice := []string{}
    madeMove := false // Track if any ant moves during this turn

    for i:=0; i < len(colony.Ants); i++{
		ant := colony.Ants[i]
        // If the ant is already at the end room, skip it
        if ant.CurrentPosition.Name == colony.End.Name {
            continue
        }

        // If the ant is at the start room, try to send it to the first room in its path

			if ant.CurrentPosition.Name == colony.Start.Name {
				move := colony.SendFromStart(ant)
				if move != "" {
					firstSlice = append(firstSlice, move)
					madeMove = true
				}
				continue
			}

        // Move the ant along its path if the next room is not occupied
        nextRoom := ant.CurrentPosition.Next
        if nextRoom != nil && !nextRoom.Occupied {
            // Mark the current room as unoccupied
            ant.CurrentPosition.Occupied = false

            // Move the ant to the next room
            ant.CurrentPosition = nextRoom
            if ant.CurrentPosition.Name != colony.End.Name {
                ant.CurrentPosition.Occupied = true
            }

            // Log the move
            firstSlice = append(firstSlice, fmt.Sprintf(logic.MovesPattern, ant.Name, ant.CurrentPosition.Name))
            madeMove = true

            // If the ant reaches the end room, update End.Quantity
            if nextRoom.Name == colony.End.Name {
                colony.End.Quantity++
            }
        }
    }

    // Add all moves from this turn to the global moves log
    if len(firstSlice) > 0 {
        logic.Moves = append(logic.Moves, firstSlice)
    }

    // If no moves were made and ants are still not at the end, exit to prevent infinite loop
    if !madeMove {
        fmt.Println("No moves possible this turn. Exiting to avoid infinite loop.")
        break
    }
}

	for _, move := range logic.Moves{
		for _, step := range move{
			fmt.Printf("%v ", step)
		}
		fmt.Println()
	}
	// fmt.Println(colony.Start.Next.Name)
	// logic.LogAntsMovement(colony)
	// for _, path := range colony.PathsSet.Paths{
	// 	for _, room := range path.Rooms{
	// 		fmt.Printf("%v : %v\n", room.Name, room.Occupied)
	// 	}
	// }
}

