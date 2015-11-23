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
	"runtime/pprof"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flags := flag.FlagSet{}
	profile := flags.Bool("profile", false, "true profile cpu/mem")
	runFile := flags.String("file", "", "file containing puzzle(s)")
	maxPuzzles := flags.Int("max-puzzles", -1, "max puzzles to solve when multiple present in a file")
	readStats := flags.String("read-stats", "", "read stats from long run, print time taken per puzzle")
	printTime := flags.Bool("time", false, "print time taken to solve")

	var err error
	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.PrintDefaults()
		log.Fatal(err)
	}

	_ = maxPuzzles

	if *readStats != "" {
		err := readStatsFile(*readStats)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	//printCompactToStandard("080009743050008010010000000800005000000804000000300006000000070030500080972400050")

	//runFile("./test_files/12_tough_20151107_173.txt")
	// 000000000000000000000000000000000000000000000000000000000000000000000000000000000
	// diabolical: 074302000000005040000607900056000790300000005027000680005701000010200000000408160
	// 8 SLNS: 080009743050008010010000000800005000000804000000300006000000070030500080972400050
	// UNSAT: 020400006400089000000007004001008060000700008030060500060000010005000300910800007
	// SAT: 020400000400089000000007004001008060000700008030060500060000010005000300910800007
	/*b, _ := loadBoard([]byte(`052400000000070100000000000000802000300000600090500000106030000000000089700000000`))
	//b.CountSolutions = true
	//b.MaxSolutions = 500
	b.PrintURL()
	kb, _ := loadBoard([]byte("652481937834679152971325864467812593315794628298563471186937245523146789749258316"))
	var ka [81]byte
	for i, v := range kb.solved {
		ka[i] = byte(v)
	}
	b.knownAnswer = &ka
	err := b.Solve()
	if err != nil {
		log.Fatal(err)
	}
	b.Print()*/

	if *runFile != "" {
		if err := runList(*runFile, *maxPuzzles); err != nil {
			log.Fatal(err)
		}
	}

	if *profile {
		if err = startProfile(); err != nil {
			log.Fatal(err)
		}
	}

	// read board from stdin before starting timer
	var boardBytes []byte
	if boardBytes, err = readBoard(os.Stdin); err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	defer func() {
		if *printTime {
			fmt.Printf("%v\n", time.Since(start))
		}
	}()

	var b *board
	if b, err = loadBoard(boardBytes); err != nil {
		if _, ok := err.(ErrUnsolvable); ok {
			fmt.Println("UNSOLVABLE")
			return
		}
		log.Fatal(err)
	}

	if err = b.Solve(); err != nil {
		if _, ok := err.(ErrUnsolvable); ok {
			fmt.Println("UNSOLVABLE")
			return
		}
		log.Fatal(err)
	}

	b.Print()

	//generate()

	if *profile {
		if err := stopProfile(); err != nil {
			log.Fatal(err)
		}
	}
}

func startProfile() error {
	f, err := os.Create("go-sudoku.pprof")
	if err != nil {
		return err
	}
	if err = pprof.StartCPUProfile(f); err != nil {
		return err
	}
	return nil
}

func stopProfile() error {
	f, err := os.Create("go-sudoku.mprof")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		return err
	}

	pprof.StopCPUProfile()

	return nil
}

func printCompactToStandard(b string) {
	i := 0
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if c != 0 {
				fmt.Print(" ")
			}
			if b[i] == '0' {
				fmt.Print("_")
			} else {
				fmt.Print(string(b[i]))
			}
			i++
		}
		fmt.Println()
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

func (b *board) SolveSAT() error {
	satInput := b.getSAT()
	satSolver, err := NewSAT(satInput, b.CountSolutions, b.MaxSolutions)
	if err != nil {
		return err
	}

	slns := satSolver.Solve()
	if slns == nil || len(slns) == 0 {
		return fmt.Errorf("could not solve with SAT %v\n", slns)
	}

	if !b.CountSolutions {
		//fmt.Printf("solved with SAT\n")
	} else {
		b.SolutionCount = len(slns)
		//fmt.Printf("solved with SAT. solution count: %d\n", len(slns))
	}

	sln1 := slns[0]
	for _, setvar := range sln1.SetVars {
		k := int(setvar.VarNum)
		v := setvar.Value
		if v {
			r := k/100 - 1
			c := (k%100)/10 - 1
			pos := r*9 + c
			if b.solved[pos] == 0 {
				val := k % 10
				b.SolvePositionNoValidate(pos, uint(val))
			}
		}
	}

	err = b.Validate()
	if err != nil {
		return err
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

func runList(fileName string, maxPuzzles int) error {
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
		fmt.Printf("----------------\nPuzzle # %d\n", i+1)
		start1 := time.Now()
		board, err := loadBoard([]byte(line))
		if err != nil {
			board.PrintHints()
			return fmt.Errorf("%s. Puzzle #: %d", err, i+1)
		}

		if err = board.Solve(); err != nil {
			fmt.Printf("%s\n", line)
			b2, err2 := loadBoard([]byte(line))
			if err2 != nil {
				b2.PrintCompact()
			}
			return fmt.Errorf("%s. Puzzle # %d", err, i+1)
		}

		if !board.isSolved() {
			board.PrintHints()
			board.PrintCompact()
			return fmt.Errorf("could not solve Puzzle # %d\n", i+1)
		}

		board.Print()
		fmt.Printf("Solve time: %v\n", time.Since(start1))

		line, err = r.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}
