package main

import "github.com/judwhite/go-sudoku/internal/bits"

func (b *board) SolveXYChain() error {
	// http://www.sudokuwiki.org/XY_Chains
	// bi-value cells linked together by one value (and visible to each other)
	// terminate when more than one away and cell shares a value with the other end
	// cells visible by both ends of the chain can have their shared value removed.
	for i := 0; i < 81; i++ {
		blit := b.blits[i]
		if bits.GetNumberOfSetBits(blit) != 2 {
			continue
		}

		bits := bits.GetBitList(blit)
		for _, bit := range bits {
			updated, err := b.xyChainTestPosition(i, bit)
			if err != nil {
				return err
			}
			if updated {
				// let simpler techniques take over
				return nil
			}
		}
	}
	return nil
}

// itemBlit&^excludeBit&firstBitInChain
// 5,9 - 5
// 5,1 - 1
// 1,9 - 9 // odd, has to fit original hint
// 19&^1=9

// 5,9 - 5
// 5,1 - 1
// 1,9 - 9
// 1,9 - 1 // even, has to fit original hint

// leftover has to fit original hint

func (b *board) xyChainTestPosition(i int, excludeBit uint) (bool, error) {
	hint := excludeBit
	lists := b.xyChainFollow([]int{i}, excludeBit, hint, 1)
	for _, list := range lists {
		startPos := list[0]
		endPos := list[len(list)-1]
		visible1 := b.getVisibleCells(startPos)
		visible2 := b.getVisibleCells(endPos)
		targets := intersect(visible1, visible2)
		if len(targets) == 0 {
			continue
		}

		/*var once1 sync.Once
		print1 := func() {
			fmt.Printf("-/- %#v hint:%d\n", list, bits.GetSingleBitValue(hint))
			for idx, chainItem := range list {
				fmt.Printf("--- chaind %d: %#2v %s\n", idx, getCoords(chainItem), bits.GetString(b.blits[chainItem]))
			}
			fmt.Printf("----- targets:\n")
			for _, target := range targets {
				fmt.Printf("----- %#2v %s\n", getCoords(target), bits.GetString(b.blits[target]))
			}
		}*/

		updated := false
	targetLoop:
		for _, target := range targets {
			// items in the chain aren't candidates (but why not? shouldn't the logic hold? TODO)
			//once1.Do(print1)
			for _, chainItem := range list {
				if target == chainItem {
					continue targetLoop
				}
			}

			targetBlit := b.blits[target]
			if targetBlit&hint == hint {
				updated = true
				if err := b.updateCandidates(target, i, ^hint); err != nil {
					return false, err
				}
			}
		}

		if updated {
			// let simpler techniques take over
			return true, nil
		}
	}
	return false, nil
}

func (b *board) xyChainFollow(chain []int, excludeBit uint, firstBitInChain uint, depth int) [][]int {
	var lists [][]int

	//firstBlit := b.blits[chain[0]]
	curPos := chain[len(chain)-1]
	curBlit := b.blits[curPos]

	visible := b.getVisibleCells(curPos)

	var filtered []int
loopVisible:
	for _, item := range visible {
		// avoid cycles
		for _, prevItem := range chain {
			if prevItem == item {
				continue loopVisible
			}
		}
		// ensure cell has 2 hints and is linked to the previous cell
		itemBlit := b.blits[item]
		if bits.GetNumberOfSetBits(itemBlit) != 2 || bits.GetNumberOfSetBits(curBlit&itemBlit) == 0 {
			continue
		}
		if curBlit&itemBlit == excludeBit {
			continue
		}
		filtered = append(filtered, item)
	}

	if len(filtered) == 0 {
		return lists
	}

	/*if len(chain) == 1 {
		fmt.Printf("* %#2v %s\n", getCoords(curPos), bits.GetString(curBlit))
	}*/

	//prefix := strings.Repeat("-", depth+1)
	for _, item := range filtered {
		//fmt.Printf("%s %#2v %s\n", prefix, getCoords(item), bits.GetString(b.blits[item]))
		itemBlit := b.blits[item]

		var newChain []int
		newChain = append(newChain, chain...)
		newChain = append(newChain, item)

		nextExcludeBit := curBlit & itemBlit &^ excludeBit

		if len(chain) > 1 {
			if itemBlit&^nextExcludeBit&firstBitInChain == firstBitInChain {
				/*itemCoords := getCoords(item)
				fmt.Printf("**** r,c:               %d,%d\n", itemCoords.row, itemCoords.col)
				fmt.Printf("**** first bit:         %09b %d\n", firstBitInChain, bits.GetSingleBitValue(firstBitInChain))
				fmt.Printf("**** item blit:         %09b\n", itemBlit)
				fmt.Printf("**** next exclude bit:  %09b %d\n", nextExcludeBit, bits.GetSingleBitValue(nextExcludeBit))
				fmt.Printf("**** ^exclude bit:      %09b\n", ^nextExcludeBit&0x1FF)
				fmt.Printf("**** itemBlit&^exclude: %09b\n", itemBlit&^nextExcludeBit)
				fmt.Printf("**** test:              %09b\n", itemBlit&^nextExcludeBit&firstBitInChain)*/

				// TODO: should we keep going? there may be longer chains
				lists = append(lists, newChain)
				continue
			}
		}

		newLists := b.xyChainFollow(newChain, nextExcludeBit, firstBitInChain, depth+1)
		lists = append(lists, newLists...)
	}

	return lists
}
