package com.dkrichards.gamesoflife;

/**
 * 
 */
public class Cell {
    public boolean alive;

    /**
    * Initialize Cell.
    */
    public Cell() {
        this.alive = false;
    }

    /**
    * Set the cell's alive state.
    */
    public void setState(Boolean state) {
        this.alive = state;
    }

    /**
    * Spawn a new cell.
    */
    public void spawn() {
        this.alive = true;
    }

    /**
    * Kill off the cell.
    */
    public void kill() {
        this.alive = false;
    }

    /**
    * Make a copy of the cell.
    */
    public Cell copy() {
        Cell cell = new Cell();
        cell.setState(this.alive);
        return cell;
    }
}