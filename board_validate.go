package main

import "fmt"

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