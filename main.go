package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/judwhite/go-sudoku/internal/bits"
)

type board struct {
	solved  [81]uint
	blits   [81]uint
	loading bool
}

func main() {

	/*for pos := 0; pos < 81; pos++ {
		startRow := ((pos / 9) / 3) * 3
		startCol := ((pos % 9) / 3) * 3
		for r := startRow; r < startRow+3; r++ {
			for c := startCol; c < startCol+3; c++ {
				fmt.Printf("%d - %d (%d, %d)\n", pos, r*9+c, r, c)
			}
		}
	}*/

	//board, err := getBoard("./test_files/12_tough_20151107_173.txt") // TODO

	//board, err getBoard("./test_files/input.txt")
	//board, err := getBoard("./test_files/input_17_clues.txt") // PASS
	//board, err := getBoard("./test_files/01_naked_single_493382.txt") // PASS
	//board, err := getBoard("./test_files/02_hidden_single_1053217.txt") // PASS
	//board, err := getBoard("./test_files/03_naked_pair_1053222.txt") // PASS
	//board, err := getBoard("./test_files/04_naked_triple_1043003.txt") // PASS
	//board, err := getBoard("./test_files/05_naked_quint_1051073.txt") // PASS
	//board, err := getBoard("./test_files/06_hidden_pair_1208057.txt") // PASS
	//board, err := getBoard("./test_files/07_hidden_triple_188899.txt") // PASS
	//board, err := getBoard("./test_files/input_supposedly_hard.txt") // PASS
	//board, err := getBoard("./test_files/08_hidden_quint_188899.txt") // PASS
	//board, err := getBoard("./test_files/09_pointing_pair_and_triple_1011509.txt") // PASS
	//board, err := getBoard("./test_files/10_xwing_1307267.txt") // PASS
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

func (b *board) Print() {
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

func (b *board) PrintHints() {
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

func (b *board) Validate() error {
	for pos := 0; pos < 81; pos++ {
		// validate row
		rowVals := make([]uint, 9)
		startRow := (pos / 9) * 9
		for r := startRow; r < startRow+9; r++ {
			rowVals[r-startRow] = b.solved[r]
		}
		if err := validate(rowVals); err != nil {
			return err
		}

		// validate column
		colVals := make([]uint, 9)
		colIndex := 0
		for c := pos % 9; c < 81; c += 9 {
			colVals[colIndex] = b.solved[c]
			colIndex++
		}
		if err := validate(colVals); err != nil {
			return err
		}

		// validate box
		startRow = ((pos / 9) / 3) * 3
		startCol := ((pos % 9) / 3) * 3
		boxVals := make([]uint, 9)
		boxIndex := 0
		for r := startRow; r < startRow+3; r++ {
			for c := startCol; c < startCol+3; c++ {
				boxVals[boxIndex] = b.solved[r*9+c]
				boxIndex++
			}
		}
		if err := validate(boxVals); err != nil {
			return err
		}
	}

	return nil
}

func validate(vals []uint) error {
	if len(vals) != 9 {
		return fmt.Errorf("len(vals) = %d", len(vals))
	}

	avail := make(map[uint]interface{})
	for i := uint(1); i <= 9; i++ {
		avail[i] = struct{}{}
	}

	for _, v := range vals {
		if v == 0 {
			continue
		}
		_, ok := avail[v]
		if !ok {
			return fmt.Errorf("val %d repeated", v)
		}
		delete(avail, v)
	}
	return nil
}

func (b *board) operateOnRow(pos int, op func(target int, source int) error) error {
	startRow := (pos / 9) * 9
	for r := startRow; r < startRow+9; r++ {
		if err := op(r, pos); err != nil {
			return err
		}
	}
	return nil
}

func (b *board) operateOnColumn(pos int, op func(target int, source int) error) error {
	for c := pos % 9; c < 81; c += 9 {
		if err := op(c, pos); err != nil {
			return err
		}
	}
	return nil
}

func (b *board) operateOnBox(pos int, op func(target int, source int) error) error {
	startRow := ((pos / 9) / 3) * 3
	startCol := ((pos % 9) / 3) * 3
	for r := startRow; r < startRow+3; r++ {
		for c := startCol; c < startCol+3; c++ {
			target := r*9 + c
			if err := op(target, pos); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *board) operateOnRCB(pos int, op func(target int, source int) error) error {
	if err := b.operateOnRow(pos, op); err != nil {
		return err
	}
	if err := b.operateOnColumn(pos, op); err != nil {
		return err
	}
	if err := b.operateOnBox(pos, op); err != nil {
		return err
	}
	return nil
}

func (b *board) operateOnCommon(pos1 int, pos2 int, op func(target int, source int) error) error {
	coords1 := getCoords(pos1)
	coords2 := getCoords(pos2)

	if coords1.row == coords2.row {
		if err := b.operateOnRow(pos1, op); err != nil {
			return err
		}
	}
	if coords1.col == coords2.col {
		if err := b.operateOnColumn(pos1, op); err != nil {
			return err
		}
	}
	if coords1.box == coords2.box {
		if err := b.operateOnBox(pos1, op); err != nil {
			return err
		}
	}
	return nil
}

func (b *board) willUpdateCandidates(targetPos int, sourcePos int, mask uint) bool {
	if targetPos == sourcePos || b.solved[targetPos] != 0 {
		return false
	}
	oldBlit := b.blits[targetPos]
	newBlit := oldBlit & mask
	if newBlit != oldBlit {
		return true
	}
	return false
}

func (b *board) updateCandidates(targetPos int, sourcePos int, mask uint) error {
	if targetPos == sourcePos || b.solved[targetPos] != 0 {
		return nil
	}
	oldBlit := b.blits[targetPos]
	newBlit := oldBlit & mask
	if newBlit != oldBlit {
		if newBlit == 0 {
			return fmt.Errorf("tried to remove last candidate from %#2v", getCoords(targetPos))
		}

		b.blits[targetPos] = newBlit
		delta := oldBlit & ^newBlit
		b.Log(false, targetPos, fmt.Sprintf("old hints: %-10s remove hint: %s remaining hints: %s", bits.GetString(oldBlit), bits.GetString(delta), bits.GetString(newBlit)))
	}
	return nil
}

func (b *board) Log(isSolve bool, pos int, msg string) {
	if b.loading {
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

func (b *board) SolvePosition(pos int, val uint) error {
	mask := uint(^(1 << (val - 1)))
	if b.solved[pos] != 0 /*&& b.solved[pos] != val*/ {
		return fmt.Errorf("pos %d has value %d, tried to set with %d", pos, b.solved[pos], val)
	}
	b.solved[pos] = val
	b.blits[pos] = 1 << (val - 1)

	b.Log(true, pos, fmt.Sprintf("set value %d mask:%09b", val, mask&0x1FF))

	if !b.loading {
		b.Print()
		b.PrintHints()
	}

	if err := b.Validate(); err != nil {
		return fmt.Errorf("%#v val:%d - %s", getCoords(pos), val, err)
	}

	if err := b.operateOnRCB(pos, b.removeCandidates(mask)); err != nil {
		return err
	}

	if !b.loading {
		b.PrintHints()
	}

	return nil
}

func (b *board) removeCandidates(mask uint) func(int, int) error {
	return func(target int, source int) error {
		if opErr := b.updateCandidates(target, source, mask); opErr != nil {
			return opErr
		}
		return nil
	}
}

func getBoard(fileName string) (*board, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	b = bytes.Replace(b, []byte{'\r'}, []byte{}, -1)
	b = bytes.Replace(b, []byte{'\n'}, []byte{' '}, -1)
	b = bytes.Replace(b, []byte{' ', ' '}, []byte{' '}, -1)

	board := &board{loading: true}
	for i := 0; i < 81; i++ {
		board.blits[i] = 0x1FF
	}

	pos := 0
	for i := 0; i < 162; i += 2 {
		if b[i] != '_' {
			val := uint(b[i] - 48)
			if err = board.SolvePosition(pos, val); err != nil {
				return board, err
			}
		}
		pos++
	}

	board.Print()
	board.PrintHints()

	board.loading = false

	return board, nil
}

func (b *board) numSolved() int {
	num := 0
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			num++
		}
	}
	return num
}

func (b *board) isSolved() bool {
	return b.numSolved() == 81
}

func (b *board) Solve() error {
	doLoop := true

	// first iteration naked single
	b.loading = true // turn off logging, this run is boring
	fmt.Println("--- NAKED SINGLE: FIRST ITERATION")
	if err := b.SolveNakedSingle(); err != nil {
		return err
	}
	b.PrintHints()
	b.loading = false

	for i := 0; doLoop; i++ {
		oldBlits := b.blits
		if err := b.runSolvers(); err != nil {
			return err
		}
		doLoop = !reflect.DeepEqual(oldBlits, b.blits)
		if b.isSolved() {
			break
		}
		fmt.Printf("doLoop: %t\n", doLoop)

		if !doLoop {
			// now we branch
			// TODO
		}
	}
	return nil
}

func (b *board) runSolvers() error {
	var err error
	fmt.Println("--- NAKED SINGLE")
	if err = b.SolveNakedSingle(); err != nil {
		return err
	}
	fmt.Println("--- HIDDEN SINGLE")
	if err = b.SolveHiddenSingle(); err != nil {
		return err
	}
	fmt.Println("--- NAKED PAIR")
	if err = b.SolveNakedN(2); err != nil {
		return err
	}
	fmt.Println("--- NAKED TRIPLE")
	if err = b.SolveNakedN(3); err != nil {
		return err
	}
	fmt.Println("--- NAKED QUAD")
	if err = b.SolveNakedN(4); err != nil {
		return err
	}
	/*fmt.Println("--- NAKED QUINT")
	if err = b.SolveNakedN(5); err != nil {
		return err
	}*/
	fmt.Println("--- HIDDEN PAIR")
	if err = b.SolveHiddenN(2); err != nil {
		return err
	}
	fmt.Println("--- HIDDEN TRIPLE")
	if err = b.SolveHiddenN(3); err != nil {
		return err
	}
	/*fmt.Println("--- HIDDEN QUAD")
	if err = b.SolveHiddenN(4); err != nil {
		return err
	}
	fmt.Println("--- HIDDEN QUINT")
	if err = b.SolveHiddenN(5); err != nil {
		return err
	}*/
	fmt.Println("--- POINTING PAIR AND TRIPLE REDUCTION")
	if err = b.SolvePointingPairAndTripleReduction(); err != nil {
		return err
	}
	fmt.Println("--- BOX LINE")
	if err = b.SolveBoxLine(); err != nil {
		return err
	}

	fmt.Println("--- XWING")
	if err = b.SolveXWing(); err != nil {
		return err
	}

	// TODO: remove, temp to cleanup SWORDFISH candidates
	fmt.Println("--- NAKED SINGLE")
	if err = b.SolveNakedSingle(); err != nil {
		return err
	}
	fmt.Println("--- HIDDEN SINGLE")
	if err = b.SolveHiddenSingle(); err != nil {
		return err
	}

	fmt.Println("*****************************************")
	b.PrintHints()
	fmt.Println("*****************************************")
	// END TODO

	fmt.Println("--- SWORDFISH")
	if err = b.SolveSwordFish(); err != nil {
		return err
	}

	fmt.Println("--- END SOLVER")

	return nil
}

func (b *board) SolveNakedSingle() error {
	// Naked Single - only hint left
	doLoop := true
	for doLoop {
		doLoop = false
		for i := 0; i < 81; i++ {
			if b.solved[i] != 0 {
				continue
			}

			blit := b.blits[i]
			if !bits.HasSingleBit(blit) {
				continue
			}

			num := bits.GetSingleBitValue(blit)

			if err := b.SolvePosition(i, uint(num)); err != nil {
				return err
			}
			doLoop = true
		}
	}

	return nil
}

func (b *board) SolveHiddenSingle() error {
	// Hidden Single - a given cell contains a candidate which is only
	// present in this cell and not in the rest of the row/column/box
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}
		blit := b.blits[i]

		var sumBlits uint
		sumHints := func(target int, source int) error {
			if target == source {
				return nil
			}
			sumBlits |= b.blits[target]

			return nil
		}

		ops := []func(int, func(int, int) error) error{
			b.operateOnRow,
			b.operateOnColumn,
			b.operateOnBox,
		}

		for opIt, op := range ops {
			sumBlits = 0
			if err := op(i, sumHints); err != nil {
				return err
			}
			leftOver := blit & ^sumBlits

			if bits.HasSingleBit(leftOver) {
				val := bits.GetSingleBitValue(leftOver)
				fmt.Printf("op-it:%d c:%#2v h:%09b sh:%09b ^sh:%09b lo:%b\n", opIt, getCoords(i), blit, sumBlits, ^sumBlits&0x1FF, leftOver)
				if err := b.SolvePosition(i, val); err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

type coords struct {
	row int
	col int
	box int
}

func getCoords(pos int) coords {
	boxRow := ((pos / 9) / 3)
	boxCol := ((pos % 9) / 3)
	box := boxRow*3 + boxCol

	return coords{row: pos / 9, col: pos % 9, box: box}
}

func (b *board) SolveNakedN(n int) error {
	if n < 2 || n > 5 {
		return fmt.Errorf("n must be between [2,5], actual=%d", n)
	}
	// When a cell has N candidates and (N-1) others combined equal
	// the N candidates, then all N candidates can be removed from
	// the other cells in common.
	// http://planetsudoku.com/how-to/sudoku-naked-triple.html
	// http://planetsudoku.com/how-to/sudoku-naked-quad.html
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}
		if bits.GetNumberOfSetBits(b.blits[i]) > uint(n) {
			continue
		}

		ops := []func(int, func(int, int) error) error{
			b.operateOnRow,
			b.operateOnColumn,
			b.operateOnBox,
		}

		for _, op := range ops {
			var pickList []int
			addToPickList := func(target int, source int) error {
				if target == source || b.solved[target] != 0 {
					return nil
				}
				pickList = append(pickList, target)
				return nil
			}

			if err := op(i, addToPickList); err != nil {
				return err
			}

			if len(pickList) <= n {
				continue
			}

			perms := b.getPermutations(n, pickList, []int{i})
			for _, list := range perms {
				var blit uint
				for _, item := range list {
					blit |= b.blits[item]
				}

				if bits.GetNumberOfSetBits(blit) != uint(n) {
					continue
				}

				removeHints := func(target int, source int) error {
					for _, item := range list {
						if item == target {
							return nil
						}
					}

					return b.updateCandidates(target, source, ^blit)
				}

				if err := op(i, removeHints); err != nil {
					return err
				}

				/*fmt.Printf("list %#2v\n", getCoords(i))
				for _, item := range list {
					fmt.Printf("- %#2v\n", getCoords(item))
				}

				// TODO*/
			}
		}
	}
	return nil
}

func (b *board) SolveHiddenN(n int) error {
	if n < 2 || n > 5 {
		return fmt.Errorf("n must be between [2,5], actual=%d", n)
	}
	// If there are N unique hints in N cells within one container,
	// then no other hints could be valid within that container.
	// http://planetsudoku.com/how-to/sudoku-hidden-triple.html
	// Triple example:
	// - N = 3,4,7,8
	// - X = 1,4,5,6,8
	// - Y = 4,5,6,7
	// - Other cells contain 1,3,4,5
	// Algo:
	// - bits.GetNumberOfSetBits(N | X | Y) >= 3
	// - bits.GetNumberOfSetBits((N | X | Y) & ^(O1 | O2 | O3 ...) == 3
	// - N,X,Y can &= ^(sum), removing 1,4,5
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}

		var pickList []int

		storePickList := func(target int, source int) error {
			if target == source || b.solved[target] != 0 {
				return nil
			}
			pickList = append(pickList, target)
			return nil
		}

		ops := []func(int, func(int, int) error) error{
			b.operateOnRow,
			b.operateOnColumn,
			b.operateOnBox,
		}

		for _, op := range ops {
			pickList = make([]int, 0)
			if err := op(i, storePickList); err != nil {
				return err
			}
			lists := b.getPermutations(n, pickList, []int{i})
			if err := b.checkHiddenPermutations(n, i, op, lists); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *board) checkHiddenPermutations(n int, source int, op func(int, func(int, int) error) error, lists [][]int) error {
	for _, list := range lists {
		var sumBits uint
		for _, pos := range list {
			sumBits |= b.blits[pos]
		}
		if bits.GetNumberOfSetBits(sumBits) < uint(n) {
			continue
		}

		sumOthers := uint(0)
		sumTheOthers := func(target int, source int) error {
			if b.solved[target] != 0 {
				return nil
			}
			for _, v := range list {
				if v == target {
					return nil
				}
			}
			sumOthers |= b.blits[target]
			return nil
		}

		if err := op(source, sumTheOthers); err != nil {
			return err
		}

		if sumOthers == 0 {
			continue
		}

		leftOver := (sumBits ^ sumOthers) & sumBits

		if bits.GetNumberOfSetBits(leftOver) == uint(n) {
			/*fmt.Printf("HIDDEN %d\n", n)
			for _, pos := range list {
				fmt.Printf("- %#2v %09b\n", getCoords(pos), b.blits[pos])
			}
			fmt.Printf("- sum:        %09b\n", sumBits)
			fmt.Printf("- sum others: %09b\n", sumOthers)
			fmt.Printf("- left over:  %09b\n", leftOver)*/

			for _, pos := range list {
				if err := b.updateCandidates(pos, source, leftOver); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (b *board) getPermutations(n int, pickList []int, curList []int) [][]int {
	output := make([][]int, 0)

	for i := 0; i < len(pickList); i++ {
		list := make([]int, len(curList))
		copy(list, curList)              // get the source list
		list = append(list, pickList[i]) // plus the current element

		if len(list) == n {
			// if this is the length we're looking for...
			output = append(output, list)
		} else {
			// otherwise, call recursively
			perms := b.getPermutations(n, pickList[i+1:], list)
			if perms != nil {
				for _, v := range perms {
					output = append(output, v)
				}
			}
		}
	}

	return output
}

func (b *board) SolvePointingPairAndTripleReduction() error {
	// http://planetsudoku.com/how-to/sudoku-pointing-pair-and-triple.html
	// "I have two or three unique HINTS within a shared box, sharing the same
	// ROW or COLUMN. Therefore that hint cannot belong anywhere else on that
	// ROW or COLUMN in any other BOXES".

	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}

		coords := getCoords(i)

		dims := []struct {
			isRow bool
			op    func(int, func(int, int) error) error
		}{
			{isRow: true, op: b.operateOnRow},
			{isRow: false, op: b.operateOnColumn},
		}

		for _, dim := range dims {
			var pickList []int
			var negateList []int
			getPickList := func(target int, source int) error {
				if target == source || b.solved[target] != 0 {
					return nil
				}

				testCoords := getCoords(target)

				if (dim.isRow && coords.row == testCoords.row) ||
					(!dim.isRow && coords.col == testCoords.col) {
					pickList = append(pickList, target)
				} else {
					negateList = append(negateList, target)
				}
				return nil
			}

			if err := b.operateOnBox(i, getPickList); err != nil {
				return err
			}

			sumNegateBits := uint(0)
			for _, item := range negateList {
				sumNegateBits |= b.blits[item]
			}

			for x := 3; x >= 2; x-- {
				perms := b.getPermutations(x, pickList, []int{i})
				sumBits := uint(0)
				for _, list := range perms {
					for _, item := range list {
						sumBits |= b.blits[item]
					}

					leftOver := sumBits & ^sumNegateBits
					nbits := bits.GetNumberOfSetBits(leftOver)
					if nbits != 2 && nbits != 3 {
						continue
					}

					//exit := false
					removeHints := func(target int, source int) error {
						if b.solved[target] != 0 {
							return nil
						}
						testCoords := getCoords(target)
						if testCoords.box == coords.box {
							return nil
						}
						/*if b.willUpdateCandidates(target, source, ^leftOver) {
							if !exit {
								fmt.Printf("i: %#2v\n", coords)
								for _, item := range list {
									fmt.Printf(" - %#2v %09b %s\n", getCoords(item), b.blits[item], bits.GetString(b.blits[item]))
								}

								fmt.Printf(" - sum:        %09b\n", sumBits)
								fmt.Printf(" - negate sum: %09b\n", sumNegateBits)
								fmt.Printf(" - left over:  %09b\n", leftOver)
							}
							fmt.Printf("-> %#2v\n", testCoords)
							exit = true
						}*/
						return b.updateCandidates(target, source, ^leftOver)
					}

					if err := dim.op(i, removeHints); err != nil {
						return err
					}

					/*if exit {
						b.PrintHints()
						os.Exit(0)
					}*/
				}
			}
		}
	}
	return nil
}

func (b *board) SolveBoxLine() error {
	// Two cells in a BOX that share a hint which isn't anywhere else on
	// the ROW or COLUMN they share can be removed as hints from other cells
	// in the same BOX.
	// http://planetsudoku.com/how-to/sudoku-box-line.html
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}

		blit := b.blits[i]
		coords := getCoords(i)

		var colPickList []int
		var rowPickList []int

		savePickLists := func(target int, source int) error {
			if target == source || b.solved[target] != 0 {
				return nil
			}
			if b.blits[target]&blit == 0 {
				// nothing shared
				return nil
			}

			targetCoords := getCoords(target)
			if targetCoords.row == coords.row {
				rowPickList = append(rowPickList, target)
			} else if targetCoords.col == coords.col {
				colPickList = append(colPickList, target)
			}
			return nil
		}

		if err := b.operateOnBox(i, savePickLists); err != nil {
			return err
		}

		dims := []struct {
			pickList  []int
			op        func(int, func(int, int) error) error
			canRemove func(int) bool
		}{
			// rows
			{
				pickList:  rowPickList,
				op:        b.operateOnRow,
				canRemove: func(target int) bool { return coords.row != getCoords(target).row },
			},
			// columns
			{
				pickList:  colPickList,
				op:        b.operateOnColumn,
				canRemove: func(target int) bool { return coords.col != getCoords(target).col },
			},
		}

		for _, dim := range dims {
			perms := b.getPermutations(1, dim.pickList, []int{})
			for _, list := range perms {
				for _, item := range list {
					sharedHints := blit & b.blits[item]
					hintList := bits.GetBitList(sharedHints)

					for _, hint := range hintList {
						safeToRemove := true
						checkLine := func(target int, source int) error {
							if target == source {
								return nil
							}
							if getCoords(target).box == coords.box {
								return nil
							}
							if b.blits[target]&hint != 0 {
								safeToRemove = false
							}
							return nil
						}

						if err := dim.op(i, checkLine); err != nil {
							return err
						}

						if safeToRemove {
							removeBoxLineHint := func(target int, source int) error {
								if !dim.canRemove(target) {
									return nil
								}

								return b.updateCandidates(target, i, ^hint)
							}

							if err := b.operateOnBox(i, removeBoxLineHint); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (b *board) SolveXWing() error {
	// When there are
	// - only two possible cells for a value in each of two different rows,
	// - and these candidates lie also in the same columns,
	// - then all other candidates for this value in the columns can be eliminated.
	// The reverse is also true for 2 columns with 2 common rows.
	// http://www.sudokuwiki.org/x_wing_strategy
	// http://planetsudoku.com/how-to/sudoku-x-wing.html
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}

		// TODO: remove, temp to cleanup XWING candidates
		if err := b.SolveNakedSingle(); err != nil {
			return err
		}
		if err := b.SolveHiddenSingle(); err != nil {
			return err
		}
		// END TODO

		blit := b.blits[i]
		c1 := getCoords(i)

		dims := []struct {
			op        func(int, func(int, int) error) error
			op2       func(int, func(int, int) error) error
			isAligned func(coords, coords) bool
		}{
			{
				op:        b.operateOnRow,
				op2:       b.operateOnColumn,
				isAligned: func(c1 coords, c2 coords) bool { return c1.row == c2.row },
			},
			{
				op:        b.operateOnColumn,
				op2:       b.operateOnRow,
				isAligned: func(c1 coords, c2 coords) bool { return c1.col == c2.col },
			},
		}

		for _, dim := range dims {
			bitList := bits.GetBitList(blit)
			for _, bit := range bitList {
				findPairs := func(items *[]int) func(target int, source int) error {
					return func(target int, source int) error {
						if target == source {
							return nil
						}
						if b.blits[target]&bit == bit {
							*items = append(*items, target)
						}
						return nil
					}
				}

				var pairs []int
				if err := dim.op(i, findPairs(&pairs)); err != nil {
					return err
				}

				if len(pairs) != 1 {
					continue
				}
				lockedPairPos := pairs[0]
				c2 := getCoords(lockedPairPos)

				var pairs21 []int
				if err := dim.op2(i, findPairs(&pairs21)); err != nil {
					return err
				}

				var pairs22 []int
				if err := dim.op2(lockedPairPos, findPairs(&pairs22)); err != nil {
					return err
				}

				// TODO: item21/item22 must only cell with hin in their shared row/column
				for _, item21 := range pairs21 {
					c3 := getCoords(item21)

					// ensure value lives in container only twice, pairs are locked
					var pairs2 []int
					if err := dim.op(item21, findPairs(&pairs2)); err != nil {
						return err
					}

					if len(pairs2) != 1 {
						continue
					}

					var shortList []int
					for _, item22 := range pairs22 {
						c4 := getCoords(item22)
						if dim.isAligned(c3, c4) {
							shortList = append(shortList, item22)
						}
					}

					if len(shortList) != 1 {
						continue
					}

					for _, item22 := range shortList {
						c4 := getCoords(item22)

						logged := false

						sourceList := []int{i, lockedPairPos, item21, item22}

						removeHints := func(target int, source int) error {
							for _, pos := range sourceList {
								if target == pos {
									return nil
								}
							}

							if b.willUpdateCandidates(target, source, ^bit) && !logged {
								logged = true
								b.PrintHints()
								fmt.Printf("xwing: val:%d\n", bits.GetSingleBitValue(bit))
								fmt.Printf("- %#2v\n", c1)
								fmt.Printf("- %#2v\n", c2)
								fmt.Printf("- %#2v\n", c3)
								fmt.Printf("- %#2v\n", c4)
							}

							return b.updateCandidates(target, source, ^bit)
						}

						for _, pos := range []int{i, lockedPairPos} {
							if err := dim.op2(pos, removeHints); err != nil {
								return err
							}
						}
					}
				}
			}
		}

	}
	return nil
}

func (b *board) SolveSwordFish() error {
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}
	}
	return nil
}
