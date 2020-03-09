<?php
    const WORLD_WIDTH = 20;
    const WORLD_HEIGHT = 10;
    const INITIAL_SPAWN_TOLERANCE = 0.4;
    const MAX_STEPS = 100;

    /**
     * Sleep for a given number of milliseconds.
     */
    function sleep_millis($millis) {
        time_nanosleep(0, $millis * 1000 * 1000);
    }

    /**
     * Clear the terminal.
     */
    function clearConsole() {
        system('clear');
    }

    /**
     * Check if a cell falls within grid bounds.
     */
    function inbounds($coords) {
        $row = $coords->row;
        $col = $coords->col;
        return $row >= 0 && $row < WORLD_HEIGHT && $col >= 0 && $col < WORLD_WIDTH;
    }

    /**
     * Find all in-bounds neighbors of a given 2D coordinate.
     */
    function neighbors($coords) {
        $row = $coords->row;
        $col = $coords->col;
        $neighbors = [];

        $triplet = [-1, 0, 1];
        $arrlength = count($triplet);
        for($i = 0; $i < $arrlength; $i++) {
            $rowDelta = $triplet[$i];
            for($j = 0; $j < $arrlength; $j++) {
                $colDelta = $triplet[$j];
                if ($rowDelta === 0 && $colDelta === 0) {
                    // do nothing
                } else {
                    $nRow = $row + $rowDelta;
                    $nCol = $col + $colDelta;
                    $nCoords = new Coord($nRow, $nCol);
                    if (inbounds($nCoords)) {
                        array_push($neighbors, $nCoords->toString());
                    }
                }
            }
        }

        return $neighbors;
    }

    /**
     * Coord describes a 2D position in a grid.
     */
    class Coord {
        public $row;
        public $col;

        /**
         * Initialize Coord.
         */
        function __construct($row, $col) {
            $this->row = $row;
            $this->col = $col;
        }

        /**
         * Stringify a Coord.
         */
        function toString() {
            return sprintf("%s,%s", $this->row, $this->col);
        }

        /**
         * Given a string of coordinates, build a new Coord.
         */
        public static function fromString($coordString) {
            $parts = explode(',', $coordString);
            return new Coord(intval($parts[0]), intval($parts[1]));
        }
    }

    /**
     * Cell represents a single living entity in the grid.
     */
    class Cell {
        public $alive;

        /**
         * Initialize Cell.
         */
        function __construct() {
            $this->alive = false;
        }

        /**
         * Set the cell's alive state.
         */
        public function setState($state) {
            $this->alive = $state;
        }

        /**
         * Kill off the cell.
         */
        public function kill() {
            $this->alive = false;
        }

        /**
         * Spawn the cell.
         */
        public function spawn() {
            $this->alive = true;
        }
    }

    /**
     * World is a 2D grid filled with Cells.
     */
    class World {
        /**
         * Build a new world with cells in a random state.
         */
        function __construct() {
            $cells = array();

            for ($row = 0; $row < WORLD_HEIGHT; $row++) {
                for ($col = 0; $col < WORLD_WIDTH; $col++) {
                    $coords = new Coord($row, $col);
                    $cell = new Cell();
                    $cell->setState(boolval(INITIAL_SPAWN_TOLERANCE >= mt_rand() / mt_getrandmax()));
                    $cells[$coords->toString()] = $cell;
                }
            }

            $this->cells = $cells;
        }

        /**
         * Apply automata rules to all cells in the grid.
         */
        public function step() {
            // Make a deep copy of the state of the world. Game of Life rules are based on current
            // timestep. We will use this to determine next state. 
            $pastState = array();
            foreach($this->cells as $coordStr => $nextCell) {
                $pastState[$coordStr] = clone $nextCell;
            }

            foreach($this->cells as $coordStr => $nextCell) {
                $pastCell = $pastState[$coordStr];
                $livingNeighbors = 0;
                $coords = Coord::fromString($coordStr);

                // Grab the number of living cells surrounding the current cell.
                foreach(neighbors($coords) as $neighborCoords) {
                    if (!array_key_exists($neighborCoords, $pastState)) {
                        continue;
                    }
                    $neighbor = $pastState[$neighborCoords];

                    if ($neighbor->alive) {
                        $livingNeighbors++;
                    }
                }

                // Apply Conway's rules.
                if ($pastCell->alive) {
                    if ($livingNeighbors < 2) {
                        // Kill due to underpopulation
                        $nextCell->kill();
                    } else if ($livingNeighbors > 3) {
                        // Kill due to overpopulation
                        $nextCell->kill();
                    }
                } else {
                    if ($livingNeighbors == 3) {
                        // Reproduce
                        $nextCell->spawn();
                    }
                }
            }
        }

        /**
         * Print the current state of the world to terminal.
         */
        public function draw() {
            $outputString = "";
            for ($row = 0; $row < WORLD_HEIGHT; $row++) {
                for ($col = 0; $col < WORLD_WIDTH; $col++) {
                    $coords = new Coord($row, $col);
                    $cell = $this->cells[$coords->toString()];
                    if ($cell->alive) {
                        $outputString .= " 0 ";
                    } else {
                        $outputString .= " . ";
                    }
                }
                $outputString .= "\n\r";
            }
            echo $outputString;
        }
    }

    /**
     * Play Conway's Game of Life.
     */
    function main() {
        $world = new World();

        for ($i = 0; $i < MAX_STEPS; $i++) {
            clearConsole();
            $world->draw();
            $world->step();
            sleep_millis(500);
        }
    }

    main();
?>