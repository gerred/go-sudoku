package sat

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type sat struct {
	vars    []int
	SetVars map[int]bool
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
	parts, err := getIntArray(strParts)
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
		SetVars: make(map[int]bool),
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
		strParts = strings.SplitN(line, " ", -1)
		parts, err := getIntArray(strParts)
		if err != nil {
			return nil, fmt.Errorf("error parsing line: %q", line)
		}
		if parts[len(parts)-1] != 0 {
			return nil, fmt.Errorf("error parsing line, must end in \"0\": %q", line)
		}
		s.Clauses[i] = parts[:len(parts)-1]
		for j := 0; j < len(s.Clauses[i]); j++ {
			v := abs(s.Clauses[i][j])
			s.vars = append(s.vars, v)
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

func getIntArray(values []string) ([]int, error) {
	list := make([]int, len(values))
	for i := 0; i < len(values); i++ {
		v, err := strconv.Atoi(values[i])
		if err != nil {
			return nil, fmt.Errorf("idx: %d unexpected format: %q", i, values[i])
		}
		list[i] = v
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
	//for k, _ := range s.vars {
	//	fmt.Printf("%d ", k)
	//}
	//fmt.Println()
	//fmt.Printf("%#2v\n", s.clauses)

	// let's try some clauses
	s2 := set(s, s.vars[0], true, 0)
	if s2 != nil {
		return s2
	}
	s2 = set(s, s.vars[0], false, 0)
	if s2 != nil {
		return s2
	}
	return nil
}

func set(s1 *sat, v int, value bool, depth int) *sat {
	s2 := &sat{
		SetVars: make(map[int]bool),
	}
	for _, k := range s1.vars {
		if k != v {
			s2.vars = append(s2.vars, k)
		}
	}
	for k, _ := range s1.SetVars {
		s2.SetVars[k] = s1.SetVars[k]
	}
	s2.SetVars[v] = value

	//prefix := strings.Repeat("-", depth+1)
	//fmt.Printf("%s %v\n", prefix, s2.setVars)

	signedV := v
	if !value {
		signedV *= -1
	}
	s2.Clauses = append(s2.Clauses, []int{signedV})
	for _, clause := range s1.Clauses {
		newClause := up(clause, v, value)
		if newClause != nil {
			if len(newClause) == 0 {
				//fmt.Printf("%s failed clause: idx:%d %v\n", prefix, idx, clause)
				return nil
			}
			s2.Clauses = append(s2.Clauses, newClause)
		}
	}

	if len(s2.vars) == 0 {
		return s2
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
	newClause := make([]int, 0)
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
	return newClause
}
