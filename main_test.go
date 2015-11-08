package main

import "testing"

func TestBoards(t *testing.T) {
	files := []string{
		"./test_files/input_17_clues.txt",
		"./test_files/01_naked_single_493382.txt",
		"./test_files/02_hidden_single_1053217.txt",
		"./test_files/03_naked_pair_1053222.txt",
		"./test_files/04_naked_triple_1043003.txt",
		"./test_files/05_naked_quint_1051073.txt",
		"./test_files/06_hidden_pair_1208057.txt",
		"./test_files/07_hidden_triple_188899.txt",
		"./test_files/input_supposedly_hard.txt",
		"./test_files/08_hidden_quint_188899.txt",
		"./test_files/09_pointing_pair_and_triple_1011509.txt",
		"./test_files/10_xwing_1307267.txt",
	}

	for _, file := range files {
		board, err := getBoard(file)
		if err != nil {
			t.Fatal(err)
			return
		}

		if err = board.Solve(); err != nil {
			t.Fatal(err)
			return
		}

		if !board.isSolved() {
			t.Fatalf("could not solve %s", file)
			return
		}
	}
}
