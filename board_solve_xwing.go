package main

import "github.com/judwhite/go-sudoku/internal/bits"

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

		blit := b.blits[i]
		//c1 := getCoords(i)

		// dims: operations using a row or column view as the starting point,
		//       the inverse being the elimination op
		dims := []struct {
			op        containerOperator
			op2       containerOperator
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

		bitList := bits.GetBitList(blit)
		for _, dim := range dims {
			for _, bit := range bitList {
				// find target cells with the same hint as the source cell,
				// populate the "items *[]int" slice with candidate cells
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

				// get cells with a matching pair in the target dimension (row or column)
				var pairs []int
				if err := dim.op(i, findPairs(&pairs)); err != nil {
					return err
				}

				// ensure only one pair exists
				if len(pairs) != 1 {
					continue
				}
				// assigned the "locked pair position"
				lockedPairPos := pairs[0]
				//c2 := getCoords(lockedPairPos)

				// find all pairs with the original cell in the inverse dimension
				var pairs21 []int
				if err := dim.op2(i, findPairs(&pairs21)); err != nil {
					return err
				}

				// find all pairs with the "locked pair" cell in the inverse dimension
				var pairs22 []int
				if err := dim.op2(lockedPairPos, findPairs(&pairs22)); err != nil {
					return err
				}

				// TODO: item21/item22 must be the only cell with hint in their shared row/column
				//       NOTE: it looks like dim.isAligned is taking care of this.
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
						//c4 := getCoords(item22)

						//logged := false

						sourceList := []int{i, lockedPairPos, item21, item22}

						removeHints := func(target int, source int) error {
							for _, pos := range sourceList {
								if target == pos {
									return nil
								}
							}

							/*if !logged && b.willUpdateCandidates(target, source, ^bit) {
								logged = true
								b.PrintHints()
								fmt.Printf("xwing: val:%d\n", bits.GetSingleBitValue(bit))
								fmt.Printf("-- %#2v\n", c1)
								fmt.Printf("-- %#2v\n", c2)
								fmt.Printf("-- %#2v\n", c3)
								fmt.Printf("-- %#2v\n", c4)
							}*/

							return b.updateCandidates(target, source, ^bit)
						}

						for _, pos := range []int{i, lockedPairPos} {
							if err := dim.op2(pos, removeHints); err != nil {
								return err
							}
						}

						if b.changed {
							// let simpler techniques take over
							//b.Log(false, -1, "xwing: change detected")
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}
