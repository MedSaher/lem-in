package logic

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Douirat/lem-in/data"
)

// Declare a structure to represent the room:
type Room struct {
	Name       string
	CorX, CorY int
	Next       *Room
	Quantity   int
	Occupied   bool
}

// declare a structre to represent the graph using the adjacency list:
type Colony struct {
	AntNum     int
	Start, End *Room
	Ants       []*Ant
	Graph      map[string]*Room
	PathsSet   *Paths
}

type Paths struct {
	Paths   []*Path           // List of paths.
	RoomMap map[*Room][]*Path // Maps each room to all paths containing it.
	// BestPaths map[string]*[]*Room // Maps each room to the shortest path that passes through it.
}

type Path struct {
	Rooms              []*Room
	IntersectionWeight int
}

// Declare a structure to represent the ant:
type Ant struct {
	Name            string
	CurrentPosition *Room
}

// Instantiate a colony:
func NewColony() *Colony {
	return &Colony{
		AntNum: 0,
		Start:  nil,
		End:    nil,
		Graph:  make(map[string]*Room),
		PathsSet: &Paths{
			Paths:   []*Path{},
			RoomMap: make(map[*Room][]*Path),
		},
		Ants: []*Ant{},
	}
}

// Instantiate a new room:
func NewRoom(str string) (*Room, error) {
	var err error
	room := new(Room)
	data := strings.Split(str, " ")
	if len(data) != 3 {
		return nil, errors.New("invalid data format")
	}
	room.Name = data[0]
	room.CorX, err = strconv.Atoi(data[1])
	if err != nil {
		return nil, err
	}
	room.CorY, err = strconv.Atoi(data[2])
	if err != nil {
		return nil, err
	}
	room.Next = nil
	room.Occupied = false
	return room, nil
}

// Add a new room to the colony:
func (colony *Colony) AddRoom(str string) (*Room, error) {
	room, err := NewRoom(str)
	if err != nil {
		return nil, err
	}
	colony.Graph[room.Name] = room
	SIZE++

	return room, nil
}

// Add a tunnel between tow rooms:
func (colony *Colony) AddTunnel(str string) error {
	data := strings.Split(str, "-")
	if len(data) != 2 {
		return errors.New("error data format, tunnels not valid")
	}
	roomSrc := &Room{}
	roomSrc.Name = colony.Graph[data[0]].Name
	roomSrc.CorX = colony.Graph[data[0]].CorX
	roomSrc.CorY = colony.Graph[data[0]].CorY
	roomDst := &Room{}
	roomDst.Name = colony.Graph[data[1]].Name
	roomDst.CorX = colony.Graph[data[1]].CorX
	roomDst.CorY = colony.Graph[data[1]].CorY
	roomDst.Next = colony.Graph[data[0]].Next
	if roomDst.Name == colony.Start.Name {
		colony.Start = roomDst
	} else if roomDst.Name == colony.End.Name {
		colony.End = roomDst
	}
	if roomSrc.Name == colony.Start.Name {
		colony.Start = roomSrc
	} else if roomSrc.Name == colony.End.Name {
		colony.End = roomSrc
	}
	colony.Graph[data[0]].Next = roomDst
	roomSrc.Next = colony.Graph[data[1]].Next
	colony.Graph[data[1]].Next = roomSrc
	return nil
}

func (colony *Colony) FindAllPathsDFS(start, end *Room) [][]string {
	result := [][]string{}
	path := []string{}
	visited := make(map[string]bool)

	// Helper function for DFS:
	var dfs func(string)
	dfs = func(current string) {
		// mark the current room as visited:
		visited[current] = true
		// Add current room to the path:
		path = append(path, current)
		// Base case: if the current room is the destination store the vertex and backtrack:
		if current == end.Name {
			result = append(result, append([]string{}, path...))
		} else {
			temp := colony.Graph[current]
			// explore neigbors:
			for temp != nil {
				if !visited[temp.Name] {
					dfs(temp.Name)
				}
				temp = temp.Next
			}
		}
		// backtrack until finding a vertex that has other unexplored neighbors:
		// Remove the current room from the track and marck it as unvisted in case an other path crosses with it:
		path = path[:len(path)-1]
		visited[current] = false
	}
	// Call the dfs:
	dfs(start.Name)
	return result
}

// Create the AntNum based on the required number:
func (colony *Colony) ExtractAntNum() {
	for i := 1; i <= colony.AntNum; i++ {
		Ant := &Ant{}
		Ant.Name = "L" + strconv.Itoa(i)
		Ant.CurrentPosition = colony.Start
		colony.Ants = append(colony.Ants, Ant)
	}
}

