package main

import (
	"fmt"

	"github.com/judwhite/go-sudoku/internal/bits"
)

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
