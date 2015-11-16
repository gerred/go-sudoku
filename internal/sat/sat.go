package sat

import (
	"bufio"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type SetVar struct {
	VarNum uint64
	Value  bool
}

type sat struct {
	vars    []uint64
	SetVars []SetVar
	Clauses [][2]uint64
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
		Clauses: make([][2]uint64, clauseCount),
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
		s.Clauses[i] = intArrayToBin(parts)

		for j := 0; j < len(parts); j++ {
			v := uint64(abs(parts[j]))
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

func intArrayToBin(list []int) [2]uint64 {
	var bin [2]uint64
	bin[0] = uint64(len(list)) << 59
	j := 0
	shift := uint(44)
	for i := 0; i < len(list); i++ {
		val := abs(list[i])
		if list[i] < 0 {
			val |= 0x400 // add sign
		}
		bin[j] |= uint64(val << shift)

		if shift == 0 {
			shift = 44
			j++
		} else {
			shift -= 11
		}
	}
	return bin
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
		//list = cut(list, len(list)-1)
		list = list[:len(list)-1]
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

const lenMask = 0xF800000000000000
const len1 = 0x800000000000000

func (s *sat) Solve() *sat {
	//fmt.Printf("len(s.vars) = %d\n", len(s.vars))
	//depth := 0
	// find a clause with a single variable
	for _, clause := range s.Clauses {
		if clause[0]&lenMask == len1 {
			val := clause[0] >> 44
			on := (val&0x400 == 0)
			val = val & 0x3FF

			//fmt.Printf("quick find [%d]: %d %t\n", depth, val, on)
			return set(s, val, on)
		}
	}

	// let's try the first variable
	s2 := set(s, s.vars[0], true)
	if s2 != nil {
		return s2
	}
	s2 = set(s, s.vars[0], false)
	if s2 != nil {
		return s2
	}
	return nil
}

func set(s1 *sat, v uint64, isOn bool) *sat {
	s2 := &sat{}
	for _, k := range s1.vars {
		if k != v {
			s2.vars = append(s2.vars, k)
		}
	}
	s2.SetVars = append(s2.SetVars, s1.SetVars...)
	s2.SetVars = append(s2.SetVars, SetVar{VarNum: v, Value: isOn})

	for _, clause := range s1.Clauses {
		newClause := up(clause, v, isOn)
		if newClause != nil {
			if newClause[0] == 0 {
				//s2.checkKnownAnswers(v, value)
				return nil
			}
			s2.Clauses = append(s2.Clauses, *newClause)
		}
	}

	if len(s2.vars) == 0 {
		return s2
	}

	// find a clause with a single variable
	for _, clause := range s2.Clauses {
		if clause[0]&lenMask == len1 {
			val := clause[0] >> 44
			on := (val&0x400 == 0)
			val = val & 0x3FF

			found := false
			for _, x := range s2.vars {
				if x == val {
					found = true
					break
				}
			}

			if !found {
				//fmt.Printf("not found %d %t\n", val, on)
				//s2.checkKnownAnswers(v, value)
				return nil
			}

			//fmt.Printf("quick find [%d]: %d %t\n", depth+1, val, on)
			return set(s2, val, on)
		}
	}

	s3 := set(s2, s2.vars[0], true)
	if s3 != nil {
		return s3
	}
	s3 = set(s2, s2.vars[0], false)
	if s3 != nil {
		return s3
	}

	//s2.checkKnownAnswers(v, value)
	return nil
}

/*func (s *sat) checkKnownAnswers(v int, val bool) {
	if s.KnownAnswers != nil {
		//fmt.Printf("len(knownAnswers)=%d\n", len(s.KnownAnswers))
		for _, ka := range s.KnownAnswers {
			if ka.VarNum == v && ka.Value == val {
				fmt.Printf("*** WTF? %d %t\n", v, val)
			}
		}
	}
}*/

func up(clause [2]uint64, v uint64, isOn bool) *[2]uint64 {
	//	fmt.Printf("-------\n")
	//	fmt.Printf("v:%d clause:%v\n", v, clause)
	var idx int
	if idx = indexOfValue(clause, v); idx != -1 {
		if isOn {
			return nil
		}
		//fmt.Printf("pos-idx:%d\n", idx)
		return cut(clause, idx)
	} else if idx = indexOfValue(clause, v|0x400); idx != -1 {
		if !isOn {
			return nil
		}
		//fmt.Printf("neg-idx:%d\n", idx)
		return cut(clause, idx)
	} else {
		//fmt.Printf("nada: %v\n", clause)
		return &clause
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

func cut(clause [2]uint64, idx int) *[2]uint64 {
	var newClause [2]uint64
	length := int((clause[0] & lenMask) >> 59)
	if idx == 0 && length == 1 {
		return &newClause
	}

	newClause[0] = uint64(length-1) << 59

	shift := uint(44)
	shift2 := uint(44)
	j, k := 0, 0
	cur := clause[0]
	for i := 0; i < length; i++ {
		if i != idx {
			curval := (cur >> shift) & 0x7FF
			newClause[k] |= curval << shift2

			if shift2 == 0 {
				shift2 = 44
				k++
			} else {
				shift2 -= 11
			}
		}

		if shift == 0 {
			shift = 44
			j++
			cur = clause[j]
		} else {
			shift -= 11
		}
	}
	return &newClause

	/*l := len(clause)
	if idx == 0 {
		if l == 1 {
			return []int{}
		}
		return clause[1:]
	} else if idx == l-1 {
		return clause[:idx]
	} else {
		var newClause []int
		newClause = append(newClause, clause[:idx]...)
		newClause = append(newClause, clause[idx+1:]...)
		return newClause
	}*/
}

func indexOfValue(clause [2]uint64, val uint64) int {
	//fmt.Printf("-- %011b %011b, %011b\n", clause[0], clause[1], val)
	length := int((clause[0] & lenMask) >> 59)
	if length == 0 {
		return -1
	}
	shift := uint(44)
	j := 0
	cur := clause[0]
	for i := 0; i < length; i++ {
		curval := (cur >> shift) & 0x7FF
		//fmt.Printf("curval[%d]: %011b\n", i, curval)
		//if curval > val {
		//	return -1
		//}
		if curval == val {
			return i
		}

		if shift == 0 {
			shift = 44
			j++
			cur = clause[j]
		} else {
			shift -= 11
		}
	}
	return -1

	// slice iterate
	/*for i := 0; i < len(clause); i++ {
		if clause[i] == val {
			return i
		}
	}
	return -1*/

	// binary search
	/*max := len(clause)
	if max == 0 {
		return -1
	}

	// do binary search
	i := 0
	min := 0
	step := max - 1
	for {
		if clause[i] == val {
			return i
		}

		if clause[i] > val {
			max = i
			i -= step
			if i < min {
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
				return -1
			}

			step >>= 1
			if step == 0 && i < max {
				step = 1
			}
		}

		if step == 0 {
			return -1
		}
	}*/
}
