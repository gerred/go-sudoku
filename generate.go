package main

import (
	"fmt"
	"math/rand"
)

func getValidBoard() (*board, error) {
	b, err := loadBoard([]byte("000000000000000000000000000000000000000000000000000000000000000000000000000000000"))
	if err != nil {
		return nil, err
	}

	for !b.isSolved() {
		n := rand.Intn(81)
		if b.solved[n] != 0 {
			continue
		}
		bitList := GetBitList(b.blits[n])
		bn := rand.Intn(len(bitList))
		val := GetSingleBitValue(bitList[bn])

		err = b.SolvePosition(n, val)
		if err != nil {
			return nil, err
		}

		err = b.Solve()
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func generate() error {
	var err error
	var b *board
	for b == nil || err != nil {
		b, err = getValidBoard()
	}

	//b.PrintHints()

	err = digHoles(b)
	if err != nil {
		return err
	}

	/*b.CountSolutions = true

	err = b.SolveSAT()
	b.Print()
	if err != nil {
		log.Fatal(err)
	}*/

	return nil
}

func digHoles(b *board) error {
	var err error
	b2 := &board{solved: b.solved, blits: b.blits}
	if err != nil {
		return err
	}

	step := 1
	failures := 0
	check := make(map[int]interface{})
	for len(check) != 81 {
		goodSolved := b2.solved
		goodBlits := b2.blits

		pos1 := rand.Intn(81)
		if step == 1 {
			if _, ok := check[pos1]; ok {
				continue
			}
			check[pos1] = struct{}{}
		}
		if b2.solved[pos1] == 0 {
			continue
		}

		coords := getCoords(pos1)
		secondRow := 8 - coords.row
		if secondRow < 0 {
			secondRow += 8
		}
		secondCol := 8 - coords.col
		if secondCol < 0 {
			secondCol += 8
		}

		if step == 4 {
			pos2 := coords.row*9 + secondCol
			pos3 := secondRow*9 + coords.col
			pos4 := secondRow*9 + secondCol

			if b2.solved[pos2] == 0 || b2.solved[pos3] == 0 || b2.solved[pos4] == 0 {
				continue
			}

			b2.solved[pos1] = 0
			b2.solved[pos2] = 0
			b2.solved[pos3] = 0
			b2.solved[pos4] = 0
		} else if step == 2 {
			pos2 := secondRow*9 + secondCol

			if b2.solved[pos2] == 0 {
				continue
			}

			b2.solved[pos1] = 0
			b2.solved[pos2] = 0
		} else {
			b2.solved[pos1] = 0
		}

		for j := 0; j < 81; j++ {
			if b2.solved[j] != 0 {
				continue
			}
			newHints, err := b2.getHints(j)
			if err != nil {
				return err
			}
			b2.blits[j] = newHints
		}

		b3 := board{solved: b2.solved, blits: b2.blits, CountSolutions: true, MaxSolutions: 2}
		err = b3.Solve()
		if err != nil {
			return err
		}

		//fmt.Printf("sln count: %d\n", b3.SolutionCount)

		if b3.SolutionCount > 1 {
			//return fmt.Errorf("bad dig, more than 1 solution")
			b2.solved = goodSolved
			b2.blits = goodBlits
			failures++
			if step > 1 && failures == 5 {
				failures = 0
				step /= 2
			}
		} else {
			//b2.PrintHints()
		}
	}

	fmt.Printf("-----------------\n")
	b2.Print()
	b2.PrintURL()
	fmt.Printf("hint count: %d\n", b2.numSolved())
	b2.CountSolutions = true
	b2.MaxSolutions = 2
	if err = b2.Solve(); err != nil {
		return err
	}
	fmt.Printf("sln count: %d\n", b2.SolutionCount)
	//b2.PrintHints()
	b2.Print()

	return nil
}

func (b *board) getHints(pos int) (uint, error) {
	check := make(map[uint]interface{})
	for i := uint(1); i <= 9; i++ {
		check[i] = struct{}{}
	}

	removeHints := func(target int, source int) error {
		if target == source {
			return nil
		}
		val := b.solved[target]
		if val == 0 {
			return nil
		}

		if _, ok := check[val]; ok {
			delete(check, val)
		}

		return nil
	}

	if err := b.operateOnRCB(pos, removeHints); err != nil {
		// TODO: return err
		return 0, err
	}

	blits := uint(0)
	for k := range check {
		blits |= 1 << (k - 1)
	}
	return blits, nil
}
