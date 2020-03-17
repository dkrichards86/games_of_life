package com.dkrichards.gamesoflife;

/**
 * Coord describes a 2D position in a grid.
 */
public class Coord {
    public int row;
    public int col;

    /**
     * Initialize Coord.
     */
    public Coord(int row, int col) {
        this.row = row;
        this.col = col;
    }

    /**
     * Stringify a Coord.
     */
    public String toString() {
        return String.format("%d,%d", this.row, this.col); 
    }

    /**
     * Given a string of coordinates, build a new Coord.
     */
    public static Coord fromString(String coordStr) {
        String[] parts = coordStr.split(",");

        return new Coord(Integer.parseInt(parts[0]), Integer.parseInt(parts[1]));
    }
}