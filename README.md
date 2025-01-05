# Ant Colony Pathfinding Project

## Overview
This project simulates an ant colony navigating through a graph of rooms to reach the end room using optimized pathfinding. The program ensures efficient movement of ants by prioritizing paths and handling edge cases such as bottlenecks or direct paths from the start to the end room.

## Features
- **Efficient Pathfinding**: Ants follow the shortest path when appropriate and avoid unnecessary moves.
- **Multiple Path Handling**: Supports multiple paths and ensures ants distribute themselves efficiently.
- **Edge Case Handling**: Includes logic for cases with no intermediate rooms or when the shortest path is blocked.

## How It Works
1. **Input**: The program reads an input file that defines the graph, including rooms, paths, and the number of ants.
2. **Initialization**:
   - Rooms and paths are parsed.
   - Ants are assigned to the start room.
3. **Simulation**:
   - Ants move along paths toward the end room.
   - The shortest path is prioritized for the last ant in the start room.
4. **Output**: The movements of ants are printed turn by turn until all ants reach the end room.

## Key Logic
### Pathfinding Algorithm
- **Shortest Path Priority**:
  - Identifies the shortest path and reserves it for the last ant.
  - Other ants distribute across available paths to minimize congestion.
- **Room Occupancy**:
  - Tracks which ants occupy which rooms.
  - Prevents multiple ants from entering the same room in a turn.

### Edge Case Handling
- **Direct Start-to-End Path**:
  - If the shortest path connects the start and end directly, ants use it efficiently.
- **Blocked Paths**:
  - Ants wait if their preferred path is blocked, ensuring fairness and preventing deadlocks.

## Input Format
The input file should follow this structure:
- **Number of Ants**: The first line contains the number of ants.
- **Rooms**: Each subsequent line defines a room.
- **Paths**: Defines connections between rooms.
- Example:
  ```
  4
  Start 0
  Room1 1
  End 2
  Start-Room1
  Room1-End
  ```

## Output Format
The output displays the movements of ants in each turn:
```
L1-Room1 L2-End
L3-Room1
```
- `L1-Room1`: Ant 1 moves to Room1.
- `L2-End`: Ant 2 reaches the end room.

## How to Run
1. Clone the repository:
   ```bash
   git clone <repository_url>
   cd <repository_name>
   ```
2. Run the program with an input file:
   ```bash
   go run main/main.go <input_file>
   ```
3. Example:
   ```bash
   go run main/main.go example00.txt
   ```

## Code Highlights
### Handling Start Room
```go
if len(colony.Ants) == 1 && path != shortPath {
    continue
}
```
- Ensures the last ant takes the shortest path.

### Moving Ants
```go
if nextRoom != nil && !nextRoom.Occupied {
    ant.CurrentPosition = nextRoom
    if nextRoom.Name == colony.End.Name {
        colony.End.Quantity++
    }
}
```
- Moves ants forward and updates room occupancy.

## Contributing
1. Fork the repository.
2. Create a feature branch:
   ```bash
   git checkout -b feature/new-feature
   ```
3. Commit changes and push:
   ```bash
   git commit -m "Add new feature"
   git push origin feature/new-feature
   ```
4. Create a pull request.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

---

Happy Coding!

