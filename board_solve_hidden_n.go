package main

import "fmt"

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

		ops := []containerOperator{
			b.operateOnRow,
			b.operateOnColumn,
			b.operateOnBox,
		}

		for _, op := range ops {
			pickList = make([]int, 0)
			if err := op(i, storePickList); err != nil {
				return err
			}
			lists := getPermutations(n, pickList, []int{i})
			if err := b.checkHiddenPermutations(n, i, op, lists); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *board) checkHiddenPermutations(n int, source int, op containerOperator, lists [][]int) error {
	for _, list := range lists {
		var sumBits uint
		for _, pos := range list {
			sumBits |= b.blits[pos]
		}
		if GetNumberOfSetBits(sumBits) < uint(n) {
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

		if GetNumberOfSetBits(leftOver) == uint(n) {
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
