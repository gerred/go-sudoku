package main

import (
	"fmt"
	"strings"

	"github.com/judwhite/go-sudoku/internal/bits"
)

type xcycle struct {
	coords coords
	// a strong link says that if this cell is OFF then the next cell
	// it points to must be ON (only two candidates in visible intersection)
	isPrevStrongLink bool
	prev             *xcycle
	next             []*xcycle
	visible          []int
	depth            int
	canSeeStart      bool
	closedLoop       bool
}

func (b *board) SolveXCycles() error {
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}

		hintList := bits.GetBitList(b.blits[i])
		for _, hint := range hintList {
			// TODO: debug, remove
			/*if hint == 0x80 {
				continue
			}*/

			visible := b.getVisibleCellsWithHint(i, hint)
			if len(visible) == 0 {
				continue
			}

			cycle := xcycle{
				coords:  getCoords(i),
				visible: visible,
			}

			b.xcyclesRecurse(&cycle, hint, i)

			fmt.Printf("+ hint: %d\n", bits.GetSingleBitValue(hint))
			modified, err := b.xcyclesPrint(&cycle, hint)
			if err != nil {
				return err
			}
			if modified {
				return nil
			}
		}
	}
	return nil
}

func (b *board) xcyclesPrint(cycle *xcycle, hint uint) (bool, error) {
	// check if we have an even number of cells > 2 and the cell can see the start
	if cycle.depth%2 == 0 && cycle.depth >= 3 && cycle.closedLoop {
		// check that there aren't two weak links in a row
		// this first pass is a simplification
		x := cycle
		valid := true
		for x.prev != nil {
			if !x.isPrevStrongLink && !x.prev.isPrevStrongLink {
				valid = false
				break
			}
			x = x.prev
		}

		if valid {
			// first pass is again a simplification
			//
			// * check for any cells which can be seen by an odd and even depth cell in the cycle,
			// * which is itself not in the cycle
			//
			// since we've validated that our cells don't contain any consecutive weak links, the
			// entire cycle can be thought of as half on, or half off, alternating.
			// maybe... it would definitely work if all links were strong links.

			var evenVisible []int
			var oddVisible []int
			var inChain []int
			x = cycle
			for x != nil {
				inChain = append(inChain, x.coords.pos)
				visible := b.getVisibleCellsWithHint(x.coords.pos, hint)
				if x.depth%2 == 0 {
					evenVisible = union(evenVisible, visible)
				} else {
					oddVisible = union(oddVisible, visible)
				}
				//prefix := strings.Repeat("+", x.depth+1)
				//fmt.Printf("%s %#2v strong:%t %s can-see-start:%t len(cycle.next)=%d\n", prefix, x.coords, x.isPrevStrongLink, bits.GetString(b.blits[x.coords.pos]), x.canSeeStart, len(x.next))
				x = x.prev
			}

			common := intersect(evenVisible, oddVisible)
			common = subtract(common, inChain)

			if len(common) != 0 {
				x = cycle
				for x != nil {
					prefix := strings.Repeat("+", x.depth+1)
					fmt.Printf("%s %#2v strong:%t %s can-see-start:%t len(cycle.next)=%d\n", prefix, x.coords, x.isPrevStrongLink, bits.GetString(b.blits[x.coords.pos]), x.canSeeStart, len(x.next))
					x = x.prev
				}
				for _, item := range common {
					if err := b.updateCandidates(item, cycle.coords.pos, ^hint); err != nil {
						return false, err
					}
					fmt.Printf("- %#2v %s\n", getCoords(item), bits.GetString(b.blits[item]))
				}
				return true, nil
			}
		}
	}

	for _, item := range cycle.next {
		modified, err := b.xcyclesPrint(item, hint)
		if err != nil {
			return false, err
		}
		if modified {
			return true, nil
		}
	}

	return false, nil
}

func (b *board) xcyclesRecurse(cycle *xcycle, hint uint, startPos int) {
	if cycle.coords.pos == startPos && cycle.prev != nil {
		return
	}
	//fmt.Printf("- %#2v depth:%d hint:%d %v\n", cycle.coords, cycle.depth, bits.GetSingleBitValue(hint), cycle)
	//fmt.Printf("- hint: %d\n", bits.GetSingleBitValue(hint))
	//fmt.Printf("-- visible %d: %#2v %s -- common: %v\n", idx, getCoords(cell), bits.GetString(b.blits[cell]), commonCoords)
	//fmt.Printf("--- %#v\n", cycle)
	for _, cell := range cycle.visible {
		visible := b.getVisibleCellsWithHint(cell, hint)
		if len(visible) == 0 {
			continue
		}

		var filtered []int
		var canSeeStart bool
		excludeFirstItem := cycle.depth%2 == 0
		for _, item := range visible {
			if !cycle.isInListBacktrack(item, excludeFirstItem) {
				filtered = append(filtered, item)
			}
			if item == startPos {
				canSeeStart = true
			}
		}

		commonList := intersect(cycle.visible, visible)

		nextCycle := &xcycle{
			coords:           getCoords(cell),
			prev:             cycle,
			visible:          filtered,
			isPrevStrongLink: len(commonList) == 0,
			depth:            cycle.depth + 1,
			canSeeStart:      canSeeStart,
			closedLoop:       cell == startPos,
		}

		cycle.next = append(cycle.next, nextCycle)

		// let's be reasonable...
		if nextCycle.depth < 6 {
			b.xcyclesRecurse(nextCycle, hint, startPos)
		}
	}
}

func (x *xcycle) isInListBacktrack(pos int, excludeFirstItem bool) bool {
	// current
	if x.coords.pos == pos {
		return true
	}

	if x.prev == nil {
		return false
	}

	// rewind
	a := x.prev
	for a.prev != nil {
		if a.coords.pos == pos {
			if excludeFirstItem && a.prev == nil {
				return false
			}
			return true
		}
		a = a.prev
	}

	return false
}
