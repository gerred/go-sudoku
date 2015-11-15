package main

import (
	"bytes"
	"fmt"
)

func getRCV(r int, c int, v int) int {
	return r*100 + c*10 + v
}

func (b *board) getSAT() string {
	var clauses int
	vars := make(map[int]interface{})

	buf := bytes.NewBufferString("")

	// apply known values

	// each row
	for r := 1; r <= 9; r++ {
		// each column
		for c := 1; c <= 9; c++ {
			offset := (r-1)*9 + (c - 1)
			v := int(b.solved[offset])
			if v != 0 {
				cur := getRCV(r, c, v)
				buf.WriteString(fmt.Sprintf("%d 0\n", cur))
				clauses++
			}
		}
	}

	// each value
	for v := 1; v <= 9; v++ {
		// each row
		for r := 1; r <= 9; r++ {
			// each column
			for c := 1; c <= 9; c++ {
				cur := getRCV(r, c, v)
				vars[cur] = struct{}{}

				buf.WriteString(fmt.Sprintf("%d ", cur))
			}
			buf.WriteString("0\n")
			clauses++
			// each combination of two values
			for c := 1; c <= 9; c++ {
				cur := getRCV(r, c, v)
				for c2 := c + 1; c2 <= 9; c2++ {
					cur2 := getRCV(r, c2, v)

					buf.WriteString(fmt.Sprintf("-%d -%d 0\n", cur, cur2))
					clauses++
				}
			}
		}

		// each column
		for c := 1; c <= 9; c++ {
			// each row
			for r := 1; r <= 9; r++ {
				cur := getRCV(r, c, v)

				buf.WriteString(fmt.Sprintf("%d ", cur))
			}
			buf.WriteString("0\n")
			clauses++
			// each combination of two values
			for r := 1; r <= 9; r++ {
				cur := getRCV(r, c, v)
				for r2 := r + 1; r2 <= 9; r2++ {
					cur2 := getRCV(r2, c, v)

					buf.WriteString(fmt.Sprintf("-%d -%d 0\n", cur, cur2))
					clauses++
				}
			}
		}

		// each box
		for b := 0; b < 9; b++ {
			rOffset := (b / 3) * 3
			cOffset := (b % 3) * 3
			//fmt.Printf("b:%d r-offset:%d c-offset:%d\n", b, rOffset, cOffset)
			// each row
			for r := rOffset + 1; r <= rOffset+3; r++ {
				// each column
				for c := cOffset + 1; c <= cOffset+3; c++ {
					cur := getRCV(r, c, v)

					buf.WriteString(fmt.Sprintf("%d ", cur))
				}
			}
			buf.WriteString("0\n")
			clauses++
			// each combination of two values
			// each row
			cheat := make(map[string]interface{})
			for r := rOffset + 1; r <= rOffset+3; r++ {
				// each column
				for c := cOffset + 1; c <= cOffset+3; c++ {
					cur := getRCV(r, c, v)
					// TODO: we really just need to check the diagnals
					// and add 4 additional clauses per box
					for r2 := rOffset + 1; r2 <= rOffset+3; r2++ {
						for c2 := cOffset + 1; c2 <= cOffset+3; c2++ {
							// already checked row/col constraint
							if r == r2 || c == c2 {
								continue
							}
							cur2 := getRCV(r2, c2, v)

							clause := fmt.Sprintf("-%d -%d 0\n", cur, cur2)
							clause_inv := fmt.Sprintf("-%d -%d 0\n", cur2, cur)
							if _, ok := cheat[clause_inv]; ok {
								continue
							}
							cheat[clause] = struct{}{}

							buf.WriteString(clause)
							clauses++
						}
					}
				}
			}
		}
	}

	// each row
	for r := 1; r <= 9; r++ {
		// each column
		for c := 1; c <= 9; c++ {
			// each value {
			for v1 := 1; v1 <= 9; v1++ {
				cur := getRCV(r, c, v1)
				for v2 := v1 + 1; v2 <= 9; v2++ {
					cur2 := getRCV(r, c, v2)

					buf.WriteString(fmt.Sprintf("-%d -%d 0\n", cur, cur2))
					clauses++
				}
			}
		}
	}

	header := fmt.Sprintf("p cnf %d %d", len(vars), clauses)
	//fmt.Printf("%s\n", header)
	input := fmt.Sprintf("%s\n%s", header, buf)
	return input
}
