package sat

import (
	"bufio"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type setvar struct {
	VarNum int
	Value  bool
}

type sat struct {
	vars []int
	//SetVars map[int]bool
	SetVars []setvar
	Clauses [][]int
}

func NewSAT(input string) (*sat, error) {
	// load CNF input

	sr := strings.NewReader(input)
	r := bufio.NewReader(sr)

	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	for strings.HasPrefix(line, "c") {
		// skip comment
		if line, err = r.ReadString('\n'); err != nil {
			return nil, err
		}
	}
	// p cnf # # // variables, clauses
	if !strings.HasPrefix(line, "p cnf ") {
		return nil, fmt.Errorf("expected first non-comment line in format\"p cnf # #\": %q", line)
	}
	// get # of variables, # of clauses from CNF header
	line = strings.Trim(line[len("p cnf"):], " \r\n\t")
	strParts := strings.SplitN(line, " ", -1)
	parts, err := getIntArray(strParts, false, false)
	if err != nil {
		return nil, fmt.Errorf("expected first non-comment line in format \"p cnf # #\": %q %q", line, err)
	}
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected first non-comment line in format \"p cnf # #\": %q", line)
	}
	variableCount := parts[0]
	clauseCount := parts[1]
	if variableCount < 0 || clauseCount < 0 {
		return nil, fmt.Errorf("variable and clause count must be non-negative: %q", line)
	}

	// TODO: validate variable, clause count
	s := &sat{
		Clauses: make([][]int, clauseCount),
		//SetVars: make([]setvar, variableCount),
		//SetVars: make(map[int]bool),
	}

	// get clauses
	if line, err = r.ReadString('\n'); err != nil {
		if err.Error() == "EOF" {
			// TODO: if clauses = 0, this is okay.
		}
		return nil, err
	}

	var lastLine bool
	for i := 0; line != ""; i++ {
		line = strings.Trim(line, " \r\n\t")
		if strings.HasPrefix(line, "c ") {
			// skip comment
			i--
			continue
		}
		strParts = strings.SplitN(line, " ", -1)
		parts, err := getIntArray(strParts, true, true)
		if err != nil {
			return nil, fmt.Errorf("error parsing line: %q", line)
		}
		s.Clauses[i] = parts

		for j := 0; j < len(s.Clauses[i]); j++ {
			v := abs(s.Clauses[i][j])
			found := false
			for _, x := range s.vars {
				if x == v {
					found = true
					break
				}
			}
			if !found {
				s.vars = append(s.vars, v)
			}
		}
		if lastLine == true {
			break
		}

		if line, err = r.ReadString('\n'); err != nil {
			if err.Error() == "EOF" {
				if line != "" {
					lastLine = true
					continue
				}
				break
			}
			return nil, err
		}
	}

	return s, nil
}

func getIntArray(values []string, sortValues bool, trimEnd bool) ([]int, error) {
	list := make([]int, len(values))
	for i := 0; i < len(values); i++ {
		v, err := strconv.Atoi(values[i])
		if err != nil {
			return nil, fmt.Errorf("idx: %d unexpected format: %q", i, values[i])
		}
		list[i] = v
	}

	if trimEnd {
		if list[len(list)-1] != 0 {
			return nil, errors.New("error parsing line, must end in \"0\"")
		}
		list = cut(list, len(list)-1)
	}

	if sortValues {
		sort.Ints(list)
	}

	return list, nil
}

func abs(v int) int {
	if v >= 0 {
		return v
	}
	return v * -1
}

func (s *sat) Solve() *sat {
	fmt.Printf("len(s.vars) = %d\n", len(s.vars))
	depth := 0

	// find a clause with a single variable
	for _, clause := range s.Clauses {
		if len(clause) == 1 {
			val := clause[0]
			on := (val > 0)
			val = abs(val)

			//fmt.Printf("quick find [%d]: %d %t\n", depth, val, on)
			return set(s, val, on, depth)
		}
	}

	// let's try the first variable
	s2 := set(s, s.vars[0], true, depth)
	if s2 != nil {
		return s2
	}
	s2 = set(s, s.vars[0], false, depth)
	if s2 != nil {
		return s2
	}
	return nil
}

