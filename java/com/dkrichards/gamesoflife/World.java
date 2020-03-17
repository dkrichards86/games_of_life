package com.dkrichards.gamesoflife;

import java.io.*;
import java.util.*; 

/**
 * World is a 2D grid filled with Cells.
 */
public class World {
    HashMap<String, Cell> cells = new HashMap<>();
    int worldWidth;
    int worldHeight;

    /**
     * Build a new world with cells in a random state.
     */
    public World() throws IOException, FileNotFoundException {
        Random rand = new Random();

        GameProperties props = GameProperties.getInstance();
        this.worldWidth = props.getWorldWidth();
        this.worldHeight = props.getWorldHeight();
        double spawnTolerance = props.getInitialSpawnTolerance();

        for (int row = 0; row < worldHeight; row++) {
            for (int col = 0; col < worldWidth; col++) {
                Coord coords = new Coord(row, col);
                Cell cell = new Cell();
                cell.setState(spawnTolerance >= rand.nextDouble());
                this.cells.put(coords.toString(), cell);
            }
        }
    }

    /**
     * Check if a cell falls within grid bounds.
     */
    private boolean inBounds(Coord coords) {
        int row = coords.row;
        int col = coords.col;
        return row >= 0 && row < this.worldHeight && col >= 0 && col < this.worldWidth;
    }

    /**
     * Find all in-bounds neighbors of a given 2D coordinate.
     */
    private List<String> neighbors(Coord coords) {
        int row = coords.row;
        int col = coords.col;

        List<String> neighbors = new ArrayList<>();
        int[] triplet = {-1, 0, 1}; 

        int neighborIdx = 0;
        for (int i = 0; i < 3; i++) {
            int rowDelta = triplet[i];
            for (int j = 0; j < 3; j++) {
                int colDelta = triplet[j];
                if (rowDelta == 0 && colDelta == 0) {
                    // this is the cell itself
                    continue;
                }

                int nRow = row + rowDelta;
                int nCol = col + colDelta;
                Coord nCoords =  new Coord(nRow, nCol);
                if (this.inBounds(nCoords)) {
                    neighbors.add(nCoords.toString());
                }

                neighborIdx++;
            }
        }

        return neighbors;
    }

    /**
     * Apply automata rules to all cells in the grid.
     */
    public void step() {
        // Make a new map containing the future state of the world. Game of Life rules are based on
        // current timestep. We will use this to maintain next state.
        HashMap<String, Cell> nextState = new HashMap<>();

        Iterator iterator = this.cells.entrySet().iterator(); 

        while (iterator.hasNext()) { 
            Map.Entry entry = (Map.Entry)iterator.next(); 
            String coordStr = ((String)entry.getKey());
            Cell pastCell = this.cells.get(coordStr);
            int livingNeighbors = 0;
            Coord coords = Coord.fromString(coordStr);
            Cell nextCell = pastCell.copy();

            // Grab the number of living cells surrounding the current cell.
            for (String neighborCoords : this.neighbors(coords)) {
                if (this.cells.containsKey(neighborCoords)) {
                    Cell neighbor = this.cells.get(neighborCoords);
                    if (neighbor.alive) {
                        livingNeighbors++;
                    }
                }
            }

            // Apply Conway's rules.
            if (pastCell.alive) {
                if (livingNeighbors < 2) {
                    // Kill due to underpopulation
                    nextCell.kill();
                } else if (livingNeighbors > 3) {
                    // Kill due to overpopulation
                    nextCell.kill();
                }
            } else {
                if (livingNeighbors == 3) {
                    // Reproduce
                    nextCell.spawn();
                }
            }
            nextState.put(coordStr, nextCell);
        }
	    this.cells = nextState;
    }

    /**
     * Print the current state of the world to terminal.
     */
    public void draw() throws IOException, FileNotFoundException {
        StringBuilder outputStr = new StringBuilder();

        GameProperties props = GameProperties.getInstance();
        int worldWidth = props.getWorldWidth();
        int worldHeight = props.getWorldWidth();

        for (int row = 0; row < worldHeight; row++) {
            for (int col = 0; col < worldWidth; col++) {
                Coord coords = new Coord(row, col);
                String coordsStr = coords.toString();
                if (!this.cells.containsKey(coordsStr)) {
                    continue;
                }
                Cell cell = this.cells.get(coordsStr);

                if (cell.alive) {
                    outputStr.append(" 0 ");
                }
                else {
                    outputStr.append(" . ");
                }
            }
            outputStr.append("\n");
        }
        System.out.println(outputStr);
    }
} 
