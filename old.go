// +build !windows

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"
)

func main() {
	boards := getBoards()
	for i, b := range boards {
		fmt.Printf("Board %d/%d\n", i+1, len(boards))
		printBoard(&b)

		fmt.Println("----------------------------")

		for i := uint(0); i < 81; i++ {
			val := b.solved[i]
			if val != 0 {
				err := b.solvePosition(i, val)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
					//log.Fatal(err)
				}
			}
		}

		boards := bruteForce(b)

		if len(boards) == 0 {
			fmt.Printf("*** no solution found\n")
			return
		}

		for _, solved := range boards {
			fmt.Println(solved.numSolved())
			printBoard(&solved)
			if err := solved.validate(); err != nil {
				log.Fatal(err)
			}
			fmt.Println("GOOD!")
		}

	}
}

func getBoards() []board {

	var boards []board
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\r\n")
		if text == "" {
			break
		}
		fmt.Println(text)

		b := board{}

		for j, c := range text {
			if j%2 == 1 {
				if c != ' ' {
					//log.Fatalf("expected space line:%d position:%d char:%#q", i+1, j+1, c)
					log.Fatalf("expected space line:%d position:%d char:%#q", len(boards)+1, j+1, c)
				}
				continue
			}

			//pos := i*9 + j/2
			pos := j / 2
			if unicode.IsDigit(c) {
				n, err := strconv.Atoi(string(c))
				if err != nil {
					log.Fatal(err)
				}
				if n < 1 {
					log.Fatal(n)
				}

				b.solved[pos] = uint(n)
				b.blits[pos] = (1 << uint(n-1))
			} else if c == '_' {
				b.blits[pos] = (1 << 9) - 1
			} else {
				//log.Fatalf("expected digit or _ line:%d position:%d char:%#q", i+1, j+1, c)
				log.Fatalf("expected digit or _ line:%d position:%d char:%#q", len(boards)+1, j+1, c)
			}
		}
		boards = append(boards, b)
		/*if b.blits[80] != 0 {
			break readloop
		}*/
	}
	return boards
}

func (b *board) validate() error {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			var tally []uint
			for rpos := r * 9; rpos < r*9+9; rpos++ {
				if b.solved[rpos] == 0 {
					return fmt.Errorf("pos %d not solved", rpos)
				}
				tally = append(tally, b.solved[rpos])
			}
			if len(tally) != 9 {
				return fmt.Errorf("len(tally) != 9 %+v", tally)
			}
			for i := uint(1); i <= 9; i++ {
				if !contains(tally, i) {
					return fmt.Errorf("tally does not contain %d", i)
				}
			}

			tally = make([]uint, 0)
			for cpos := r % 9; cpos < 81; cpos += 9 {
				if b.solved[cpos] == 0 {
					return fmt.Errorf("pos %d not solved", cpos)
				}
				tally = append(tally, b.solved[cpos])
			}
			if len(tally) != 9 {
				return fmt.Errorf("len(tally) != 9 %+v", tally)
			}
			for i := uint(1); i <= 9; i++ {
				if !contains(tally, i) {
					return fmt.Errorf("tally does not contain %d", i)
				}
			}
		}
	}
	return nil
}

var abort int32

func bruteForce(b board) []board {
	if b.numSolved() == 81 {
		return []board{b}
	}

	minBitsPos := uint(0) // TODO: this should always get set, but... should handle if it doesn't.
	minBits := uint(10)
	for i := uint(0); i < 81; i++ {
		if b.solved[i] != 0 {
			continue
		}

		bits := numberOfSetBits(b.blits[i])
		//fmt.Println(bits)
		if bits < minBits {
			minBits = bits
			minBitsPos = i
		}
	}

	if minBits == 10 {
		log.Fatal("min bits position not found")
	}

	var values []uint
	bits := b.blits[minBitsPos]
	for i := uint(1); i <= 9; i++ {
		if bits&0x01 == 1 {
			values = append(values, i)
		}
		bits >>= 1
	}

	solution := make(chan board)

	var wg sync.WaitGroup
	for _, val := range values {
		//fmt.Printf("%d %d %d %b\n", minBitsPos, val, minBits, b.blits[minBitsPos])
		wg.Add(1)
		go func(b2 board, val uint) {
			defer wg.Done()
			if atomic.LoadInt32(&abort) == 1 {
				return
			}

			err := b2.solvePosition(minBitsPos, val)
			if atomic.LoadInt32(&abort) == 1 {
				return
			}

			if err != nil {
				return
			}
			if b2.numSolved() == 81 {
				solution <- b2
			}

			if atomic.LoadInt32(&abort) == 1 {
				return
			}

			boards := bruteForce(b2)

			for _, solved := range boards {
				if solved.numSolved() == 81 && solved.validate() == nil {
					atomic.StoreInt32(&abort, 1)
					solution <- solved
				}
			}
		}(b, val)
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	var solutions []board
loop:
	//for {
	select {
	case s := <-solution:
		solutions = append(solutions, s)
	case <-done:
		break loop
	}
	//}

	var distinctSolutions []board
	for _, s := range solutions {
		found := false
		for _, e := range distinctSolutions {
			if reflect.DeepEqual(e, s) {
				found = true
				break
			}
		}
		if !found {
			distinctSolutions = append(distinctSolutions, s)
		}
	}

	return distinctSolutions
}

func printBoard(b *board) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			val := b.solved[i*9+j]
			if val != 0 {
				fmt.Print(val)
			} else {
				fmt.Print("_")
			}
			if j != 8 {
				fmt.Print(" ")
			} else {
				fmt.Println("")
			}
		}
	}

	if b.numSolved() == 81 {
		return
	}

	fmt.Println("")

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			val := b.solved[i*9+j]
			if val != 0 {
				fmt.Print(val)
			} else {
				fmt.Printf("%b", b.blits[i*9+j])
			}
			if j != 8 {
				fmt.Print(" ")
			} else {
				fmt.Println("")
			}
		}
	}
}

