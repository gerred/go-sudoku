package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
)

func main() {
	runFile("./test_files/12_tough_20151107_173.txt") // TODO
	//runFile("./test_files/input.txt")

	//runFile("./test_files/28_xcycles.txt")

	//runTop95()
}

func runFile(fileName string) {
	board, err := getBoard(fileName)
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
		return
	}
	if err = board.Solve(); err != nil {
		fmt.Printf("ERROR - %s\n", err)
		board.PrintHints()
		return
	}

	if !board.isSolved() {
		board.PrintHints()
		board.PrintURL()
		fmt.Println("could not solve")
	} else {
		board.Print()
	}
}

func runTop95() {
	b, err := ioutil.ReadFile("./test_files/top95.txt")
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
		return
	}

	r := bufio.NewReader(bytes.NewReader(b))
	line, _ := r.ReadString('\n')
	for i := 0; line != ""; i++ {
		// TOOD: debug
		if i != 6 {
			line, _ = r.ReadString('\n')
			continue
		}

		board, err := loadBoard([]byte(line))
		if err != nil {
			board.PrintHints()
			fmt.Printf("ERROR - %d - %s\n", i+1, err)
			return
		}
		if err = board.Solve(); err != nil {
			board.PrintHints()
			fmt.Printf("ERROR - %d - %s\n", i+1, err)
			return
		}

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
}
