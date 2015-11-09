package main

import (
	"fmt"
	"testing"
)

func TestBoards(t *testing.T) {
	files := []string{
		"./test_files/01_naked_single_493382.txt",
		"./test_files/02_hidden_single_1053217.txt",
		"./test_files/03_naked_pair_1053222.txt",
		"./test_files/04_naked_triple_1043003.txt",
		"./test_files/05_naked_quint_1051073.txt",
		"./test_files/06_hidden_pair_1208057.txt",
		"./test_files/07_hidden_triple_188899.txt",
		"./test_files/08_hidden_quint_188899.txt",
		"./test_files/09_pointing_pair_and_triple_1011509.txt",
		"./test_files/10_xwing_1307267.txt",
		//"./test_files/12_tough_20151107_173.txt", // TODO
		"./test_files/11_swordfish_1280430.txt",
		"./test_files/13_swordfish_008009000300057001000100009230000070005406100060000038900003000700840003000700600.txt",
		"./test_files/14_swordfish_980010020002700000000009010700040800600107002009030005040900000000005700070020039.txt",
		"./test_files/15_swordfish_108000067000050000000000030006100040450000900000093000200040010003002700807001005.txt",
		"./test_files/16_swordfish_107300040800006000050870630090000510000000007700060080000904000080100002410000000.txt",
		"./test_files/17_swordfish_300040000000007048000000907010003080400050020050008070500300000000000090609025300.txt",
		"./test_files/18_swordfish.txt",
		"./test_files/19_supposedly_hard.txt",
		"./test_files/20_17_clues.txt",
		"./test_files/21_ywing.txt",
		"./test_files/22_xychain.txt",
	}

	for _, file := range files {
		board, err := getBoard(file)
		if err != nil {
			t.Fatalf("%s: %s", file, err)
			return
		}

		if err = board.Solve(); err != nil {
			t.Fatalf("%s: %s", file, err)
			return
		}

		if !board.isSolved() {
			t.Fatalf("%s: could not solve", file)
			return
		}
	}
	fmt.Printf("solved %d puzzles\n", len(files))
}
