# go-sudoku

## Run

## Board format

## Additional arguments

Run `go-sudoku -help` for a full list of optional arguments.

| Argument                | Description
|-------------------------|-------------
| `-time`                 | print time to solve (TODO)
| `-generate [level]`     | generate a sudoku puzzle with the given difficulty. valid options are `easy`, `medium`, and `hard`. (TODO)
| `-grade`                | grades a puzzle as `easy`, `medium`, or `hard` (TODO)
| `-profile`              | enable CPU and memory profiling
| `-board [compact-form]` | accepts a board in compact form (TODO)
| `-compact`              | prints the solved board in compact form (TODO)
| `-convert-compact`      | prints the board from stdin in compact form without solving (TODO)
| `-file`                 | run a set of Sudoku puzzles from a file.
| `-max-puzzles`          | use with `-file` to limit the number of puzzles executed
| `-time-only`            | outputs a tab delimited list of puzzle numbers and time. useful for finding puzzles which are challenging. (TODO)
| `-sat`                  | only use SAT to solve the puzzle(s) (TODO)
| `-sat-output`           | prints the CNF format used by many SAT solvers such as minisat
| `-stategies`            | list. valid values listed below.

## How it works

`go-sudoku` first attempts human strategy and ultimately falls back on a SAT solver

## Terminology

- Row
- Column
- Box
- Candidate / Hint

## Strategies

| Name               | Argument        | Brief Description
|--------------------|-----------------|-------------------------
| **Basic**          |                 |
| Naked Singles      | `n1`            | A cell has only one hint left; it can be solved with this hint.
| Hidden Singles     | `h1`            | A cell is the only one in a row, column, or box with a given hint.
| **Easy**           |                 |
| Naked Pairs        | `n2`            |
| Naked Triples      | `n3`            |
| Naked Quads        | `n4`            |
| Naked Quints       | `n5`            |
| Hidden Pairs       | `h2`            |
| Hidden Triples     | `h3`            |
| Hidden Quads       | `h4`            |
| Hidden Quints      | `h5`            |
| **Moderate**       |                 |
| Pointing Pairs     | `pointing-pair` |
| Box/Line Reduction | `box-line`      |
| X-Wing             | `xwing`         |
| **Tough**          |                 |
| Simple Coloring    | `simple-color`  |
| Y-Wing             | `ywing`         |
| Sword-Fish         | `swordfish`     |
| XY-Chain           | `xychain`       |
| Empty Rectangles   | `empty-rect`    |

Note:
- Difficulties are based on listings at [sudokiwiki.org](http://www.sudokuwiki.org/sudoku.htm).
- Difficulties are subjective to the player.
- Robust grading algorithms are often commercial products used in validating the worthiness of a puzzle for publication.

## Things learned

### Go

pprof - cpu, memory - invaluable for finding the largest bottle necks
vet caught errors (in generate comparing 
go lint removed redundant code. for example, `[][]SetVar{SetVar{VarNum: 1, Value: true}}` became `[][]SetVar{{VarNum: 1, Value: true}}`

### Sudoku

## Resources

- http://www.sudokuwiki.org/Strategy_Families
- https://en.wikipedia.org/wiki/Simplex_algorithm
- http://www.nature.com/articles/srep00725#f3
- http://arxiv.org/ftp/arxiv/papers/0805/0805.0697.pdf
- Satisfiability Solvers: http://www.cs.cornell.edu/gomes/papers/satsolvers-kr-handbook.pdf
- http://ocw.mit.edu/courses/electrical-engineering-and-computer-science/6-005-elements-of-software-construction-fall-2011/assignments/MIT6_005F11_ps4.pdf
- https://en.wikipedia.org/wiki/Conjunctive_normal_form
- https://en.wikipedia.org/wiki/DPLL_algorithm
- https://en.wikipedia.org/wiki/Backtracking
- https://en.wikipedia.org/wiki/Unit_propagation
- http://ocw.mit.edu/courses/electrical-engineering-and-computer-science/6-005-elements-of-software-construction-fall-2011/lecture-notes/
- [Puzzle Generation](http://zhangroup.aporc.org/images/files/Paper_3485.pdf)
- http://www.sudokuwiki.org/sudoku_creation_and_grading.pdf

http://www.websudoku.com/faqs.php
http://planetsudoku.com/how-to/sudoku-squirmbag.html
http://www.sudoku-solutions.com/index.php?page=background
https://gophers.slack.com/files/mem/F0DHMJBML/top95.txt
http://www.websudoku.com/
https://en.wikipedia.org/wiki/Exact_cover#Sudoku
