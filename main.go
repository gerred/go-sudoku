package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/judwhite/go-sudoku/internal/sat"
)

func main() {
	flags := flag.FlagSet{}
	profile := flags.Bool("profile", false, "true profile cpu/mem")
	max_iterations := flags.Int("run", -1, "max iterations")
	read_stats := flags.String("readstats", "", "read stats from long run, print time taken per puzzle")
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	_ = max_iterations

	if *read_stats != "" {
		err := readStats(*read_stats)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	start := time.Now()
	if *profile {
		startProfile()
	}

	//printCompactToStandard("210300000000060050000000000300000702004050000000000100000102080036000000000700000")

	//runFile("./test_files/12_tough_20151107_173.txt")

	// UNSAT: 020400006400089000000007004001008060000700008030060500060000010005000300910800007
	// SAT: 020400000400089000000007004001008060000700008030060500060000010005000300910800007
	b, _ := loadBoard([]byte("123456789456000000789000000000000000000000000000000000000000000000000000000000000"))
	b.Solve()
	b.Print()

	//start := time.Now()
	//runFile("./test_files/29_ben.txt")
	//runFile("./test_files/12_tough_20151107_173.txt")
	//runFile("./test_files/input_no_solution.txt")
	//runList("./test_files/top95.txt", *max_iterations)
	//runList("./test_files/sudoku17.txt", *max_iterations)

	if *profile {
		stopProfile()
	}
	fmt.Printf("%v\n", time.Since(start))
}

func startProfile() {
	f, err := os.Create("go-sudoku.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
}

func stopProfile() {
	f2, err := os.Create("go-sudoku.mprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f2)
	f2.Close()

	pprof.StopCPUProfile()
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

func readStats(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')
	const prefix = "Solve time: "
	puzzle := 1
	for line != "" && err == nil {
		if strings.HasPrefix(line, prefix) {
			line = strings.Trim(line[len(prefix):], " \n\r")
			d, err := time.ParseDuration(line)
			if err != nil {
				return err
			}
			fmt.Printf("%d\t%v\n", puzzle, d.Nanoseconds()/int64(time.Millisecond))
			puzzle++
		}
		line, err = r.ReadString('\n')
	}
	return nil
}

func (b *board) SolveSAT() error {
	satInput := b.getSAT()
	satSolver, err := sat.NewSAT(satInput)
	if err != nil {
		return err
	}
	satSolver = satSolver.Solve()
	if satSolver == nil {
		return fmt.Errorf("could not solve with SAT\n")
	} else {
		fmt.Printf("solved with SAT\n")

		for _, setvar := range satSolver.SetVars {
			k := int(setvar.VarNum)
			v := setvar.Value
			if v {
				//fmt.Printf("%d %v\n", k, v)
				r := k/100 - 1
				c := (k%100)/10 - 1
				pos := r*9 + c
				if b.solved[pos] == 0 {
					val := k % 10
					//fmt.Printf("r:%d c:%d val:%d\n", r, c, val)
					b.SolvePositionNoValidate(pos, uint(val))
				}
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

	err = board.SolveSAT()
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
		return
	}

	/*if err = board.Solve(); err != nil {
		fmt.Printf("ERROR - %s\n", err)
		board.PrintHints()
		return
	}*/

	if !board.isSolved() {
		board.PrintHints()
		board.PrintURL()
		fmt.Println("could not solve")
	} else {
		board.Print()
	}
}

func runList(fileName string, max_iterations int) {
	b, err := ioutil.ReadFile(fileName)
	//b, err := ioutil.ReadFile("./test_files/sudoku17.txt")
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
		return
	}

	r := bufio.NewReader(bytes.NewReader(b))
	line, _ := r.ReadString('\n')
	for i := 0; line != "" && (max_iterations == -1 || i < max_iterations); i++ {
		fmt.Printf("----------------\nPuzzle # %d\n", i+1)
		start1 := time.Now()
		board, err := loadBoard([]byte(line))
		if err != nil {
			board.PrintHints()
			fmt.Printf("ERROR - %d - %s\n", i+1, err)
			return
		}

		if err = board.Solve(); err != nil {
			fmt.Printf("ERROR - %s\n", err)
			return
		}

		/*if err = board.Solve(); err != nil {
			board.PrintHints()
			fmt.Printf("ERROR - %d - %s\n", i+1, err)
			return
		}*/

		if !board.isSolved() {
			board.PrintHints()
			board.PrintURL()
			fmt.Printf("could not solve - %d\n", i+1)
			break
		} else {
			board.Print()
			fmt.Printf("Solve time: %v\n", time.Since(start1))
		}

		line, _ = r.ReadString('\n')
	}
}
