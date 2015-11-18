package sat

import (
	"bufio"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const lenMask = 0xF800000000000000
const len1 = 0x800000000000000
const len2 = 0x1000000000000000

var notFound *[2]uint64 = &[2]uint64{}

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
		list = list[:len(list)-1]
	}

	if sortValues {
		sort.Ints(list)
	}

	return list, nil
}

func abs(x int) int {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0 // return correctly abs(-0)
	}
	return x
}

func (s *sat) getNextVar() (*uint64, *bool) {
	// find a clause with a single variable
	for _, clause := range s.Clauses {
		length := clause[0] & lenMask
		if length == len1 {
			val := clause[0] >> 44
			on := (val&0x400 == 0)
			val = val & 0x3FF

			// TODO: is this necessary?
			found := false
			for _, x := range s.vars {
				if x == val {
					found = true
					break
				}
			}

			if !found {
				return nil, nil
			}

			return &val, &on
		}
	}

	for _, clause := range s.Clauses {
		length := clause[0] & lenMask
		if length == len2 {
			val := clause[0] >> 44
			on := (val&0x400 == 0)
			if on {
				continue
			}
			val = val & 0x3FF

			// TODO: is this necessary?
			found := false
			for _, x := range s.vars {
				if x == val {
					found = true
					break
				}
			}

			if !found {
				return nil, nil
			}

			return &val, nil
		}
	}

	// let's try the first variable
	return &s.vars[0], nil
}

func (s *sat) Solve() *sat {
	val, on := s.getNextVar()
	var s2 *sat
	if on != nil {
		s2 = set(s, *val, *on)
	} else {
		s2 = set(s, *val, true)
		if s2 == nil {
			s2 = set(s, *val, false)
		}
	}
	return s2
}

func set(s1 *sat, v uint64, isOn bool) *sat {
	s2 := &sat{}
	for _, k := range s1.vars {
		if k != v {
			s2.vars = append(s2.vars, k)
		}
	}

	for _, clause := range s1.Clauses {
		newClause := up(&clause, v, isOn)
		if newClause != nil {
			if newClause != notFound {
				s2.Clauses = append(s2.Clauses, *newClause)
			}
		} else {
			return nil
		}
	}

	if len(s2.vars) == 0 {
		s2.SetVars = append(s2.SetVars, SetVar{VarNum: v, Value: isOn})
		return s2
	}

	val, on := s2.getNextVar()

	if on != nil {
		s3 := set(s2, *val, *on)
		if s3 != nil {
			s3.SetVars = append(s3.SetVars, SetVar{VarNum: v, Value: isOn})
			return s3
		}
	} else {
		s3 := set(s2, *val, true)
		if s3 != nil {
			s3.SetVars = append(s3.SetVars, SetVar{VarNum: v, Value: isOn})
			return s3
		}
		s3 = set(s2, *val, false)
		if s3 != nil {
			s3.SetVars = append(s3.SetVars, SetVar{VarNum: v, Value: isOn})
			return s3
		}
	}

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

func up(clause *[2]uint64, v uint64, isOn bool) *[2]uint64 {
	var idx int
	if idx = indexOfValue(clause, v); idx != -1 {
		if isOn {
			return notFound
		}
		return cut(clause, idx)
	} else if idx = indexOfValue(clause, v|0x400); idx != -1 {
		if !isOn {
			return notFound
		}
		return cut(clause, idx)
	} else {
		return clause
	}
}

func cut(clause *[2]uint64, idx int) *[2]uint64 {
	if clause[0]&lenMask == len1 {
		return nil
	}

	var newClause [2]uint64

	length := int((clause[0] & lenMask) >> 59)
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

}

func indexOfValue(clause *[2]uint64, val uint64) int {
	length := int((clause[0] & lenMask) >> 59)
	if length == 0 {
		return -1
	}

	shift := uint(44)
	cur := clause[0]
	for i := 0; i < length; i++ {
		curval := (cur >> shift) & 0x7FF
		if curval == val {
			return i
		}

		if shift == 0 {
			shift = 44
			cur = clause[1]
		} else {
			shift -= 11
		}
	}
	return -1

	// binary search
	/*i := 0
	min := 0
	max := length
	step := length - 1
	var findVal int
	findVal = int(val)
	if val&0x400 == 0x400 {
		findVal = (findVal & 0x3FF) * -1 // flip sign
	}
	//fmt.Printf("index-of: %b %b val:%d\n", clause[0], clause[1], findVal)
	for {
		// 44, 33, 22, 11 0
		var curval int
		if i >= 5 {
			shift := uint(44 - ((i - 5) * 11))
			curval = int((clause[1] >> shift) & 0x7FF)
		} else {
			shift := uint(44 - (i * 11))
			curval = int((clause[0] >> shift) & 0x7FF)
		}

		if curval&0x400 == 0x400 {
			curval = (curval & 0x3FF) * -1
			//fmt.Printf("-- curval: %d\n", curval)
		}

		if curval == findVal {
			return i
		}

		if curval > findVal {
			max = i
			i -= step
			if i < min {
				//fmt.Printf("i < min\n")
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
				//fmt.Printf("i > max\n")
				return -1
			}

			step >>= 1
			if step == 0 && i < max {
				step = 1
			}
		}

		//fmt.Printf("- new step: %d new i: %d\n", step, i)
		if step == 0 {
			return -1
		}
	}
	return -1*/
}