func (b *board) numSolved() int {
	solved := 0
	for i := 0; i < 81; i++ {
		if b.solved[i] != 0 {
			solved++
		}
	}
	return solved
}

func (b *board) solvePosition(pos uint, val uint) error {
	mask := uint(^(1 << (val - 1)))
	//fmt.Printf("-- r:%d c:%d val:%d\n", pos/9, pos%9, val)
	if b.solved[pos] != 0 && b.solved[pos] != val {
		return fmt.Errorf("pos %d has value %d, tried to set with %d", pos, b.solved[pos], val)
	}
	b.solved[pos] = val
	b.blits[pos] = 1 << (val - 1)
	//b.blits[pos] = 0

	//printBoard(b)
	var changed []uint

	for r := (pos / 9) * 9; r < (pos/9)*9+9; r++ {
		if r == pos || b.solved[r] != 0 {
			continue
		}
		newBlit := b.blits[r] & mask
		if newBlit != b.blits[r] {
			b.blits[r] = newBlit
			if b.solved[r] == 0 && singleBit(b.blits[r]) {
				//fmt.Printf("singlebit r %d %b\n", r, b.blits[r])
				m := getSingleBitValue(b.blits[r])
				if err := b.solvePosition(r, m); err != nil {
					return err
				}
			}
		}
		changed = append(changed, r)
	}
	for c := pos % 9; c < 81; c += 9 {
		if c == pos || b.solved[c] != 0 {
			continue
		}
		newBlit := b.blits[c] & mask
		if newBlit != b.blits[c] {
			b.blits[c] = newBlit
			if b.solved[c] == 0 && singleBit(b.blits[c]) {
				//fmt.Printf("singlebit c %d %b\n", c, b.blits[c])
				m := getSingleBitValue(b.blits[c])
				if err := b.solvePosition(c, m); err != nil {
					return err
				}
			}
		}
		changed = append(changed, c)
	}

	//for _, c := range changed {
	for c := uint(0); c < 81; c++ {
		if b.solved[c] != 0 {
			continue
		}
		if err := b.checkForOnlyBit(c); err != nil {
			return err
		}

		if b.solved[c] != 0 {
			continue
		}
		if err := b.checkForEqualMasks(c); err != nil {
			return err
		}

		//fmt.Printf("%d %b\n", c, b.blits[c])
	}

	return nil

	// if a blit has a flag that's the only one in its row OR column, it must be that value

	//b.checkColumn(pos%9, val)
	//b.checkRow(pos/9, val)
}

func (b *board) checkForEqualMasks(pos uint) error {
	// check for positions with equal masks.
	// if the number of positions with equal masks
	// equals the bits in those masks, the other positions
	// can remove that mask.
	tally := make(map[uint][]uint)
	for r := (pos / 9) * 9; r < (pos/9)*9+9; r++ {
		if b.solved[r] != 0 {
			continue
		}
		mask := b.blits[r]

		if v, ok := tally[mask]; ok {
			tally[mask] = append(v, r)
		} else {
			tally[mask] = []uint{r}
		}
	}

	for k, v := range tally {
		if numberOfSetBits(k) == uint(len(v)) {
			for r := (pos / 9) * 9; r < (pos/9)*9+9; r++ {
				if contains(v, r) {
					continue
				}
				//old := b.blits[r]
				b.blits[r] &= ^k
				/*if old != b.blits[r] {
					fmt.Println("it helped r")
				}*/
			}
		}
	}

	tally = make(map[uint][]uint)
	for c := pos % 9; c < 81; c += 9 {
		if b.solved[c] != 0 {
			continue
		}
		mask := b.blits[c]

		if v, ok := tally[mask]; ok {
			tally[mask] = append(v, c)
		} else {
			tally[mask] = []uint{c}
		}
	}

	for k, v := range tally {
		if numberOfSetBits(k) == uint(len(v)) {
			for c := pos % 9; c < 81; c += 9 {
				if contains(v, c) {
					continue
				}
				//old := b.blits[c]
				b.blits[c] &= ^k
				/*if old != b.blits[c] {
					fmt.Println("it helped c")
				}*/
			}
		}
	}

	return nil
}

func contains(list []uint, val uint) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

func (b *board) checkForOnlyBit(pos uint) error {
	//fmt.Printf("checkForOnlyBit r:%d c:%d\n", pos/9, pos%9)
	var mask uint
	for r := (pos / 9) * 9; r < (pos/9)*9+9; r++ {
		if r == pos {
			continue
		}
		mask |= b.blits[r]
	}
	mask = (b.blits[pos] ^ mask) & b.blits[pos]
	if singleBit(mask) {
		val := getSingleBitValue(mask)
		//fmt.Printf("r only bit - %d\n", val)
		if err := b.solvePosition(pos, val); err != nil {
			return err
		}
	}

	mask = 0
	for c := pos % 9; c < 81; c += 9 {
		if c == pos {
			continue
		}
		mask |= b.blits[c]
	}
	mask = (b.blits[pos] ^ mask) & b.blits[pos]
	if singleBit(mask) {
		val := getSingleBitValue(mask)
		//fmt.Printf("c only bit - %d\n", val)
		if err := b.solvePosition(pos, val); err != nil {
			return err
		}
	}

	return nil
}
