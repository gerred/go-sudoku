package main

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
			op    containerOperator
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
				perms := getPermutations(x, pickList, []int{i})
				sumBits := uint(0)
				for _, list := range perms {
					for _, item := range list {
						sumBits |= b.blits[item]
					}

					leftOver := sumBits & ^sumNegateBits
					nbits := GetNumberOfSetBits(leftOver)
					//if nbits != 2 && nbits != 3 {
					if nbits == 0 || nbits > 3 {
						continue
					}

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
				}
			}
		}
	}
	return nil
}
