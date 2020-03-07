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
)

func clearConsole() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
}

func inbounds(coords *Coord) bool {
	row := coords.Row
	col := coords.Col
	return row >= 0 && row < WorldHeight && col >= 0 && col < WorldWidth
}

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

type Coord struct {
	Row, Col int
}

func (c *Coord) String() string {
	return fmt.Sprintf("%d,%d", c.Row, c.Col)
}

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

type Cell struct {
	Alive bool
}

func NewCell() *Cell {
	return &Cell{false}
}

func (c *Cell) SetState(state bool) {
	c.Alive = state
}

func (c *Cell) Spawn() {
	c.Alive = true
}

func (c *Cell) Kill() {
	c.Alive = false
}

func (c *Cell) Copy() *Cell {
	return &Cell{c.Alive}
}

type World struct {
	Cells map[string]*Cell
}

func NewWorld() *World {
	cells := make(map[string]*Cell)
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	for row := 0; row < WorldHeight; row++ {
		for col := 0; col < WorldWidth; col++ {
			coords := Coord{row, col}
			cell := NewCell()
			cells[coords.String()] = cell
			cell.SetState(InitialSpawnTolerance >= randGen.Float64())
		}
	}

	return &World{cells}
}

func (w *World) Step() {
	pastState := make(map[string]*Cell)
	for coordStr, cell := range w.Cells {
		pastState[coordStr] = cell.Copy()
	}

	for coordStr, nextCell := range w.Cells {
		pastCell := pastState[coordStr]
		livingNeighbors := 0
		coords := CoordFromString(coordStr)

		for _, neighborCoords := range neighbors(coords) {
			neighbor, ok := pastState[neighborCoords]
			if ok && neighbor.Alive {
				livingNeighbors++
			}
		}

		if pastCell.Alive {
			if livingNeighbors < 2 {
				// underpopulation
				nextCell.Kill()
			} else if livingNeighbors > 3 {
				// overpopulation
				nextCell.Kill()
			}
		} else {
			if livingNeighbors == 3 {
				// reproduce
				nextCell.Spawn()
			}
		}
	}
}

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

func main() {
	w := NewWorld()
	for {
		clearConsole()
		w.Draw()
		w.Step()
		time.Sleep(time.Second / 2)
	}
}