func set(s1 *sat, v int, value bool, depth int) *sat {
	s2 := &sat{}
	for _, k := range s1.vars {
		if k != v {
			s2.vars = append(s2.vars, k)
		}
	}
	s2.SetVars = append(s2.SetVars, s1.SetVars...)
	s2.SetVars = append(s2.SetVars, setvar{VarNum: v, Value: value})
	//s2.vars[depth] = k
	//s2.SetVars[depth] = setvar{VarNum: v, Value: value}

	signedV := v
	if !value {
		signedV *= -1
	}
	//s2.Clauses = append(s2.Clauses, []int{signedV})
	for _, clause := range s1.Clauses {
		newClause := up(clause, v, value)
		if newClause != nil {
			if len(newClause) == 0 {
				return nil
			}
			s2.Clauses = append(s2.Clauses, newClause)
		}
	}

	if len(s2.vars) == 0 {
		return s2
	}

	// find a clause with a single variable
	for _, clause := range s2.Clauses {
		if len(clause) == 1 {
			val := clause[0]
			on := (val > 0)
			val = abs(val)

			found := false
			for _, x := range s2.vars {
				if x == val {
					found = true
					break
				}
			}

			if !found {
				fmt.Printf("not found %d %t\n", val, on)
				return nil
			}

			//fmt.Printf("quick find [%d]: %d %t\n", depth+1, val, on)
			return set(s2, val, on, depth+1)
		}
	}

	s3 := set(s2, s2.vars[0], true, depth+1)
	if s3 != nil {
		return s3
	}
	s3 = set(s2, s2.vars[0], false, depth+1)
	if s3 != nil {
		return s3
	}

	return nil
}

func up(clause []int, v int, val bool) []int {
	//	fmt.Printf("-------\n")
	//	fmt.Printf("v:%d clause:%v\n", v, clause)
	var idx int
	if idx = indexOfValue(clause, v); idx != -1 {
		if val {
			return nil
		}
		//fmt.Printf("pos-idx:%d\n", idx)
		return cut(clause, idx)
	} else if idx = indexOfValue(clause, -v); idx != -1 {
		if !val {
			return nil
		}
		//fmt.Printf("neg-idx:%d\n", idx)
		return cut(clause, idx)
	} else {
		//fmt.Printf("nada: %v\n", clause)
		return clause
	}

	/*newClause := make([]int, 0)
	for _, k := range clause {
		if k == v {
			// true statement
			if val {
				return nil
			}
		} else if k == -v {
			// true statement
			if !val {
				return nil
			}
		} else {
			newClause = append(newClause, k)
		}
	}
	return newClause*/
}

func cut(clause []int, idx int) []int {
	/*//newClause := make([]int, len(clause)-1)
	newClause := make([]int, 0)
	//i := 0
	for j := 0; j < len(clause); j++ {
		if j != idx {
			newClause = append(newClause, clause[j])
			//newClause[i] = clause[j]
			//i++
		}
	}
	return newClause*/

	l := len(clause)
	if idx == 0 {
		if l == 1 {
			return []int{}
		}
		return clause[1:]
	} else if idx == l-1 {
		return clause[:idx]
	} else {
		return append(clause[:idx], clause[idx+1:]...)
		/*newClause := make([]int, len(clause)-1)
		i := 0
		for j := 0; j < len(clause); j++ {
			if j != idx {
				//newClause = append(newClause, clause[j])
				newClause[i] = clause[j]
				i++
			}
		}
		return newClause*/
	}
}

func indexOfValue(clause []int, val int) int {
	/*slowIdx := -1
	for i := 0; i < len(clause); i++ {
		if clause[i] == val {
			slowIdx = i
			break
		}
	}*/
	//slowIdx = -1 // return -1

	max := len(clause)
	if max == 0 {
		/*if slowIdx != -1 {
			fmt.Printf("%v val:%d (%d != %d)\n", clause, val, slowIdx, -1)
		}*/
		return -1
	}

	// do binary search
	i := 0
	min := 0
	step := max - 1
	for {
		if clause[i] == val {
			/*if slowIdx != i {
				fmt.Printf("%v val:%d (%d != %d)\n", clause, val, slowIdx, i)
			}*/
			return i
		}

		if clause[i] > val {
			max = i
			i -= step
			if i < min {
				/*if slowIdx != -1 {
					fmt.Printf("%v val:%d (%d != %d)\n", clause, val, slowIdx, -1)
				}*/
				return -1
			}

			step >>= 1
			if step == 0 && i > min {
				step = 1
			}
		} else {
			min = i
			i += step
			if i > max {
				/*if slowIdx != -1 {
					fmt.Printf("%v val:%d (%d != %d)\n", clause, val, slowIdx, -1)
				}*/
				return -1
			}

			step >>= 1
			if step == 0 && i < max {
				step = 1
			}
		}

		if step == 0 {
			/*if slowIdx != -1 {
				fmt.Printf("%v val:%d (%d != %d) step == 0\n", clause, val, slowIdx, -1)
			}*/
			return -1
		}
	}
}
