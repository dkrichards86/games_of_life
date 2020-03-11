extern crate rand;

use std::collections::HashMap;
use std::fmt;
use std::process::Command;
use std::{thread, time};

const WORLDWIDTH: i32 = 20;
const WORLDHEIGHT: i32 = 10;
const INITIALSPAWNTOLERANCE: f32 = 0.4;
const MAXSTEPS: i32 = 100;

fn clear_console() {
    // Clear the terminal.
    let output = Command::new("clear").output().unwrap();
    println!("{}", String::from_utf8_lossy(&output.stdout));
}

fn inbounds(coords: Coord) -> bool {
    // Check if a cell falls within grid bounds.
    let row = coords.row;
    let col = coords.col;
	return row >= 0 && row < WORLDHEIGHT && col >= 0 && col < WORLDWIDTH
}

fn neighbors(coords: &Coord) -> Vec<String> {
    // Find all in-bounds neighbors of a given 2D coordinate.
    let row = coords.row;
    let col = coords.col;

    let mut neighbors: Vec<String> = Vec::new();
    let triplets: [i32; 3] = [-1, 0, 1];

    for row_delta in &triplets {
        for col_delta in &triplets {
            if *row_delta == 0 && *col_delta == 0 {
                // this is the cell itself
                continue
            }

            let n_row = row + *row_delta;
            let n_col = col + *col_delta;
            let n_coords = Coord::new(n_row, n_col);

            if inbounds(n_coords) {
                neighbors.push(n_coords.to_string());
            }
        }
    }

    neighbors
}

#[derive(Clone, Copy)]
// Coord describes a 2D position in a grid.
struct Coord {
    row: i32,
    col: i32
}

impl Coord {
    fn new(row: i32, col: i32) -> Self {
        // Initialize Coord.
        Coord {row: row, col: col}
    }

    fn from_string(coord_str: String) -> Self {
        // Given a string of coordinates, build a new Coord.
        let parts:Vec<&str>= coord_str.split(",").collect();
        let row = parts[0].parse::<i32>().unwrap();
        let col = parts[1].parse::<i32>().unwrap();
        Coord {row: row, col: col}
    }
}

// Stringify a Coord.
impl fmt::Display for Coord {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{},{}", self.row, self.col)
    }
}

#[derive(Clone, Copy)]
// Cell represents a single living entity in the grid.
struct Cell {
    alive: bool
}

impl Cell {
    fn new() -> Self {
        // Initialize Cell.
        Cell {alive: false}
    }

    fn set_state(&mut self, state: bool) {
        // Set the cell's alive state.
        self.alive = state;
    }

    fn spawn(&mut self) {
        // Spawn a new cell.
        self.alive = true;
    }

    fn kill(&mut self) {
        // Kill off the cell.
        self.alive = false;
    }

    fn copy(self) -> Self {
        // Make a copy of the cell.
        Cell {alive: self.alive}
    }
}

impl fmt::Display for Cell {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Cell {}", self.alive)
    }
}

// World is a 2D grid filled with Cells.
struct World {
    cells: HashMap<String, Cell>
}

impl World {
    fn new() -> Self {
        // Build a new world with cells in a random state.

        let mut cells = HashMap::new();

        for row in 0..WORLDHEIGHT {
            for col in 0..WORLDWIDTH {
                let coords = Coord::new(row, col);
                let mut cell = Cell::new();
                cell.set_state(INITIALSPAWNTOLERANCE >= rand::random::<f32>());
                cells.insert(coords.to_string(), cell);
            }
        }

        World {cells: cells}
    }

    fn step(&mut self) {
        // Apply automata rules to all cells in the grid.

        // Make a new map containing the future state of the world. Game of Life rules are based on
        // current timestep. We will use this to maintain next state.
        let mut next_state = HashMap::new();

        for coord_str in self.cells.keys() {
            let past_cell = self.cells.get(coord_str).unwrap();
            let mut living_neighbors = 0;
            let coords = Coord::from_string(coord_str.to_string());
            let mut next_cell = past_cell.copy();

            // Grab the number of living cells surrounding the current cell.
            for neighbor_coords in neighbors(&coords) {
                let neighbor = self.cells.get(&neighbor_coords);
                match neighbor {
                    Some(x) => {
                        if x.alive {
                            living_neighbors += 1;
                        }
                    },
                    _ => {}
                }
            }

            // Apply Conway's rules.
            if past_cell.alive {
                if living_neighbors < 2 {
                    // Kill due to underpopulation
                    next_cell.kill();
                } else if living_neighbors > 3 {
                    // Kill due to overpopulation
                    next_cell.kill();
                }
            } else {
                if living_neighbors == 3 {
                    // Reproduce
                    next_cell.spawn();
                }
            }

            next_state.insert(coord_str.to_string(), next_cell);
        }

        self.cells = next_state;
    }

    fn draw(&self) {
        // Print the current state of the world to terminal.

        let mut output_str = String::from("");

        for row in 0..WORLDHEIGHT {
            for col in 0..WORLDWIDTH {
                let coords = Coord::new(row, col);
                let cell = &self.cells.get(&coords.to_string()).unwrap();

                if cell.alive {
                    output_str.push_str(" 0 ");
                } else {
                    output_str.push_str(" . ");
                }
            }
            output_str.push_str("\n\r");
        }

        println!("{}", output_str);
    }
}

fn main() {
    let mut world = World::new();
    for _ in 0..MAXSTEPS {
        clear_console();
        world.draw();
        world.step();
        thread::sleep(time::Duration::from_millis(500));
    }
}
