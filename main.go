package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/judwhite/go-sudoku/internal/sat"
)

func main() {

	//runFile("./test_files/12_tough_20151107_173.txt") // TODO
	//runFile("./test_files/input.txt")

	//runFile("./test_files/28_xcycles.txt")

	runTop95()
}

func (b *board) SolveSAT() error {
	satInput := b.getSAT()
	satSolver, err := sat.NewSAT(satInput)
	if err != nil {
		return err
	}
	satSolver = satSolver.Solve()
	if satSolver == nil {
		return fmt.Errorf("could not solve with SAT\n")
	} else {
		fmt.Printf("solved with SAT\n")

		for _, setvar := range satSolver.SetVars {
			k := setvar.VarNum
			v := setvar.Value
			if v {
				//fmt.Printf("%d %v\n", k, v)
				r := k/100 - 1
				c := (k%100)/10 - 1
				pos := r*9 + c
				if b.solved[pos] == 0 {
					val := k % 10
					fmt.Printf("r:%d c:%d val:%d\n", r, c, val)
					err := b.SolvePosition(pos, uint(val))
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func runFile(fileName string) {
	board, err := getBoard(fileName)
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
		return
	}

	err = board.SolveSAT()
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
		return
	}

	/*if err = board.Solve(); err != nil {
		fmt.Printf("ERROR - %s\n", err)
		board.PrintHints()
		return
	}*/

	if !board.isSolved() {
		board.PrintHints()
		board.PrintURL()
		fmt.Println("could not solve")
	} else {
		board.Print()
	}
}

func runTop95() {
	start := time.Now()
	b, err := ioutil.ReadFile("./test_files/top95.txt")
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
		return
	}

	r := bufio.NewReader(bytes.NewReader(b))
	line, _ := r.ReadString('\n')
	for i := 0; line != ""; i++ {
		board, err := loadBoard([]byte(line))
		if err != nil {
			board.PrintHints()
			fmt.Printf("ERROR - %d - %s\n", i+1, err)
			return
		}

		if err = board.SolveSAT(); err != nil {
			fmt.Printf("ERROR - %s\n", err)
			return
		}

		/*if err = board.Solve(); err != nil {
			board.PrintHints()
			fmt.Printf("ERROR - %d - %s\n", i+1, err)
			return
		}*/

		if !board.isSolved() {
			board.PrintHints()
			board.PrintURL()
			fmt.Printf("could not solve - %d\n", i+1)
			break
		} else {
			board.Print()
		}

		line, _ = r.ReadString('\n')
	}
	fmt.Printf("%v\n", time.Since(start))
}
