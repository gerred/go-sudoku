package main

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

		ops := []containerOperator{
			b.operateOnRow,
			b.operateOnColumn,
			b.operateOnBox,
		}

		for _, op := range ops {
			sumBlits = 0
			if err := op(i, sumHints); err != nil {
				return err
			}
			leftOver := blit & ^sumBlits

			if HasSingleBit(leftOver) {
				val := GetSingleBitValue(leftOver)
				if err := b.SolvePosition(i, val); err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}
