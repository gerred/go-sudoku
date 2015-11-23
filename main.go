package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flags := flag.FlagSet{}
	profile := flags.Bool("profile", false, "profile cpu/mem")
	runFile := flags.String("file", "", "file containing puzzle(s)")
	maxPuzzles := flags.Int("max-puzzles", -1, "max puzzles to solve when multiple present in a file")
	//readStats := flags.String("read-stats", "", "read stats from long run, print time taken per puzzle")
	showSteps := flags.Bool("steps", false, "show solve steps")
	showSolveTime := flags.Bool("time", false, "print time taken to solve")

	var err error
	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.PrintDefaults()
		log.Fatal(err)
	}

	_ = maxPuzzles

	/*if *readStats != "" {
		err := readStatsFile(*readStats)
		if err != nil {
			log.Fatal(err)
		}
		return
	}*/

	if *profile {
		if err = startProfiler(); err != nil {
			log.Fatal(err)
		}
	}

	var start time.Time
	defer func() {
		if *showSolveTime {
			fmt.Printf("time=%v\n", time.Since(start))
		}
		if *profile {
			if err := stopProfiler(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	if *runFile == "" {
		// read board from stdin (before starting timer)
		var boardBytes []byte
		if boardBytes, err = readBoard(os.Stdin); err != nil {
			log.Fatal(err)
		}

		start = time.Now()

		var b *board
		if b, err = loadBoard(boardBytes); err != nil {
			if _, ok := err.(ErrUnsolvable); ok {
				fmt.Printf("UNSOLVABLE\n")
				return
			}
			log.Fatal(err)
		}

		if *showSteps {
			b.showSteps = true
		}

		if err = b.Solve(); err != nil {
			if _, ok := err.(ErrUnsolvable); ok {
				fmt.Printf("UNSOLVABLE\n")
				return
			}
			log.Fatal(err)
		}

		b.Print()
	} else {
		// read compact board(s) from file
		start = time.Now()
		if err := runList(*runFile, *maxPuzzles, *showSolveTime, *showSteps); err != nil {
			log.Fatal(err)
		}
	}
}

func readStatsFile(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	const prefix = "Solve time: "

	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}

	puzzle := 1
	for line != "" {
		if strings.HasPrefix(line, prefix) {
			line = strings.Trim(line[len(prefix):], " \n\r")
			var d time.Duration
			if d, err = time.ParseDuration(line); err != nil {
				return err
			}
			fmt.Printf("%d\t%v\n", puzzle, d.Nanoseconds()/int64(time.Millisecond))
			puzzle++
		}
		if line, err = r.ReadString('\n'); err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

func runFile(fileName string) {
	board, err := getBoard(fileName)
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
		return
	}

	if err = board.Solve(); err != nil {
		fmt.Printf("ERROR - %s\n", err)
		return
	}

	if !board.isSolved() {
		board.PrintHints()
		board.PrintCompact()
		fmt.Println("could not solve")
	} else {
		board.Print()
	}
}

func runList(fileName string, maxPuzzles int, showSolveTime, showSteps bool) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	r := bufio.NewReader(bytes.NewReader(b))
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}
	for i := 0; line != "" && (maxPuzzles == -1 || i < maxPuzzles); i++ {
		firstLog = true
		logLastStepReducedHints = false
		logLastBoardWithHints = ""
		logLastHeader = ""

		fmt.Printf("-----------------\nPuzzle # %d:\n", i+1)
		start1 := time.Now()
		board, err := loadBoard([]byte(line))
		if err != nil {
			if board != nil {
				board.PrintHints()
			}
			return fmt.Errorf("puz=%d err=%q", i+1, err)
		}

		board.showSteps = showSteps

		if err = board.Solve(); err != nil {
			fmt.Printf("%s\n", line)
			b2, err2 := loadBoard([]byte(line))
			if err2 != nil {
				b2.PrintCompact()
			}
			return fmt.Errorf("puz=%d err=%q", i+1, err)
		}

		if !board.isSolved() {
			board.PrintHints()
			board.PrintCompact()
			return fmt.Errorf("could not solve Puzzle # %d", i+1)
		}

		board.Print()

		if showSolveTime {
			fmt.Printf("puz=%d time=%v\n", i+1, time.Since(start1))
		}

		line, err = r.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}
