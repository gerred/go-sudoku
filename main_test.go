package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/judwhite/go-sudoku/internal/bits"
	"github.com/judwhite/go-sudoku/internal/sat"
)

func testXCycles(t *testing.T) {
	// arrange
	hintBoard := `
|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
|r,c|               0               1               2 |               3               4               5 |               6               7               8 |
|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
| 0 |               8           (1,9)               4 |               5               3               7 |         (1,6,9)         (1,2,6)           (1,2) |
| 1 |           (7,9)               2               3 |               6               1               4 |           (7,9)               8               5 |
| 2 |               6           (1,7)               5 |               9               8               2 |           (1,7)               3               4 |
|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
| 3 |         (3,4,9)         (3,4,6)         (2,6,9) |               1         (4,6,9)               5 |               8               7           (2,9) |
| 4 |               5           (4,9)           (1,2) |               7           (4,9)               8 |               3           (1,2)               6 |
| 5 |         (1,7,9)               8       (1,6,7,9) |               2           (6,9)               3 |               4               5           (1,9) |
|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
| 6 |               2         (4,6,7)         (1,6,7) |               8               5               9 |           (1,6)         (1,4,6)               3 |
| 7 |           (4,9)               5           (6,9) |               3               7               1 |               2         (4,6,9)               8 |
| 8 |         (1,3,9)           (3,9)               8 |               4               2               6 |               5           (1,9)               7 |
|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
`

	b := loadBoardWithHints(t, hintBoard)

	b.PrintHints()
	b.PrintURL()

	newBoard, err := TrialAndError(*b)
	if err != nil {
		t.Fatal(err)
	}
	if newBoard != nil {
		newBoard.Print()
	}

	// check board is in expected initial state
	testHint(t, b, 1, 0, []uint{7, 9})
	testHint(t, b, 8, 0, []uint{1, 3, 9})
	testHint(t, b, 8, 7, []uint{1, 9})
	testHint(t, b, 8, 1, []uint{3, 9})
	testHint(t, b, 4, 1, []uint{4, 9})
	testHint(t, b, 0, 1, []uint{1, 9})

	// TODO: testHint on correct cells in cycle

	// act
	if err := b.SolveXCycles(); err != nil {
		t.Fatal(err)
	}

	b.PrintHints()

	// assert

	// test for absence of defect in X-CYCLES:
	testHint(t, b, 3, 0, []uint{3, 4, 9})
	testHint(t, b, 5, 0, []uint{1, 7, 9})
	testHint(t, b, 7, 0, []uint{4, 9})

	// test for expected state after X-CYCLES applied
	// note: this test may fail in the future if X-CYCLES is modified
	//       since it could pick a different cycle to operate on
	/*testHint(t, b, 3, 1, []uint{3})
	testHint(t, b, 5, 1, []uint{3, 6})
	testHint(t, b, 8, 2, []uint{2, 9})*/
}

func testXYChain(t *testing.T) {
	// arrange
	hintBoard := `
|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
|r,c|               0               1               2 |               3               4               5 |               6               7               8 |
|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
| 0 |               4               8               7 |               3               1               2 |           (5,6)               9           (5,6) |
| 1 |           (5,9)           (5,9)               3 |               6           (4,8)           (4,8) |               2               7               1 |
| 2 |               1               2               6 |           (5,7)               9           (5,7) |               3               8               4 |
|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
| 3 |               7           (3,4)               5 |           (8,9)         (3,4,8)         (4,8,9) |               1               6               2 |
| 4 |           (6,9)               1           (4,9) |               2         (3,4,6)           (5,7) |               8           (3,4)           (5,7) |
| 5 |           (2,8)         (3,4,6)           (2,8) |           (5,7)         (3,4,6)               1 |           (5,7)           (3,4)               9 |
|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
| 6 |           (5,8)           (4,5)               1 |           (4,8)               7               6 |               9               2               3 |
| 7 |               3           (6,7)           (8,9) |               1               2           (8,9) |               4               5           (6,7) |
| 8 |           (2,6)       (4,6,7,9)         (2,4,9) |           (4,9)               5               3 |           (6,7)               1               8 |
|---|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
	`
	b := loadBoardWithHints(t, hintBoard)
	b.quiet = true

	// check board is in expected initial state
	testHint(t, b, 5, 6, []uint{5, 7})
	testHint(t, b, 3, 1, []uint{3, 4})
	testHint(t, b, 5, 1, []uint{3, 4, 6})
	testHint(t, b, 8, 2, []uint{2, 4, 9})

	// act
	if err := b.SolveXYChain(); err != nil {
		t.Fatal(err)
	}

	// assert

	// test for absence of defect in XY-Chain:
	// - R5C6: old hints: 5,7        remove hint: 5 remaining hints: 7
	testHint(t, b, 5, 6, []uint{5, 7})

	// test for expected state after XY-Chain applied
	// note: this test may fail in the future if XY-Chain is modified
	//       since it could pick a different chain to operate on
	// - R3C1: old hints: 3,4        remove hint: 4 remaining hints: 3
	// - R5C1: old hints: 3,4,6      remove hint: 4 remaining hints: 3,6
	// - R8C2: old hints: 2,4,9      remove hint: 4 remaining hints: 2,9
	testHint(t, b, 3, 1, []uint{3})
	testHint(t, b, 5, 1, []uint{3, 6})
	testHint(t, b, 8, 2, []uint{2, 9})
}

