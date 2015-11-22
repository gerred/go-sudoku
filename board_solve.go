package main

import "fmt"

func (b *board) SolveWithSolversList(solvers []solver) error {
	// first iteration naked single
	b.loading = true // turn off logging, this run is boring

	if b.verbose {
		b.PrintHints()
		b.Log(false, -1, "NAKED SINGLE: FIRST ITERATION")
	}

	if err := b.SolveNakedSingle(); err != nil {
		return err
	}

	if b.verbose {
		b.PrintHints()
	}

	b.loading = false

	if err := b.runSolvers(solvers); err != nil {
		return err
	}

	return nil
}

func (b *board) Solve() error {
	return b.SolveWithSolversList(b.getSolvers())
}

type solver struct {
	run        func() error
	name       string
	printBoard bool
}

func (b *board) getSimpleSolvers() []solver {
	solvers := []solver{
		{name: "NAKED SINGLE", run: b.SolveNakedSingle},
		{name: "HIDDEN SINGLE", run: b.SolveHiddenSingle},
		{name: "NAKED PAIR", run: b.getSolverN(b.SolveNakedN, 2)},
		{name: "NAKED TRIPLE", run: b.getSolverN(b.SolveNakedN, 3)},
		{name: "NAKED QUAD", run: b.getSolverN(b.SolveNakedN, 4)},
		{name: "NAKED QUINT", run: b.getSolverN(b.SolveNakedN, 5)},
		{name: "HIDDEN PAIR", run: b.getSolverN(b.SolveHiddenN, 2)},
		{name: "HIDDEN TRIPLE", run: b.getSolverN(b.SolveHiddenN, 3)},
		{name: "HIDDEN QUAD", run: b.getSolverN(b.SolveHiddenN, 4)},
		{name: "HIDDEN QUINT", run: b.getSolverN(b.SolveHiddenN, 5)},
		{name: "POINTING PAIR AND TRIPLE REDUCTION", run: b.SolvePointingPairAndTripleReduction},
		{name: "BOX LINE", run: b.SolveBoxLine},
	}

	return solvers
}

func (b *board) getSolvers() []solver {
	solvers := []solver{
		{name: "NAKED SINGLE", run: b.SolveNakedSingle},
		{name: "HIDDEN SINGLE", run: b.SolveHiddenSingle},
		{name: "NAKED PAIR", run: b.getSolverN(b.SolveNakedN, 2)},
		{name: "NAKED TRIPLE", run: b.getSolverN(b.SolveNakedN, 3)},
		{name: "NAKED QUAD", run: b.getSolverN(b.SolveNakedN, 4)},
		{name: "NAKED QUINT", run: b.getSolverN(b.SolveNakedN, 5)},
		{name: "HIDDEN PAIR", run: b.getSolverN(b.SolveHiddenN, 2)},
		{name: "HIDDEN TRIPLE", run: b.getSolverN(b.SolveHiddenN, 3)},
		{name: "HIDDEN QUAD", run: b.getSolverN(b.SolveHiddenN, 4)},
		{name: "HIDDEN QUINT", run: b.getSolverN(b.SolveHiddenN, 5)},
		{name: "POINTING PAIR AND TRIPLE REDUCTION", run: b.SolvePointingPairAndTripleReduction},
		{name: "BOX LINE", run: b.SolveBoxLine},
		{name: "X-WING", run: b.SolveXWing},
		{name: "SIMPLE-COLORING", run: b.SolveSimpleColoring},
		{name: "Y-WING", run: b.SolveYWing},
		{name: "SWORDFISH", run: b.SolveSwordFish},
		{name: "XY-CHAIN", run: b.SolveXYChain},
		{name: "EMPTY RECTANGLES", run: b.SolveEmptyRectangles},
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

func (b *board) runSolvers(solvers []solver) error {
mainLoop:
	for !b.isSolved() {
		b.changed = false
		for _, solver := range solvers {
			if solver.printBoard {
				b.PrintHints()
				b.PrintURL()
			}

			if b.verbose {
				b.Log(false, -1, solver.name)
			}

			if err := solver.run(); err != nil {
				return err
			}
			if b.changed {
				/*if solver.safetyCheck && testBoard != nil {
					for i := 0; i < 81; i++ {
						if b.blits[i]&testBoard.blits[i] != testBoard.blits[i] {
							fmt.Printf("-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/\n")
							fmt.Printf("%#2v\n", getCoords(i))
							fmt.Printf("%09b\n", b.blits[i])
							fmt.Printf("%09b\n", testBoard.blits[i])
							b.PrintHints()
							testBoard.Print()
							fmt.Printf("-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/\n")
							return fmt.Errorf("error at %#v", getCoords(i))
						}
					}
					testBoard, _ = TrialAndError(*b)
				}*/

				if solver.printBoard || b.verbose {
					b.PrintHints()
					b.PrintURL()
				}

				continue mainLoop
			}
			if b.isSolved() {
				break
			}
		}

		if !b.isSolved() && !b.SkipSAT {
			//b.Log(false, -1, "Starting SAT solver...")
			//start := time.Now()
			err := b.SolveSAT()
			//fmt.Printf("sat time: %v\n", time.Since(start))
			if err != nil {
				return err
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

	//fmt.Println("--- END SOLVER")

	return nil
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

func (b *board) SolvePositionNoValidate(pos int, val uint) {
	b.solved[pos] = val
	b.blits[pos] = 1 << (val - 1)
}

func (b *board) SolvePosition(pos int, val uint) error {
	mask := uint(^(1 << (val - 1)))
	if b.solved[pos] != 0 /*&& b.solved[pos] != val*/ {
		return fmt.Errorf("pos %d has value %d, tried to set with %d", pos, b.solved[pos], val)
	}
	b.solved[pos] = val
	b.blits[pos] = 1 << (val - 1)

	if b.verbose {
		b.Log(true, pos, fmt.Sprintf("set value %d mask:%09b", val, mask&0x1FF))
	}

	if err := b.Validate(); err != nil {
		return fmt.Errorf("%#v val:%d - %s", getCoords(pos), val, err)
	}

	if err := b.operateOnRCB(pos, b.removeCandidates(mask)); err != nil {
		return err
	}

	if !b.loading && b.verbose {
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
		b.changed = true
		if newBlit == 0 {
			return fmt.Errorf("tried to remove last candidate from %#2v", getCoords(targetPos))
		}

		b.blits[targetPos] = newBlit
		if b.verbose {
			delta := oldBlit & ^newBlit
			b.Log(false, targetPos, fmt.Sprintf("old hints: %-10s remove hint: %s remaining hints: %s", GetBitsString(oldBlit), GetBitsString(delta), GetBitsString(newBlit)))
		}
		return b.Validate()
	}
	return nil
}

func getPermutations(n int, pickList []int, curList []int) [][]int {
	var output [][]int

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
