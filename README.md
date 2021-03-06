# Games of Life
The goal of this repository is to compare various languages using Conway's Game of Life as a basis.
As much as possible, I tried to create an apples to apples comparison. In some cases, I forgo
language idioms in favor polyglot approaches. I also try to use standard library for all functionality.

## Game of Life Rules
Game of Life replicates cellular automation in a two dimensional grid. The rules are as follows:
- a living cell with less than two live neighbors dies, as if by underpopulation.
- a living cell with more than three live neighbors dies, as if by overpopulation.
- any dead cell with three living neighbors becomes a live cell, as if by reproduction.

## Why Game of Life
Game of Life provides simple rules, but affords us the opportunity to explore a language's features.
It:
- includes basic language syntax anc constructs like conditionals and loops,
- manipulates string,
- manipualtes objects,
- includes cursory functional programming patterns,
- leverages `static` methods,
- works with pointers/borrowership/memory management,
- draws to console

## Languages Represented
Game of Life is currently translated into the following langauges:
- Python
- JavaScript
- Go
- PHP
- Ruby
- Java
