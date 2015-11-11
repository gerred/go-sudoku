package main

import (
	"fmt"

	"github.com/judwhite/go-sudoku/internal/bits"
)

func (b *board) Print() {
	if b.quiet {
		return
	}

	fmt.Print("|-------|-------|-------|\n| ")
	for i := 0; i < len(b.solved); i++ {
		if b.solved[i] == 0 {
			fmt.Print("_ ")
		} else {
			fmt.Printf("%d ", b.solved[i])
		}
		if (i+1)%9 == 0 {
			fmt.Print("|\n|")
			if (i+1)%27 == 0 {
				fmt.Print("-------|-------|-------|\n")
				if i != 80 {
					fmt.Print("| ")
				}
			} else {
				fmt.Print(" ")
			}
		} else if (i+1)%3 == 0 {
			fmt.Print("| ")
		}
	}
}

func (b *board) PrintURL() {
	if b.quiet {
		return
	}

	for i := 0; i < len(b.solved); i++ {
		fmt.Print(b.solved[i])
	}
	fmt.Println()
}

func (b *board) PrintHints() {
	if b.quiet {
		return
	}

	fmt.Print("|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|\n")
	fmt.Printf("|r,c| %15d %15d %15d | %15d %15d %15d | %15d %15d %15d |\n", 0, 1, 2, 3, 4, 5, 6, 7, 8)
	fmt.Print("|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|\n| 0 | ")
	for i := 0; i < len(b.solved); i++ {
		if b.solved[i] == 0 {
			fmt.Printf("%15s ", fmt.Sprintf("(%s)", bits.GetString(b.blits[i])))
		} else {
			fmt.Printf("%15d ", b.solved[i])
		}
		if (i+1)%9 == 0 {
			fmt.Printf("|\n|")
			if (i+1)%27 == 0 {
				fmt.Print("---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|\n")
				if i != 80 {
					fmt.Printf("| %d | ", (i+1)/9)
				}
			} else {
				fmt.Printf(" %d | ", (i+1)/9)
			}
		} else if (i+1)%3 == 0 {
			fmt.Print("| ")
		}
	}
}

func (b *board) Log(isSolve bool, pos int, msg string) {
	if b.loading || b.quiet {
		return
	}

	var prefix string
	if isSolve {
		prefix = ">"
	} else {
		prefix = "-"
	}

	coords := getCoords(pos)
	fmt.Printf("%s R%dC%d: %s\n", prefix, coords.row, coords.col, msg)
}
