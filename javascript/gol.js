const WORLD_HEIGHT = 10;
const WORLD_WIDTH = 20;
const INITIAL_SPAWN_TOLERANCE = 0.4;

const sleep = (millis) => {
    return new Promise(resolve => setTimeout(resolve, millis));
};

const inbounds = (coords) => {
	const row = coords.row;
	const col = coords.col;
	return row >= 0 && row < WORLD_HEIGHT && col >= 0 && col < WORLD_WIDTH;
};

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

class Coord {
    constructor(row, col) {
        this.row = row;
        this.col = col;
    }

    toString() {
        return `${this.row},${this.col}`;
    }

    static fromString(coordString) {
        const parts = coordString.split(',');
        return new Coord(parseInt(parts[0]), parseInt(parts[1]));
    }
}

class Cell {
    constructor() {
        this.alive = false;
    }

    setState(state) {
        this.alive = state;
    }

    kill() {
        this.alive = false;
    }

    spawn(){
        this.alive = true;
    }

    copy() {
        let cell = new Cell();
        cell.setState(this.alive);
        return cell;
    }
}

class World {
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

    step() {
        const pastState = {};
        Object.entries(this.cells).forEach(entry => {
            const [ coordStr, nextCell ] = entry;
            pastState[coordStr] = nextCell.copy();
        });

        Object.entries(this.cells).forEach(entry => {
            const [ coordStr, nextCell ] = entry;
            const coords = Coord.fromString(coordStr);
            const pastCell = pastState[coordStr];
            let livingNeighbors = 0;

            neighbors(coords).forEach(neighborCoords => {
                const neighbor = pastState[neighborCoords];
                if (neighbor && neighbor.alive) {
                    livingNeighbors++;
                }
            });

            if (pastCell.alive) {
                if (livingNeighbors < 2) {
                    // underpopulation
                    nextCell.kill();
                } else if (livingNeighbors > 3) {
                    // overpopulation
                    nextCell.kill();
                }
            } else {
                if (livingNeighbors == 3) {
                    // reproduce
                    nextCell.spawn();
                }
            }
        });
    }

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

const main = async () => {
    let world = new World();
    while (true) {
        console.clear();
        world.draw()
        world.step()
        await sleep(500);
    }
};

main();
