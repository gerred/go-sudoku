package main

import (
	"fmt"

	"github.com/judwhite/go-sudoku/internal/bits"
)

func (b *board) SolveSwordFish() error {
	// http://www.sudokuwiki.org/Sword_Fish_Strategy
	// find a 3x3 which share the same candidate
	// hint cannot be repeated on container (row or col depending on orientation)
	// 333,332,322,222 = all valid.
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}

		//coords := getCoords(i)

		// create dimensions for looking for SwordFish in the row and column dimension
		dims := []struct {
			op          func(int, func(int, int) error) error
			op_inverted func(int, func(int, int) error) error
		}{
			{
				op:          b.operateOnRow,
				op_inverted: b.operateOnColumn,
			},
			{
				op:          b.operateOnColumn,
				op_inverted: b.operateOnRow,
			},
		}

		for dim_it, dim := range dims {
			bitList := bits.GetBitList(b.blits[i])

			var pickList []int

			for _, bit := range bitList {
				getPickList := func(target int, source int) error {
					if b.solved[target] != 0 {
						return nil
					}
					if b.blits[target]&bit == bit {
						pickList = append(pickList, target)
					}
					return nil
				}

				if err := dim.op(i, getPickList); err != nil {
					return err
				}

				if len(pickList) != 2 && len(pickList) != 3 {
					continue
				}

				// pickList contains all cells in the row/col container with
				// the current `bit` in its hints. this number has to be 2 or 3.
				fmt.Printf("swordfish: dim:%d num:%d\n", dim_it, bits.GetSingleBitValue(bit))
				for _, pick := range pickList {
					fmt.Printf("- %#2v %s\n", getCoords(pick), bits.GetString(b.blits[pick]))
				}

				// get other rows/cols
			}
		}
	}
	return nil
}

func (b *board) getSwordFishContainer(bit uint, rowcol int, inspect inspector) error {
	// TODO
	return nil
}
