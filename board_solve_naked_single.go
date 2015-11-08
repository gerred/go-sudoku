package main

import "github.com/judwhite/go-sudoku/internal/bits"

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
