package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	WorldWidth            = 20
	WorldHeight           = 10
	InitialSpawnTolerance = 0.4
	MaxSteps              = 100
)

// Clear the terminal.
func clearConsole() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
}

// Check if a cell falls within grid bounds.
func inbounds(coords *Coord) bool {
	row := coords.Row
	col := coords.Col
	return row >= 0 && row < WorldHeight && col >= 0 && col < WorldWidth
}

// Find all in-bounds neighbors of a given 2D coordinate.
func neighbors(coords *Coord) []string {
	row := coords.Row
	col := coords.Col

	neighbors := make([]string, 0)

	triplet := [3]int{-1, 0, 1}
	for _, rowDelta := range triplet {
		for _, colDelta := range triplet {
			if rowDelta == 0 && colDelta == 0 {
				// this is the cell itself
				continue
			}

			nRow := row + rowDelta
			nCol := col + colDelta
			nCoords := &Coord{nRow, nCol}
			if inbounds(nCoords) {
				neighbors = append(neighbors, nCoords.String())
			}
		}
	}

	return neighbors
}

// Coord describes a 2D position in a grid.
type Coord struct {
	Row, Col int
}

// Initialize Coord.
func NewCoord(row, col int) *Coord {
	return &Coord{row, col}
}

// Stringify a Coord.
func (c *Coord) String() string {
	return fmt.Sprintf("%d,%d", c.Row, c.Col)
}

// Given a string of coordinates, build a new Coord.
func CoordFromString(coordStr string) *Coord {
	parts := strings.Split(coordStr, ",")

	row, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}
	col, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}

	return &Coord{row, col}
}

// Cell represents a single living entity in the grid.
type Cell struct {
	Alive bool
}

// Initialize Cell.
func NewCell() *Cell {
	return &Cell{false}
}

// Set the cell's alive state.
func (c *Cell) SetState(state bool) {
	c.Alive = state
}

// Spawn a new cell.
func (c *Cell) Spawn() {
	c.Alive = true
}

// Kill off the cell.
func (c *Cell) Kill() {
	c.Alive = false
}

// Make a copy of the cell.
func (c *Cell) Copy() *Cell {
	return &Cell{c.Alive}
}

// World is a 2D grid filled with Cells.
type World struct {
	Cells map[string]*Cell
}

// Build a new world with cells in a random state.
func NewWorld() *World {
	cells := make(map[string]*Cell)
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	for row := 0; row < WorldHeight; row++ {
		for col := 0; col < WorldWidth; col++ {
			coords := NewCoord(row, col)
			cell := NewCell()
			cells[coords.String()] = cell
			cell.SetState(InitialSpawnTolerance >= randGen.Float64())
		}
	}

	return &World{cells}
}

// Apply automata rules to all cells in the grid.
func (w *World) Step() {
	// Make a new map containing the future state of the world. Game of Life rules are based on
	// current timestep. We will use this to maintain next state.
	nextState := make(map[string]*Cell)

	for coordStr, _ := range w.Cells {
		pastCell := w.Cells[coordStr]
		livingNeighbors := 0
		coords := CoordFromString(coordStr)
		nextCell := pastCell.Copy()

		// Grab the number of living cells surrounding the current cell.
		for _, neighborCoords := range neighbors(coords) {
			neighbor, ok := w.Cells[neighborCoords]
			if ok && neighbor.Alive {
				livingNeighbors++
			}
		}

		// Apply Conway's rules.
		if pastCell.Alive {
			if livingNeighbors < 2 {
				// Kill due to nderpopulation
				nextCell.Kill()
			} else if livingNeighbors > 3 {
				// Kill due to overpopulation
				nextCell.Kill()
			}
		} else {
			if livingNeighbors == 3 {
				// Reproduce
				nextCell.Spawn()
			}
		}
		nextState[coordStr] = nextCell
	}
	w.Cells = nextState
}

// Print the current state of the world to terminal.
func (w *World) Draw() {
	var buf bytes.Buffer
	for row := 0; row < WorldHeight; row++ {
		for col := 0; col < WorldWidth; col++ {
			coords := Coord{row, col}
			cell, ok := w.Cells[coords.String()]

			if !ok {
				continue
			}

			buf.WriteByte(' ')
			b := byte('.')
			if cell.Alive {
				b = '0'
			}
			buf.WriteByte(b)
			buf.WriteByte(' ')
		}
		buf.WriteByte('\n')
	}
	fmt.Print(buf.String())
}

// Play Conway's Game of Life.
func main() {
	w := NewWorld()
	for i := 0; i < MaxSteps; i++ {
		clearConsole()
		w.Draw()
		w.Step()
		time.Sleep(time.Second / 2)
	}
}
