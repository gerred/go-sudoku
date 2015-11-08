package main

import "fmt"

func main() {

	//board, err := getBoard("./test_files/12_tough_20151107_173.txt") // TODO

	//board, err getBoard("./test_files/input.txt")
	board, err := getBoard("./test_files/11_swordfish_1280430.txt")
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
		return
	}
	if err = board.Solve(); err != nil {
		fmt.Printf("ERROR - %s\n", err)
		board.Print()
		board.PrintHints()
		return
	}

	if !board.isSolved() {
		board.PrintHints()
	}
	board.Print()
}
