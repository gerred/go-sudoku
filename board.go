package main

import (
	"bytes"
	"io/ioutil"
)

type board struct {
	solved  [81]uint
	blits   [81]uint
	loading bool
}

type coords struct {
	row int
	col int
	box int
}

type inspector func(int, int) error
type containerOperator func(int, inspector) error

func getCoords(pos int) coords {
	boxRow := ((pos / 9) / 3)
	boxCol := ((pos % 9) / 3)
	box := boxRow*3 + boxCol

	return coords{row: pos / 9, col: pos % 9, box: box}
}

func (b *board) operateOnRow(pos int, op inspector) error {
	startRow := (pos / 9) * 9
	for r := startRow; r < startRow+9; r++ {
		if err := op(r, pos); err != nil {
			return err
		}
	}
	return nil
}

func (b *board) operateOnColumn(pos int, op inspector) error {
	for c := pos % 9; c < 81; c += 9 {
		if err := op(c, pos); err != nil {
			return err
		}
	}
	return nil
}

func (b *board) operateOnBox(pos int, op inspector) error {
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

func (b *board) operateOnRCB(pos int, op inspector) error {
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

func (b *board) operateOnCommon(pos1 int, pos2 int, op inspector) error {
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

func getBoard(fileName string) (*board, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return loadBoard(b)
}

func loadBoard(b []byte) (*board, error) {
	b = bytes.Replace(b, []byte{'\r'}, []byte{}, -1)
	b = bytes.Replace(b, []byte{'\n'}, []byte{}, -1)
	b = bytes.Replace(b, []byte{' '}, []byte{}, -1)

	board := &board{loading: true}
	for i := 0; i < 81; i++ {
		board.blits[i] = 0x1FF
	}

	for i := 0; i < 81; i++ {
		if b[i] != '_' && b[i] != '0' {
			val := uint(b[i] - 48)
			if err := board.SolvePosition(i, val); err != nil {
				return board, err
			}
		}
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

func intersect(a []int, b []int) []int {
	var list []int
	for _, i := range a {
		for _, j := range b {
			if i == j {
				list = append(list, i)
			}
		}
	}
	return list
}
