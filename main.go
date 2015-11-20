package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/judwhite/go-sudoku/internal/bits"
	"github.com/judwhite/go-sudoku/internal/sat"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	/*puzzles := bytes.Buffer{}

	for i := 174; i > 0; i-- {
		fmt.Printf("%d\n", i)
		resp, err := http.DefaultClient.Get("http://www.sudokuwiki.org/feed/scanraid/ASSudokuWeekly.asp?wp=" + strconv.Itoa(i))
		if err != nil {
			log.Fatal(err)
		}
		r := bufio.NewReader(resp.Body)
		find := "load_from_script(false,'e"
		for {
			line, err := r.ReadString('\n')
			if line == "" || err != nil {
				break
			}
			idx := strings.Index(line, find)
			if idx != -1 {
				line = line[idx+len(find) : len(line)-7]
				_, err = puzzles.WriteString(line + "\n")
				if err != nil {
					log.Fatal(err)
				}

				break
			}
		}
	}

	err := ioutil.WriteFile("weekly_unsolvable.txt", puzzles.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	return*/

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

	//printCompactToStandard("080009743050008010010000000800005000000804000000300006000000070030500080972400050")

	//runFile("./test_files/12_tough_20151107_173.txt")
	// 000000000000000000000000000000000000000000000000000000000000000000000000000000000
	// 8 SLNS: 080009743050008010010000000800005000000804000000300006000000070030500080972400050
	// UNSAT: 020400006400089000000007004001008060000700008030060500060000010005000300910800007
	// SAT: 020400000400089000000007004001008060000700008030060500060000010005000300910800007
	/*b, _ := loadBoard([]byte(`097003420604020000300700000080407000503002006010086040400005008000170000000040270`))
	b.PrintURL()
	err := b.Solve()
	if err != nil {
		log.Fatal(err)
	}
	b.Print()*/

	//start := time.Now()
	//runFile("./test_files/29_ben.txt")
	//runFile("./test_files/12_tough_20151107_173.txt")
	//runFile("./test_files/input.txt")
	//runFile("./test_files/input_no_solution.txt")
	//runList("./test_files/top95.txt", *max_iterations)
	//runList("./test_files/weekly_unsolvable.txt", *max_iterations)
	//runList("./test_files/sudokus.txt", *max_iterations)

	generate()

	if *profile {
		stopProfile()
	}
	fmt.Printf("%v\n", time.Since(start))
}

