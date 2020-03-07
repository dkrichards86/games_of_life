import os
import random
from copy import deepcopy
from time import sleep

WORLD_WIDTH = 20
WORLD_HEIGHT = 10
INITIAL_SPAWN_TOLERANCE = 0.4


def clear_console():
    os.system("clear")


def inbounds(coords):
    row = coords.row
    col = coords.col
    return row >= 0 and row < WORLD_HEIGHT and col >= 0 and col < WORLD_WIDTH


def neighbors(coords):
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
    def __init__(self, row, col):
        self.row = row
        self.col = col

    @classmethod
    def from_string(cls, coord_str):
        parts = coord_str.split(',')
        return cls(int(parts[0]), int(parts[1]))

    def __str__(self):
        return "{},{}".format(self.row, self.col)


class Cell:
    def __init__(self):
        self.alive = False

    def set_state(self, state):
        self.alive = state

    def spawn(self):
        self.alive = True

    def kill(self):
        self.alive = False

    def copy(self):
        cell = Cell()
        cell.set_state(self.alive)
        return cell

class World:
    def __init__(self):
        cells = dict()

        for row in range(WORLD_HEIGHT):
            for col in range(WORLD_WIDTH):
                coords = Coord(row, col)
                cell = Cell()
                cells[str(coords)] = cell
                cell.set_state(INITIAL_SPAWN_TOLERANCE >= random.random())
        self.cells = cells

    def step(self):
        past_state = dict()
        for coord_str, next_cell in self.cells.items():
            past_state[coord_str] = next_cell.copy()

        for coord_str, next_cell in self.cells.items():
            past_cell = past_state[coord_str]
            living_neighbors = 0
            coords = Coord.from_string(coord_str)

            for neighbor_coords in neighbors(coords):
                try:
                    neighbor = past_state[neighbor_coords]
                    if neighbor.alive:
                        living_neighbors += 1
                except KeyError:
                    continue

            if past_cell.alive:
                if living_neighbors < 2:
                    # underpopulation
                    next_cell.kill()
                elif living_neighbors > 3:
                    # overpopulation
                    next_cell.kill()
            else:
                if living_neighbors == 3:
                    # reproduce
                    next_cell.spawn()

    def draw(self):
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


if __name__ == '__main__':
    world = World()
    while True:
        clear_console()
        world.draw()
        world.step()
        sleep(0.5)
