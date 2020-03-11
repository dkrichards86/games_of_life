import os
import random
from copy import deepcopy
from time import sleep

WORLD_WIDTH = 20
WORLD_HEIGHT = 10
INITIAL_SPAWN_TOLERANCE = 0.4
MAX_STEPS = 100


def clear_console():
    """Clear the terminal."""
    os.system("clear")


def inbounds(coords):
    """Check if a cell falls within grid bounds."""
    row = coords.row
    col = coords.col
    return row >= 0 and row < WORLD_HEIGHT and col >= 0 and col < WORLD_WIDTH


def neighbors(coords):
    """Find all in-bounds neighbors of a given 2D coordinate."""
    row = coords.row
    col = coords.col

    neighbors = []
    triplet = [-1, 0, 1]
    for row_delta in triplet:
        for col_delta in triplet:
            if row_delta == 0 and col_delta == 0:
                # this is the cell itself
                continue

            n_row = row + row_delta
            n_col = col + col_delta
            n_coords = Coord(n_row, n_col)
            if inbounds(n_coords):
                neighbors.append(str(n_coords))

    return neighbors


class Coord:
    """Coord describes a 2D position in a grid."""

    def __init__(self, row, col):
        """Initialize Coord."""
        self.row = row
        self.col = col

    @classmethod
    def from_string(cls, coord_str):
        """Given a string of coordinates, build a new Coord."""
        parts = coord_str.split(',')
        return cls(int(parts[0]), int(parts[1]))

    def __str__(self):
        """Stringify a Coord."""
        return "{},{}".format(self.row, self.col)


class Cell:
    """Cell represents a single living entity in the grid."""

    def __init__(self):
        """Initialize Cell."""
        self.alive = False

    def set_state(self, state):
        """Set the cell's alive state."""
        self.alive = state

    def spawn(self):
        """Spawn the cell."""
        self.alive = True

    def kill(self):
        """Kill off the cell."""
        self.alive = False

    def copy(self):
        """Make a copy of the cell.."""
        cell = Cell()
        cell.set_state(self.alive)
        return cell


class World:
    """World is a 2D grid filled with Cells."""

    def __init__(self):
        """Build a new world with cells in a random state."""
        cells = dict()

        for row in range(WORLD_HEIGHT):
            for col in range(WORLD_WIDTH):
                coords = Coord(row, col)
                cell = Cell()
                cells[str(coords)] = cell
                cell.set_state(INITIAL_SPAWN_TOLERANCE >= random.random())
        self.cells = cells

    def step(self):
        """Apply automata rules to all cells in the grid."""
        # Make a new map containing the future state of the world. Game of Life rules are based on
        # current timestep. We will use this to maintain next state.
        next_state = dict()

        for coord_str in self.cells.keys():
            past_cell = self.cells[coord_str]
            living_neighbors = 0
            coords = Coord.from_string(coord_str)
            next_cell = past_cell.copy()

            # Grab the number of living cells surrounding the current cell.
            for neighbor_coords in neighbors(coords):
                try:
                    neighbor = self.cells[neighbor_coords]
                    if neighbor.alive:
                        living_neighbors += 1
                except KeyError:
                    continue

            # Apply Conway's rules.
            if past_cell.alive:
                if living_neighbors < 2:
                    # Kill due to underpopulation
                    next_cell.kill()
                elif living_neighbors > 3:
                    # Kill due to overpopulation
                    next_cell.kill()
            else:
                if living_neighbors == 3:
                    # Reproduce
                    next_cell.spawn()

            next_state[coord_str] = next_cell
        self.cells = next_state

    def draw(self):
        """Print the current state of the world to terminal."""
        output_str = ""

        for row in range(WORLD_HEIGHT):
            for col in range(WORLD_WIDTH):
                coords = Coord(row, col)
                cell = self.cells[str(coords)]

                if cell.alive:
                    output_str += " 0 "
                else:
                    output_str += " . "

            output_str += "\n\r"

        print(output_str, end=" ")


def main():
    """Play Conway's Game of Life."""
    world = World()
    for _ in range(MAX_STEPS):
        clear_console()
        world.draw()
        world.step()
        sleep(0.5)

if __name__ == '__main__':
    main()