func testHint(t *testing.T, b *board, row, col int, hints []uint) {
	actual := b.blits[row*9+col]
	var expected uint
	for _, hint := range hints {
		expected |= 1 << (hint - 1)
	}
	if expected != actual {
		t.Fatalf("R%dC%d, expected %v actual %v", row, col, hints, bits.GetString(actual))
	}
}

func loadBoardWithHints(t *testing.T, hintBoard string) (b *board) {
	// read the text board, apply hints
	sr := strings.NewReader(hintBoard)
	r := bufio.NewReader(sr)

	// skip header
	for i := 0; i < 4; i++ {
		r.ReadString('\n')
	}

	b = &board{}

	for i := 0; i < 9; i++ {
		line, _ := r.ReadString('\n')
		if strings.HasPrefix(line, "|---|") {
			line, _ = r.ReadString('\n')
		}
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		line = line[6 : len(line)-2]
		line = strings.Replace(line, " |", "", -1)

		start := 0
		cells := make([]string, 9)
		for j := 0; j < 9; j++ {
			end := start + 15
			cell := strings.Trim(line[start:end], " ")
			cells[j] = cell
			start = end + 1

			pos := i*9 + j
			if strings.HasPrefix(cell, "(") {
				// get hints
				hints := strings.Split(cell[1:len(cell)-1], ",")
				for _, hint := range hints {
					val, _ := strconv.Atoi(hint)
					b.blits[pos] |= 1 << uint(val-1)
				}
			} else {
				// solved cell
				val, _ := strconv.Atoi(cell)
				b.solved[pos] = uint(val)
				b.blits[pos] = 1 << uint(val-1)
			}
		}
	}

	return b
}

func TestBoards(t *testing.T) {
	files := []string{
		"./test_files/input.txt",
		//"./test_files/01_naked_single_493382.txt",
		/*"./test_files/02_hidden_single_1053217.txt",
		"./test_files/03_naked_pair_1053222.txt",
		"./test_files/04_naked_triple_1043003.txt",
		"./test_files/05_naked_quint_1051073.txt",
		"./test_files/06_hidden_pair_1208057.txt",
		"./test_files/07_hidden_triple_188899.txt",
		"./test_files/08_hidden_quint_188899.txt",
		"./test_files/09_pointing_pair_and_triple_1011509.txt",
		"./test_files/10_xwing_1307267.txt",*/
		//"./test_files/12_tough_20151107_173.txt", // TODO
		/*"./test_files/11_swordfish_1280430.txt",
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
		"./test_files/23_xychain.txt",
		"./test_files/24_xychain.txt",
		"./test_files/25_xychain.txt",
		"./test_files/26_xychain.txt",*/
	}

	for _, file := range files {
		board, err := getBoard(file)
		if err != nil {
			t.Fatalf("%s: %s", file, err)
			return
		}
		board.PrintURL()

		satInput := board.getSAT()
		//fmt.Println(satInput)
		sat, err := sat.NewSAT(satInput)
		if err != nil {
			t.Fatalf("%s: %s", file, err)
			return
		}
		sat = sat.Solve()
		if sat == nil {
			t.Fatalf("%s: could not solve with SAT", file)
		} else {
			fmt.Printf("%s: solved with SAT\n", file)

			for _, setvar := range sat.SetVars {
				k := setvar.VarNum
				v := setvar.Value
				if v {
					fmt.Printf("%d %v\n", k, v)
					r := k/100 - 1
					c := (k%100)/10 - 1
					pos := r*9 + c
					if board.solved[pos] == 0 {
						val := k % 10
						fmt.Printf("r:%d c:%d val:%d\n", r, c, val)
						err := board.SolvePosition(pos, uint(val))
						if err != nil {
							t.Fatalf("%s: %s", file, err)
							return
						}
					}
				}
			}

			//fmt.Println(satInput)

			board.PrintHints()
		}

		/*if err = board.Solve(); err != nil {
			t.Fatalf("%s: %s", file, err)
			return
		}*/

		if !board.isSolved() {
			t.Fatalf("%s: could not solve", file)
			return
		}
	}
	fmt.Printf("solved %d puzzles\n", len(files))
}
