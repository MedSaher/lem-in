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
	colony.Start.Quantity = colony.AntNum
	colony.SendAnts()
	for _, move := range logic.Moves {
		for _, step := range move {
			fmt.Printf("%v ", step)
		}
		fmt.Println()
	}
}
