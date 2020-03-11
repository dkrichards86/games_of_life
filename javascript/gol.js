const WORLD_HEIGHT = 10;
const WORLD_WIDTH = 20;
const INITIAL_SPAWN_TOLERANCE = 0.4;
const MAX_STEPS = 100;

/**
 * Sleep for a given number of milliseconds.
 */
const sleep = (millis) => new Promise(resolve => setTimeout(resolve, millis));

/**
 * Clear the terminal.
 */
const clearConsole = () => console.clear();

/**
 * Check if a cell falls within grid bounds.
 */
const inbounds = (coords) => {
	const row = coords.row;
	const col = coords.col;
	return row >= 0 && row < WORLD_HEIGHT && col >= 0 && col < WORLD_WIDTH;
};

/**
 * Find all in-bounds neighbors of a given 2D coordinate.
 */
const neighbors = (coords) => {
	const row = coords.row;
	const col = coords.col;
    let neighbors = [];

    const triplet = [-1, 0, 1];
    triplet.forEach(rowDelta => {
        triplet.forEach(colDelta => {
            if (rowDelta === 0 && colDelta === 0) {
                // do nothing
            } else {
                const nRow = row + rowDelta;
                const nCol = col + colDelta;
                const nCoords = new Coord(nRow, nCol);
                if (inbounds(nCoords)) {
                    neighbors.push(nCoords.toString());
                }
            }
        })
    });

	return neighbors;
};

/**
 * Coord describes a 2D position in a grid.
 */
class Coord {
    /**
     * Initialize Coord.
     */
    constructor(row, col) {
        this.row = row;
        this.col = col;
    }

    /**
     * Stringify a Coord.
     */
    toString() {
        return `${this.row},${this.col}`;
    }

    /**
     * Given a string of coordinates, build a new Coord.
     */
    static fromString(coordString) {
        const parts = coordString.split(',');
        return new Coord(parseInt(parts[0]), parseInt(parts[1]));
    }
}

/**
 * Cell represents a single living entity in the grid.
 */
class Cell {
    /**
     * Initialize Cell.
     */
    constructor() {
        this.alive = false;
    }

    /**
     * Set the cell's alive state.
     */
    setState(state) {
        this.alive = state;
    }

    /**
     * Kill off the cell.
     */
    kill() {
        this.alive = false;
    }

    /**
     * Spawn the cell.
     */
    spawn(){
        this.alive = true;
    }

    /**
     * Make a copy of the cell.
     */
    copy() {
        let cell = new Cell();
        cell.setState(this.alive);
        return cell;
    }
}

/**
 * World is a 2D grid filled with Cells.
 */
class World {
    /**
     * Build a new world with cells in a random state.
     */
    constructor() {
        let cells = {}

        for (let row = 0; row < WORLD_HEIGHT; row++) {
            for (let col = 0; col < WORLD_WIDTH; col++) {
                const coords = new Coord(row, col);
                const cell = new Cell();
                cells[coords.toString()] = cell;
                cell.setState(INITIAL_SPAWN_TOLERANCE >= Math.random());
            }
        }

        this.cells = cells;
    }

    /**
     * Apply automata rules to all cells in the grid.
     */
    step() {
        // Make a new map containing the future state of the world. Game of Life rules are based on
        // current timestep. We will use this to maintain next state.
        const nextState = {};

        Object.keys(this.cells).forEach(coordStr => {
            const pastCell = this.cells[coordStr];
            let livingNeighbors = 0;
            const coords = Coord.fromString(coordStr);
            const nextCell = pastCell.copy();

            // Grab the number of living cells surrounding the current cell.
            neighbors(coords).forEach(neighborCoords => {
                const neighbor = this.cells[neighborCoords];
                if (neighbor && neighbor.alive) {
                    livingNeighbors++;
                }
            });

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

            nextState[coordStr] = nextCell;
        });

        this.cells = nextState;
    }

    /**
     * Print the current state of the world to terminal.
     */
    draw() {
        let outputString = "";
        for (let row = 0; row < WORLD_HEIGHT; row++) {
            for (let col = 0; col < WORLD_WIDTH; col++) {
                const coords = new Coord(row, col);
                const cell = this.cells[coords.toString()];
                if (cell.alive) {
                    outputString += " 0 ";
                } else {
                    outputString += " . ";
                };
            }
            outputString += "\n\r"
        }
        console.log(outputString);
    }
}

/**
 * Play Conway's Game of Life.
 */
const main = async () => {
    let world = new World();
    for (let i = 0; i < MAX_STEPS; i++) {
        clearConsole();
        world.draw()
        world.step()
        await sleep(500);
    }
};

main();
