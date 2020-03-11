WORLD_WIDTH = 20
WORLD_HEIGHT = 10
INITIAL_SPAWN_TOLERANCE = 0.4
MAX_STEPS = 100

"""Clear the terminal."""
def clear_console
    system("clear")
end

"""Check if a cell falls within grid bounds."""
def inbounds(coords)
    row = coords.row
    col = coords.col
    row >= 0 and row < WORLD_HEIGHT and col >= 0 and col < WORLD_WIDTH
end

"""Find all in-bounds neighbors of a given 2D coordinate."""
def neighbors(coords)
    row = coords.row
    col = coords.col

    neighbors = []
    triplet = [-1, 0, 1]
    for row_delta in triplet
        for col_delta in triplet
            if row_delta == 0 and col_delta == 0
                # this is the cell itself
                next
            end

            n_row = row + row_delta
            n_col = col + col_delta
            n_coords = Coord.new(n_row, n_col)
            if inbounds(n_coords)
                neighbors.push(n_coords.to_s)
            end
        end
    end

    neighbors
end

"""Coord describes a 2D position in a grid."""
class Coord
    """Initialize Coord."""
    def initialize(row, col)
        @row = row
        @col = col
    end

    """Row getter."""
    def row
        @row
    end

    """Col getter."""
    def col
        @col
    end

    """Stringify a Coord."""
    def to_s
        "#{@row},#{@col}"
    end

    """Given a string of coordinates, build a new Coord."""
    def self.from_s(coord_str)
        parts = coord_str.split(',')
        coord = self.new(parts[0].to_i, parts[1].to_i)
        coord
    end
end

"""Cell represents a single living entity in the grid."""
class Cell
    """Initialize Cell."""
    def initialize
        @alive = false
    end

    """Alive getter."""
    def alive
        @alive
    end

    """Set the cell's alive state."""
    def set_state(state)
        @alive = state
    end

    """Spawn the cell."""
    def spawn
        @alive = true
    end

    """Kill off the cell."""
    def kill
        @alive = false
    end

    """Make a copy of the cell."""
    def copy
        cell = Cell.new
        cell.set_state(@alive)
        cell
    end
end

"""World is a 2D grid filled with Cells."""
class World
    """Build a new world with cells in a random state."""
    def initialize()
        cells = Hash.new

        for row in 0..WORLD_HEIGHT
            for col in 0..WORLD_WIDTH
                coords = Coord.new(row, col)
                cell = Cell.new
                cells[coords.to_s] = cell
                cell.set_state(INITIAL_SPAWN_TOLERANCE >= rand())
            end
        end

        @cells = cells
    end

    """Apply automata rules to all cells in the grid."""
    def step()
        # Make a new map containing the future state of the world. Game of Life rules are based on
        # current timestep. We will use this to maintain next state.
        next_state = Hash.new

        @cells.each do |coord_str, next_cell|
            past_cell = @cells[coord_str]
            living_neighbors = 0
            coords = Coord.from_s(coord_str)
            next_cell = past_cell.copy

            # Grab the number of living cells surrounding the current cell.
            for neighbor_coords in neighbors(coords)
                if @cells.key?(neighbor_coords)
                    neighbor = @cells[neighbor_coords] 
                    if neighbor.alive
                        living_neighbors += 1
                    end
                end
            end

            # Apply Conway's rules.
            if past_cell.alive
                if living_neighbors < 2
                    # Kill due to underpopulation
                    next_cell.kill()
                elsif living_neighbors > 3
                    # Kill due to overpopulation
                    next_cell.kill()
                end
            else
                if living_neighbors == 3
                    # Reproduce
                    next_cell.spawn()
                end
            end

            next_state[coord_str] = next_cell
        end

        @cells = next_state
    end

    """Print the current state of the world to terminal."""
    def draw()
        output_str = ""

        for row in 0..WORLD_HEIGHT
            for col in 0..WORLD_WIDTH
                coords = Coord.new(row, col)
                cell = @cells[coords.to_s]

                if cell.alive
                    output_str += " 0 "
                else
                    output_str += " . "
                end
            end

            output_str += "\n\r"
        end

        puts output_str
    end
end

"""Play Conway's Game of Life."""
def main
    world = World.new
    for _ in 0..MAX_STEPS
        clear_console
        world.draw
        world.step
        sleep(0.5)
    end

end

main()
