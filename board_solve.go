package main

import (
	"fmt"
	"reflect"

	"github.com/judwhite/go-sudoku/internal/bits"
)

func (b *board) Solve() error {
	// first iteration naked single
	b.loading = true // turn off logging, this run is boring
	fmt.Println("--- NAKED SINGLE: FIRST ITERATION")
	if err := b.SolveNakedSingle(); err != nil {
		return err
	}
	b.PrintHints()
	b.loading = false

	if err := b.runSolvers(); err != nil {
		return err
	}

	return nil
}

func (b *board) runSolver(solver func() error) (bool, error) {
	oldBlits := b.blits
	if err := solver(); err != nil {
		return false, err
	}
	changed := !reflect.DeepEqual(oldBlits, b.blits)
	return changed, nil
}

func (b *board) runSolverN(solver func(int) error, n int) (bool, error) {
	oldBlits := b.blits
	if err := solver(n); err != nil {
		return false, err
	}
	changed := !reflect.DeepEqual(oldBlits, b.blits)
	return changed, nil
}

func (b *board) runSolvers() error {
	var err error
	var changed bool

	for !b.isSolved() {
		fmt.Println("--- NAKED SINGLE")
		if changed, err = b.runSolver(b.SolveNakedSingle); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- HIDDEN SINGLE")
		if changed, err = b.runSolver(b.SolveHiddenSingle); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- NAKED PAIR")
		if changed, err = b.runSolverN(b.SolveNakedN, 2); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- NAKED TRIPLE")
		if changed, err = b.runSolverN(b.SolveNakedN, 3); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- NAKED QUAD")
		if changed, err = b.runSolverN(b.SolveNakedN, 4); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- NAKED QUINT") // not seen in any tests yet
		if changed, err = b.runSolverN(b.SolveNakedN, 5); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- HIDDEN PAIR")
		if changed, err = b.runSolverN(b.SolveHiddenN, 2); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- HIDDEN TRIPLE")
		if changed, err = b.runSolverN(b.SolveHiddenN, 3); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- HIDDEN QUAD") // not seen in any tests yet
		if changed, err = b.runSolverN(b.SolveHiddenN, 4); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- HIDDEN QUINT")
		if changed, err = b.runSolverN(b.SolveHiddenN, 5); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- POINTING PAIR AND TRIPLE REDUCTION")
		if changed, err = b.runSolver(b.SolvePointingPairAndTripleReduction); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- BOX LINE")
		if changed, err = b.runSolver(b.SolveBoxLine); err != nil {
			return err
		}
		if changed {
			continue
		}

		fmt.Println("--- XWING")
		if changed, err = b.runSolver(b.SolveXWing); err != nil {
			return err
		}
		if changed {
			continue
		}

		// TODO: remove, temp to cleanup SWORDFISH candidates
		fmt.Println("*****************************************")
		b.PrintHints()
		fmt.Println("*****************************************")
		// END TODO

		fmt.Println("--- SWORDFISH")
		if changed, err = b.runSolver(b.SolveSwordFish); err != nil {
			return err
		}
		if changed {
			continue
		}

		break
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