// Formulating the colny graph based on the input extracted from the file:
func (colony *Colony) RockAndRoll(fileName string) error {
	started := false
	ended := false
	// end := false
	data, err := data.ReadFile(fileName)
	if err != nil {
		return err
	}
	colony.AntNum, err = strconv.Atoi(data[0])
	if err != nil {
		return err
	}
	data = data[1:]

	// Define regex patterns for each type of line
	patterns := map[string]*regexp.Regexp{
		"start":   regexp.MustCompile(`^##start$`),
		"room":    regexp.MustCompile(`^(.*)\s+(\d+)\s+(\d+)$`),
		"end":     regexp.MustCompile(`^##end$`),
		"tunnel":  regexp.MustCompile(`^([a-zA-Z0-9]+)-([a-zA-Z0-9]+)$`),
		"comment": regexp.MustCompile(`^#.*`),
	}

	// Iterate over the input data
	for _, str := range data {
		// Check each pattern
		for key, rg := range patterns {
			if rg.MatchString(str) {
				// Process based on the matched pattern
				switch key {
				case "start":
					if colony.Start != nil || started {
						return errors.New("wrond data format, the graph is already started")
					}
					started = true
					continue
				case "room":
					// Handle the room pattern
					if started && colony.Start == nil {
						colony.Start, err = colony.AddRoom(str)
						if err != nil {
							return err
						}
						started = false
						continue
					}
					if ended && colony.End == nil {
						colony.End, err = colony.AddRoom(str)
						if err != nil {
							return err
						}
						ended = false
						// end = true
						continue
					}
					_, err = colony.AddRoom(str)
					if err != nil {
						return err
					}
				case "end":
					if colony.End != nil || ended {
						return errors.New("wrond data format, the graph is already ended")
					}
					// Handle the end pattern
					ended = true
					// fmt.Println("End found:", str)
					continue
				case "tunnel":
					if colony.End == nil {
						return errors.New("wrond data format, tunnel before end flag")
					}
					// Handle the tunnel pattern
					colony.AddTunnel(str)
				case "comment":
					// Handle the comment pattern
					continue
				}
			}
		}
	}
	return nil
}

func (colony *Colony) RoomMap() {
	for _, path := range colony.PathsSet.Paths {
		// fmt.Println(len(path.Rooms))
		for _, room := range path.Rooms {
			if room == colony.Start || room == colony.End {
				continue
			}
			colony.PathsSet.RoomMap[room] = append(colony.PathsSet.RoomMap[room], path)
		}
	}
}

func (colony *Colony) JoinPaths(allPaths [][]string) {
	for _, path := range allPaths {
		tempPath := &Path{}
		for _, room := range path {
			roomPath := colony.Graph[room]
			tempPath.Rooms = append(tempPath.Rooms, roomPath)
			tempPath.IntersectionWeight = 0
		}
		colony.PathsSet.Paths = append(colony.PathsSet.Paths, tempPath)

	}
}

// Method to calculate intersection weight
func (colony *Colony) CalcInterWeight() {
	for _, path := range colony.PathsSet.Paths {
		for _, room := range path.Rooms {
			if room.Name == colony.Start.Name || room.Name == colony.End.Name {
				continue
			}
			path.IntersectionWeight += len(colony.PathsSet.RoomMap[room]) - 1
		}
		// fmt.Println(path.IntersectionWeight)
	}
}

// Method to sort paths based on IntersectionWeight
func (colony *Colony) SortByIntersectionWeight() {
	sort.Slice(colony.PathsSet.Paths, func(i, j int) bool {
		return colony.PathsSet.Paths[i].IntersectionWeight < colony.PathsSet.Paths[j].IntersectionWeight
	})
}

// Method to sort paths based on IntersectionWeight
func (colony *Colony) SortByLength() {
	sort.Slice(colony.PathsSet.Paths, func(i, j int) bool {
		return len(colony.PathsSet.Paths[i].Rooms) < len(colony.PathsSet.Paths[j].Rooms)
	})
}

func (colony *Colony) EditNextRoom() {
	for _, path := range colony.PathsSet.Paths {
		for i := 0; i < len(path.Rooms); i++ {
			if path.Rooms[i].Name == colony.End.Name {
				path.Rooms[i].Next = nil
				break
			}
			path.Rooms[i].Next = path.Rooms[i+1]
		}
	}
}

func (colony *Colony) SendFromStart(ant *Ant) string {
	// We want to send ants to available paths, checking for the first available one
	for _, path := range colony.PathsSet.Paths {
		firstRoom := path.Rooms[1] // The first room in the path after the start room
		// If the first room is unoccupied, send the ant there
		if !firstRoom.Occupied {
			ant.CurrentPosition = firstRoom
			firstRoom.Occupied = true
			colony.Start.Quantity--
			return fmt.Sprintf(MovesPattern, ant.Name, firstRoom.Name)
		}
	}

	// If no path is available, return empty
	return ""
}

//	func (colony *Colony) Sendrest(){
//		if colony
//	}
func LogAntsMovement(colony *Colony) {
	fmt.Println("Current State of Ants:")
	for _, ant := range colony.Ants {
		fmt.Printf("Ant %s: %s\n", ant.Name, ant.CurrentPosition.Name)
	}
	fmt.Println("Current State of Rooms:")
	for _, room := range colony.Graph {
		fmt.Printf("Room %s: Occupied: %v\n", room.Name, room.Occupied)
	}
}
