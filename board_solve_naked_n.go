package main

import (
	"fmt"

	"github.com/judwhite/go-sudoku/internal/bits"
)

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

		ops := []containerOperator{
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

			perms := getPermutations(n, pickList, []int{i})
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
			}
		}
	}
	return nil
}
