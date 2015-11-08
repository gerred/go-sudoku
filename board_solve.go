package main

import (
	"fmt"
	"reflect"

	"github.com/judwhite/go-sudoku/internal/bits"
)

func (b *board) Solve() error {
	doLoop := true

	// first iteration naked single
	b.loading = true // turn off logging, this run is boring
	fmt.Println("--- NAKED SINGLE: FIRST ITERATION")
	if err := b.SolveNakedSingle(); err != nil {
		return err
	}
	b.PrintHints()
	b.loading = false

	for i := 0; doLoop; i++ {
		oldBlits := b.blits
		if err := b.runSolvers(); err != nil {
			return err
		}
		doLoop = !reflect.DeepEqual(oldBlits, b.blits)
		if b.isSolved() {
			break
		}
		fmt.Printf("doLoop: %t\n", doLoop)

		if !doLoop {
			// now we branch
			// TODO
		}
	}
	return nil
}

func (b *board) runSolvers() error {
	var err error
	fmt.Println("--- NAKED SINGLE")
	if err = b.SolveNakedSingle(); err != nil {
		return err
	}
	fmt.Println("--- HIDDEN SINGLE")
	if err = b.SolveHiddenSingle(); err != nil {
		return err
	}
	fmt.Println("--- NAKED PAIR")
	if err = b.SolveNakedN(2); err != nil {
		return err
	}
	fmt.Println("--- NAKED TRIPLE")
	if err = b.SolveNakedN(3); err != nil {
		return err
	}
	fmt.Println("--- NAKED QUAD")
	if err = b.SolveNakedN(4); err != nil {
		return err
	}
	/*fmt.Println("--- NAKED QUINT")
	if err = b.SolveNakedN(5); err != nil {
		return err
	}*/
	fmt.Println("--- HIDDEN PAIR")
	if err = b.SolveHiddenN(2); err != nil {
		return err
	}
	fmt.Println("--- HIDDEN TRIPLE")
	if err = b.SolveHiddenN(3); err != nil {
		return err
	}
	/*fmt.Println("--- HIDDEN QUAD")
	if err = b.SolveHiddenN(4); err != nil {
		return err
	}
	fmt.Println("--- HIDDEN QUINT")
	if err = b.SolveHiddenN(5); err != nil {
		return err
	}*/
	fmt.Println("--- POINTING PAIR AND TRIPLE REDUCTION")
	if err = b.SolvePointingPairAndTripleReduction(); err != nil {
		return err
	}
	fmt.Println("--- BOX LINE")
	if err = b.SolveBoxLine(); err != nil {
		return err
	}

	fmt.Println("--- XWING")
	if err = b.SolveXWing(); err != nil {
		return err
	}

	// TODO: remove, temp to cleanup SWORDFISH candidates
	fmt.Println("--- NAKED SINGLE")
	if err = b.SolveNakedSingle(); err != nil {
		return err
	}
	fmt.Println("--- HIDDEN SINGLE")
	if err = b.SolveHiddenSingle(); err != nil {
		return err
	}

	fmt.Println("*****************************************")
	b.PrintHints()
	fmt.Println("*****************************************")
	// END TODO

	fmt.Println("--- SWORDFISH")
	if err = b.SolveSwordFish(); err != nil {
		return err
	}

	fmt.Println("--- END SOLVER")

	return nil
}

func (b *board) SolvePosition(pos int, val uint) error {
	mask := uint(^(1 << (val - 1)))
	if b.solved[pos] != 0 /*&& b.solved[pos] != val*/ {
		return fmt.Errorf("pos %d has value %d, tried to set with %d", pos, b.solved[pos], val)
	}
	b.solved[pos] = val
	b.blits[pos] = 1 << (val - 1)

	b.Log(true, pos, fmt.Sprintf("set value %d mask:%09b", val, mask&0x1FF))

	if !b.loading {
		b.Print()
		b.PrintHints()
	}

	if err := b.Validate(); err != nil {
		return fmt.Errorf("%#v val:%d - %s", getCoords(pos), val, err)
	}

	if err := b.operateOnRCB(pos, b.removeCandidates(mask)); err != nil {
		return err
	}

	if !b.loading {
		b.PrintHints()
	}

	return nil
}

func (b *board) removeCandidates(mask uint) func(int, int) error {
	return func(target int, source int) error {
		if opErr := b.updateCandidates(target, source, mask); opErr != nil {
			return opErr
		}
		return nil
	}
}

func (b *board) updateCandidates(targetPos int, sourcePos int, mask uint) error {
	if targetPos == sourcePos || b.solved[targetPos] != 0 {
		return nil
	}
	oldBlit := b.blits[targetPos]
	newBlit := oldBlit & mask
	if newBlit != oldBlit {
		if newBlit == 0 {
			return fmt.Errorf("tried to remove last candidate from %#2v", getCoords(targetPos))
		}

		b.blits[targetPos] = newBlit
		delta := oldBlit & ^newBlit
		b.Log(false, targetPos, fmt.Sprintf("old hints: %-10s remove hint: %s remaining hints: %s", bits.GetString(oldBlit), bits.GetString(delta), bits.GetString(newBlit)))
	}
	return nil
}

func getPermutations(n int, pickList []int, curList []int) [][]int {
	output := make([][]int, 0)

	for i := 0; i < len(pickList); i++ {
		list := make([]int, len(curList))
		copy(list, curList)              // get the source list
		list = append(list, pickList[i]) // plus the current element

		if len(list) == n {
			// if this is the length we're looking for...
			output = append(output, list)
		} else {
			// otherwise, call recursively
			perms := getPermutations(n, pickList[i+1:], list)
			if perms != nil {
				for _, v := range perms {
					output = append(output, v)
				}
			}
		}
	}

	return output
}
