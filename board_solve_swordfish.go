package main

type swordfishOperation struct {
	blockType                    string
	op                           containerOperator
	opInverted                   containerOperator
	nextContainer                func(int) int
	doesOverlap                  func(x int, y int) bool
	isInSameDimension            func(x int, y int) bool
	maxMisses                    func(overlapSet []int) int
	getMaxPosFromOverlapSet      func(overlapSet []int) int
	getInvertedDimensionPosition func(pos int) int
}

func (b *board) SolveSwordFish() error {
	// http://www.sudokuwiki.org/Sword_Fish_Strategy
	// find a 3x3 which share the same candidate
	// hint cannot be repeated on container (row or col depending on orientation)
	// 333,332,322,222 = all valid.

	// create dimensions for looking for SwordFish in the row and column dimension
	dims := []swordfishOperation{
		{
			blockType:  "row",
			op:         b.operateOnRow,
			opInverted: b.operateOnColumn,
			nextContainer: func(cur int) int {
				// next start row index
				next := cur + 9
				if next >= 81 {
					return -1
				}
				return next
			},
			doesOverlap: func(x int, y int) bool {
				coords1 := getCoords(x)
				coords2 := getCoords(y)
				return coords1.col == coords2.col
			},
			isInSameDimension: func(x int, y int) bool {
				coords1 := getCoords(x)
				coords2 := getCoords(y)
				return coords1.row == coords2.row
			},
			maxMisses: func(overlapSet []int) int {
				// count distinct cols
				cols := make(map[int]interface{})
				for _, pos := range overlapSet {
					cols[getCoords(pos).col] = struct{}{}
				}
				return 3 - len(cols)
			},
			getMaxPosFromOverlapSet: func(overlapSet []int) int {
				maxRow := 0
				for _, pos := range overlapSet {
					coords := getCoords(pos)
					if coords.row > maxRow {
						maxRow = coords.row
					}
				}
				return maxRow * 9
			},
			getInvertedDimensionPosition: func(pos int) int {
				return getCoords(pos).col
			},
		},
		{
			blockType:  "column",
			op:         b.operateOnColumn,
			opInverted: b.operateOnRow,
			nextContainer: func(cur int) int {
				// next start column index
				next := cur + 1
				if next >= 9 {
					return -1
				}
				return next
			},
			doesOverlap: func(x int, y int) bool {
				coords1 := getCoords(x)
				coords2 := getCoords(y)
				return coords1.row == coords2.row
			},
			isInSameDimension: func(x int, y int) bool {
				coords1 := getCoords(x)
				coords2 := getCoords(y)
				return coords1.col == coords2.col
			},
			maxMisses: func(overlapSet []int) int {
				// count distinct rows
				rows := make(map[int]interface{})
				for _, pos := range overlapSet {
					rows[getCoords(pos).row] = struct{}{}
				}
				return 3 - len(rows)
			},
			getMaxPosFromOverlapSet: func(overlapSet []int) int {
				maxCol := 0
				for _, pos := range overlapSet {
					coords := getCoords(pos)
					if coords.col > maxCol {
						maxCol = coords.col
					}
				}
				return maxCol
			},
			getInvertedDimensionPosition: func(pos int) int {
				return getCoords(pos).row
			},
		},
	}

	for _, dim := range dims {

		// TODO: DEBUG, remove
		/*if dim.blockType != "column" {
			continue
		}*/
		// END debug

		for i := 0; i != -1; i = dim.nextContainer(i) {
			// i = start of row/column

			// TODO: debug, remove
			/*if i != 2 {
				continue
			}*/
			// END debug

			// hint, list of positions
			initialCandidates := make(map[uint][]int)

			// operate on all cells in container
			// get list of candidate cells: map(hint, []pos)
			extractHints := func(target int, source int) error {
				bitList := GetBitList(b.blits[target])
				if len(bitList) < 2 {
					return nil
				}

				for _, bit := range bitList {
					list, ok := initialCandidates[bit]
					if !ok {
						initialCandidates[bit] = []int{target}
					} else {
						initialCandidates[bit] = append(list, target)
					}
				}
				return nil
			}

			if err := dim.op(i, extractHints); err != nil {
				return err
			}

			// get a list of hints containing permutations of 2 or 3 cells
			candidatePerms := swordfishGetCandidatePermutations(initialCandidates)
			if len(candidatePerms) == 0 {
				continue
			}

			//fmt.Printf("%#v\n", candidatePerms)

			for hint, v := range candidatePerms {
				// TODO: REMOVE, debug
				/*if hint != 1 {
					continue
				}*/
				// END debug

				/*var once1 sync.Once
				print1 := func() {
					fmt.Printf("-----------------------------------------\n")
					fmt.Printf("swordfish: dim:%s num:%d\n", dim.blockType, bits.GetSingleBitValue(hint))
				}*/

				for _, c := range v {
					// TODO: DEBUG, remove
					/*if !reflect.DeepEqual(c, []int{2, 11}) {
						continue
					}
					fmt.Println("GOT IT")*/
					// END debug

					/*var once2 sync.Once
					print2 := func() {
						once1.Do(print1)
						fmt.Printf("\n")
						fmt.Printf("-- top level:\n")
					}*/

					sfPerms, err := b.swordfishGetPermutations(dim, c, hint, i)
					if err != nil {
						return err
					}

					/*var once3 sync.Once
					print3 := func() {
						once2.Do(print2)
						for _, pick := range c {
							fmt.Printf("-- %#2v %s\n", getCoords(pick), bits.GetString(b.blits[pick]))
						}
					}*/

					for _, x := range sfPerms {
						// TODO: DEBUG, remove
						/*if !reflect.DeepEqual(x, []int{3, 12, 30}) {
							continue
						}*/
						// END debug

						/*var once4 sync.Once
						print4 := func() {
							once3.Do(print3)
							fmt.Printf("---- second level:\n")
							for _, pos := range x {
								fmt.Printf("---- %#2v %s\n", getCoords(pos), bits.GetString(b.blits[pos]))
							}
						}*/

						// combine first and second level to get new 'one must overlap' set
						// (in this case, two must overlap...)
						var overlapSet []int
						overlapSet = append(overlapSet, c...)
						overlapSet = append(overlapSet, x...)

						//fmt.Printf("%#v\n", overlapSet)

						nextPos := dim.getMaxPosFromOverlapSet(overlapSet)
						sfPerms2, err := b.swordfishGetPermutations(dim, overlapSet, hint, nextPos)
						if err != nil {
							return err
						}

						if len(sfPerms2) != 0 {
							//once4.Do(print4)
							for _, y := range sfPerms2 {
								// TODO: DEBUG, remove
								/*if !reflect.DeepEqual(y, []int{16, 34}) {
									continue
								}*/
								// END debug

								/*fmt.Printf("------ third level:\n")
								for _, pos := range y {
									fmt.Printf("------ %#2v %s\n", getCoords(pos), bits.GetString(b.blits[pos]))
								}*/

								// Here we go...
								// if we pulled sets from columns, given a hint, remove that hint
								// from all cells in the superset of ROWS covered by our swordfish
								// set. still confused? me too.
								// http://www.sudokuwiki.org/Sword_Fish_Strategy

								err := b.swordfishApply(dim, hint, c, x, y)
								if err != nil {
									return err
								}
								if b.changed {
									// try simpler techniques before re-trying swordfish
									//fmt.Println("**** we did it!")
									return nil
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (b *board) swordfishGetPermutations(dim swordfishOperation, mustOverlap []int, hint uint, pos int) ([][]int, error) {
	var perms [][]int

	// make sure we're still on the board
	nextPos := dim.nextContainer(pos)
	if nextPos == -1 {
		return perms, nil
	}

	var candidates []int
	getCandidates := func(target int, source int) error {
		// cell must contain the hint we're looking for
		if b.blits[target]&hint != hint {
			return nil
		}

		candidates = append(candidates, target)
		return nil
	}

	if err := dim.op(nextPos, getCandidates); err != nil {
		return perms, err
	}

	tmpPerms := swordfishGetTwosAndThrees(candidates)

	// TODO: debug, look for dupes between mustOverlap and tmpPerms
	/*fmt.Printf("--------------------\n")
	fmt.Printf("blockType:   %s\n", dim.blockType)
	fmt.Printf("mustOverlap: %#2v\n", mustOverlap)
	fmt.Printf("tmpPerms:    %#2v\n", tmpPerms)
	fmt.Printf("pos:         %d\n", pos)
	fmt.Printf("nextPos:     %d\n", nextPos)
	fmt.Printf("--------------------\n")*/
	// END debug

	maxMisses := dim.maxMisses(mustOverlap)
	if maxMisses < 0 {
		return perms, nil
	}

	for _, perm := range tmpPerms {
		misses := 0
		for _, item := range perm {
			doesOverlap := false
			for _, existingSetItem := range mustOverlap {
				if dim.doesOverlap(item, existingSetItem) {
					doesOverlap = true
					break
				}
			}
			if !doesOverlap {
				misses++
				if misses > maxMisses {
					break
				}
			}
		}
		if misses <= maxMisses {
			perms = append(perms, perm)
		}
	}

	permsN, err := b.swordfishGetPermutations(dim, mustOverlap, hint, nextPos)
	if err != nil {
		return perms, err
	}

	perms = append(perms, permsN...)

	return perms, nil
}

func swordfishGetCandidatePermutations(orig map[uint][]int) map[uint][][]int {
	filtered := make(map[uint][][]int)
	for k, v := range orig {
		list := swordfishGetTwosAndThrees(v)
		if len(list) != 0 {
			filtered[k] = list
		}
	}
	return filtered
}

func swordfishGetTwosAndThrees(v []int) [][]int {
	var emptyList [][]int

	fours := getPermutations(4, v, []int{})
	if len(fours) != 0 {
		return emptyList
	}

	threes := getPermutations(3, v, []int{})
	if len(threes) != 0 {
		return threes
	}

	twos := getPermutations(2, v, []int{})
	if len(twos) != 0 {
		return twos
	}

	return emptyList
}

func (b *board) swordfishApply(sf swordfishOperation, hint uint, set1 []int, set2 []int, set3 []int) error {
	var overlap []int
	overlap = append(overlap, set1...)
	overlap = append(overlap, set2...)
	overlap = append(overlap, set3...)

	dupeCheck := make(map[int]interface{})
	for _, item := range overlap {
		key := sf.getInvertedDimensionPosition(item)
		_, found := dupeCheck[key]
		if found {
			continue
		}
		dupeCheck[key] = struct{}{}

		removeHint := func(target int, source int) error {
			for _, pos := range overlap {
				if sf.isInSameDimension(target, pos) {
					return nil
				}
			}

			return b.updateCandidates(target, source, ^hint)
		}

		if err := sf.opInverted(item, removeHint); err != nil {
			return err
		}
	}

	return nil
}
