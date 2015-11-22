package main

import "sort"

func (b *board) SolveSimpleColoring() error {
	//b.PrintHints()
valueLoop:
	for v := uint(1); v <= 9; v++ {
		hint := uint(1 << (v - 1))
		// cellPeers will contain a list of positions for the given value 'v'
		// and positions visible to it in a container ONLY when that container
		// contains only those two hints
		cellPeers := make(map[int][]int)
		for r := 0; r < 9; r++ {
			for c := 0; c < 9; c++ {
				pos := r*9 + c
				if b.solved[pos] != 0 {
					continue
				}
				if b.blits[pos]&hint == 0 {
					continue
				}

				var links []int
				getSingleLink := func(target int, source int) error {
					if target == source {
						return nil
					}
					if b.solved[target] != 0 {
						return nil
					}
					if b.blits[target]&hint != hint {
						return nil
					}

					links = append(links, target)
					return nil
				}

				allLinks := make(map[int]interface{})

				// row
				if err := b.operateOnRow(pos, getSingleLink); err != nil {
					return err
				}
				if len(links) == 1 {
					for _, item := range links {
						allLinks[item] = struct{}{}
					}
				}

				// column
				links = make([]int, 0)
				if err := b.operateOnColumn(pos, getSingleLink); err != nil {
					return err
				}
				if len(links) == 1 {
					for _, item := range links {
						allLinks[item] = struct{}{}
					}
				}

				// box
				links = make([]int, 0)
				if err := b.operateOnBox(pos, getSingleLink); err != nil {
					return err
				}
				if len(links) == 1 {
					for _, item := range links {
						allLinks[item] = struct{}{}
					}
				} else if len(links) > 1 {
					// delete links if there's more than 1 in a box
					for _, item := range links {
						if _, ok := allLinks[item]; ok {
							//fmt.Printf("-- delete: %#v\n", getCoords(item))
							delete(allLinks, item)
						}
					}
				}

				if len(allLinks) != 0 {
					links = make([]int, 0)
					for k := range allLinks {
						links = append(links, k)
					}
					sort.Ints(links)
					cellPeers[pos] = links
				}
			}
		}

		// we need to consider only contiguous chains
		// it's possible to have two distinct chains with the same hint
		for len(cellPeers) != 0 {
			posColor := make(map[int]int)
			//fmt.Printf("len(cellPeers) = %d, val: %d\n", len(cellPeers), bits.GetSingleBitValue(hint))
			i := 0
			for k, v := range cellPeers {
				color, ok := posColor[k]
				if !ok {
					if i != 0 {
						continue
					}
					color = 0
					posColor[k] = color
				}
				i++

				flippedColor := 1 - color

				for _, peer := range v {
					peerColor, peerOK := posColor[peer]
					if peerOK {
						if peerColor != flippedColor {
							// contradiction
							continue valueLoop
						}
					} else {
						posColor[peer] = flippedColor
					}
				}
			}

			var color0 []int
			var color1 []int
			for k, color := range posColor {
				delete(cellPeers, k)
				if color == 0 {
					color0 = append(color0, k)
				} else {
					color1 = append(color1, k)
				}
			}

			if len(color0) != len(color1) {
				continue
			}

			/*fmt.Printf("* color0: ")
			for _, pos0 := range color0 {
				fmt.Printf("%#v ", getCoords(pos0))
			}
			fmt.Printf("\n* color1: ")
			for _, pos1 := range color1 {
				fmt.Printf("%#v ", getCoords(pos1))
			}
			fmt.Println()*/

			for _, pos0 := range color0 {
				for _, pos1 := range color1 {
					vis0 := b.getVisibleCellsWithHint(pos0, hint)
					vis1 := b.getVisibleCellsWithHint(pos1, hint)

					both := intersect(vis0, vis1)
					both = subtract(both, color0)
					both = subtract(both, color1)

					if len(both) > 0 {
						/*b.PrintHints()
						b.PrintURL()
						fmt.Printf("-- c0=%#v c1=%#v val=%d:\n", getCoords(pos0), getCoords(pos1), bits.GetSingleBitValue(hint))*/
						for _, elem := range both {
							//fmt.Printf("---- eliminate: %#v\n", getCoords(elem))
							if err := b.updateCandidates(elem, pos0, ^hint); err != nil {
								//fmt.Printf("%s\n", err)
								return err
							}
						}
						// let simpler techniques take over
						//fmt.Printf("return\n")
						return nil
					}
				}
			}
		}
	}
	//fmt.Printf("return\n")
	return nil
}
