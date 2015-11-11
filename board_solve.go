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

type solver struct {
	run        func() error
	name       string
	printBoard bool
}

func (b *board) getSolvers() []solver {
	solvers := []solver{
		{name: "NAKED SINGLE", run: b.SolveNakedSingle},
		{name: "HIDDEN SINGLE", run: b.SolveHiddenSingle},
		{name: "NAKED PAIR", run: b.getSolverN(b.SolveNakedN, 2)},
		{name: "NAKED TRIPLE", run: b.getSolverN(b.SolveNakedN, 3)},
		{name: "NAKED QUAD", run: b.getSolverN(b.SolveNakedN, 4)},
		{name: "NAKED QUINT", run: b.getSolverN(b.SolveNakedN, 5)}, // NOTE: not seen in any tests yet
		{name: "HIDDEN PAIR", run: b.getSolverN(b.SolveHiddenN, 2)},
		{name: "HIDDEN TRIPLE", run: b.getSolverN(b.SolveHiddenN, 3)},
		{name: "HIDDEN QUAD", run: b.getSolverN(b.SolveHiddenN, 4)}, // NOTE: not seen in any tests yet
		{name: "HIDDEN QUINT", run: b.getSolverN(b.SolveHiddenN, 5)},
		{name: "POINTING PAIR AND TRIPLE REDUCTION", run: b.SolvePointingPairAndTripleReduction},
		{name: "BOX LINE", run: b.SolveBoxLine},
		{name: "X-WING", run: b.SolveXWing},
		{name: "Y-WING", run: b.SolveYWing},
		//{name: "SWORDFISH", run: b.SolveSwordFish},
		{name: "X-CYCLES", run: b.SolveXCycles},
		//{name: "XY-CHAIN", run: b.SolveXYChain, printBoard: true},
	}

	return solvers
}

func (b *board) getSolverN(solver func(int) error, n int) func() error {
	return func() error {
		if err := solver(n); err != nil {
			return err
		}
		return nil
	}
}

func (b *board) runSolvers() error {
	solvers := b.getSolvers()

mainLoop:
	for !b.isSolved() {
		for _, solver := range solvers {
			oldBlits := b.blits
			fmt.Printf("--- %s\n", solver.name)
			if solver.printBoard {
				b.PrintHints()
				b.PrintURL()
			}

			if err := solver.run(); err != nil {
				return err
			}
			if !reflect.DeepEqual(oldBlits, b.blits) {
				if solver.printBoard {
					b.PrintHints()
					b.PrintURL()
				}
				continue mainLoop
			}
			if b.isSolved() {
				break
			}
		}

		/*newBoard, err := TrialAndError(*b)
		if err != nil {
			return err
		}
		if newBoard == nil {
			return fmt.Errorf("board has no solution")
		}

		b.solved = newBoard.solved
		b.blits = newBoard.blits*/
		break
	}

	fmt.Println("--- END SOLVER")

	return nil
}

func TrialAndError(b board) (*board, error) {
	// what type of cell would make a good candidate?
	// one where choosing a candidate would result in a lot of eliminations
	// we can use b.getVisibleCellsWithHint

	// cell pos, hint, number of eliminations
	store := make(map[int]map[uint]int)

	maxEliminations := 0
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}

		hintEliminations := make(map[uint]int)
		store[i] = hintEliminations

		bitList := bits.GetBitList(b.blits[i])
		for _, hint := range bitList {
			n := len(b.getVisibleCellsWithHint(i, hint))
			hintEliminations[hint] = n
			if n > maxEliminations {
				maxEliminations = n
			}
		}
	}

	for eliminations := maxEliminations; eliminations > 0; eliminations-- {
		fmt.Printf("eliminations: %d\n", eliminations)
		posHints := getCandidates(eliminations, store)
		if len(posHints) == 0 {
			continue
		}
		for pos, hints := range posHints {
			for _, hint := range hints {
				fmt.Printf("%#2v hint:%d\n", getCoords(pos), bits.GetSingleBitValue(hint))
				// try the elimination
				trialBoard := &board{solved: b.solved, blits: b.blits, loading: false, quiet: true}
				trialBoard.SolvePosition(pos, bits.GetSingleBitValue(hint))
				if err := trialBoard.Solve(); err != nil {
					fmt.Printf("... that board didn't work\n")
					// nope.
				} else {
					trialBoard.quiet = false
					return trialBoard, nil
				}
			}
		}
	}

	return nil, nil
}

func getCandidates(n int, store map[int]map[uint]int) map[int][]uint {
	list := make(map[int][]uint)
	for pos, hintEliminations := range store {
		for hint, eliminations := range hintEliminations {
			if eliminations == n {
				hintList := list[pos]
				list[pos] = append(hintList, hint)
			}
		}
	}
	return list
}

func (b *board) SolvePosition(pos int, val uint) error {
	mask := uint(^(1 << (val - 1)))
	if b.solved[pos] != 0 /*&& b.solved[pos] != val*/ {
		return fmt.Errorf("pos %d has value %d, tried to set with %d", pos, b.solved[pos], val)
	}
	b.solved[pos] = val
	b.blits[pos] = 1 << (val - 1)

	b.Log(true, pos, fmt.Sprintf("set value %d mask:%09b", val, mask&0x1FF))

	/*if !b.loading {
		b.Print()
		b.PrintHints()
	}*/

	if err := b.Validate(); err != nil {
		return fmt.Errorf("%#v val:%d - %s", getCoords(pos), val, err)
	}

	if err := b.operateOnRCB(pos, b.removeCandidates(mask)); err != nil {
		return err
	}

	/*if !b.loading {
		b.PrintHints()
	}*/

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
		return b.Validate()
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