func getValidBoard() (*board, error) {
	b, err := loadBoard([]byte("000000000000000000000000000000000000000000000000000000000000000000000000000000000"))
	if err != nil {
		return nil, err
	}
	/*for i := 0; i < 8; i++ {
		err := b.SolvePosition(i, uint(i+1))
		if err != nil {
			return nil, err
		}
	}*/
	for !b.isSolved() {
		n := rand.Intn(81)
		if b.solved[n] != 0 {
			continue
		}
		bitList := bits.GetBitList(b.blits[n])
		bn := rand.Intn(len(bitList))
		val := bits.GetSingleBitValue(bitList[bn])

		err = b.SolvePosition(n, val)
		if err != nil {
			return nil, err
		}

		err = b.Solve()
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func generate() {
	var err error
	var b *board
	for b == nil || err != nil {
		b, err = getValidBoard()
	}

	//b.PrintHints()

	digHoles(b)

	/*b.CountSolutions = true

	err = b.SolveSAT()
	b.Print()
	if err != nil {
		log.Fatal(err)
	}*/
}

func digHoles(b *board) error {
	var err error
	var b2 *board
restartLoop:
	for {
		b2, err = loadBoard([]byte("000000000000000000000000000000000000000000000000000000000000000000000000000000000"))
		if err != nil {
			return err
		}

		// solve digits 1-9, minimum guarantee for a unique solution
		j := 1
		for b2.numSolved() != 9 {
			n := rand.Intn(81)
			if b2.solved[n] != 0 || b.solved[n] != uint(j) {
				continue
			}
			j++
			err = b2.SolvePosition(n, b.solved[n])
			if err != nil {
				return err
			}
		}

		// solve to 17 cells
		for b2.numSolved() < 17 {
			b2.Print()
			fmt.Printf("num solved: %d\n", b2.numSolved())
			check := make(map[int]interface{})
			minslncount := -1
			minslncount_n := -1
			for {
				if len(check) == 81 {
					break
				}

				n := rand.Intn(81)
				if _, ok := check[n]; ok {
					continue
				}
				check[n] = struct{}{}

				if b2.solved[n] != 0 {
					continue
				}

				b3 := &board{solved: b2.solved, blits: b2.blits, CountSolutions: true, MaxSolutions: 100}

				// see if the sln count has changed
				err = b3.SolvePosition(n, b.solved[n])
				if err != nil {
					return err
				}
				if err = b3.SolveSAT(); err != nil {
					return err
				}

				if minslncount == -1 || minslncount > b3.SolutionCount {
					minslncount = b3.SolutionCount
					minslncount_n = n
				}
			}

			if minslncount == -1 {
				continue restartLoop
			}

			// solve the cell in b2
			if minslncount != 0 {
				err = b2.SolvePosition(minslncount_n, b.solved[minslncount_n])
				if err != nil {
					return err
				}
			}
		}

		// solve the rest needed for a unique solution
		//var oldslncount int
		for {
			b2.Print()
			fmt.Printf("num solved: %d\n", b2.numSolved())
			if b2.numSolved() > 30 {
				fmt.Printf("restartd, %d solved..\n", b2.numSolved())
				continue restartLoop
			}

			b3 := &board{solved: b2.solved, blits: b2.blits, CountSolutions: true, MaxSolutions: 500}
			err = b3.Solve()
			if err != nil {
				continue restartLoop
			}

			//oldslncount = b3.SolutionCount

			if err == nil && b3.isSolved() && b3.SolutionCount <= 1 {
				break restartLoop
			}

			check := make(map[int]interface{})
			minslncount := -1
			minslncount_n := -1
			for {
				if len(check) == 81 {
					break
				}

				n := rand.Intn(81)
				if _, ok := check[n]; ok {
					continue
				}
				check[n] = struct{}{}

				if b2.solved[n] != 0 {
					continue
				}

				b3 := &board{solved: b2.solved, blits: b2.blits, CountSolutions: true, MaxSolutions: 500}

				// see if the sln count has changed
				err = b3.SolvePosition(n, b.solved[n])
				if err != nil {
					return err
				}
				if err = b3.SolveSAT(); err != nil {
					return err
				}

				if minslncount == -1 || minslncount > b3.SolutionCount {
					minslncount = b3.SolutionCount
					minslncount_n = n
				}
			}

			if minslncount == -1 {
				continue restartLoop
			}

			// solve the cell in b2
			err = b2.SolvePosition(minslncount_n, b.solved[minslncount_n])
			if err != nil {
				return err
			}
		}
	}

	b3 := &board{solved: b2.solved, blits: b2.blits, SkipSAT: true}
	b3.Print()
	b3.PrintURL()
	fmt.Printf("num solved: %d\n", b3.numSolved())
	err = b3.Solve()
	if err != nil {
		log.Fatal(err)
	}
	b3.Print()

	return nil
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
	satSolver, err := sat.NewSAT(satInput, b.CountSolutions, b.MaxSolutions)
	if err != nil {
		return err
	}
	slns := satSolver.Solve()
	if slns == nil || len(slns) == 0 {
		return fmt.Errorf("could not solve with SAT %v\n", slns)
	} else {
		if !b.CountSolutions {
			fmt.Printf("solved with SAT\n")
		} else {
			b.SolutionCount = len(slns)
			fmt.Printf("solved with SAT. solution count: %d\n", len(slns))
		}

		sln1 := slns[0]
		for _, setvar := range sln1.SetVars {
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
		//b.SolveNakedSingle()
	}

	/*if !b.isSolved() {
		b.changed = true
		for b.changed {
			b.changed = false
			b.SolveNakedSingle()
			b.SolveHiddenSingle()
		}
	}*/

	err = b.Validate()
	if err != nil {
		return err
	}

	/*var vars []int

	for _, item := range satSolver.SetVars {
		vars = append(vars, int(item.VarNum))
	}

	sort.Ints(vars)

	for _, v := range vars {
		for _, item := range satSolver.SetVars {
			if v == int(item.VarNum) {
				fmt.Printf("%d %t\n", item.VarNum, item.Value)
			}
		}
	}*/

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
