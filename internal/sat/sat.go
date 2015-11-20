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

var satisfied *[2]uint64 = &[2]uint64{0xFF, 0xFF}

type SetVar struct {
	VarNum uint64
	Value  bool
}

type sat struct {
	SetVars               []SetVar
	Clauses               [][2]uint64
	FindMultipleSolutions bool
}

func NewSAT(input string, findMultipleSolutions bool) (*sat, error) {
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
		Clauses:               make([][2]uint64, clauseCount),
		FindMultipleSolutions: findMultipleSolutions,
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

/*func (s *sat) getRemainingVars() []int {
	vars := make(map[int]interface{})

	for _, clause := range s.Clauses {
		length := int(clause[0]&lenMask) >> 59
		if length == 0 {
			continue
		}
		shift := uint(44)
		cur := clause[0]
		for i := 0; i < length; i++ {
			curval := (cur >> shift) & 0x7FF
			neg := (curval & 0x400) == 0x400
			curval &= 0x3FF

			var val int
			val = int(curval)
			if neg {
				val = -val
			}
			vars[val] = struct{}{}

			if shift == 0 {
				shift = 44
				cur = clause[1]
			} else {
				shift -= 11
			}
		}
	}

	var list []int
	for k, _ := range vars {
		list = append(list, k)
	}
	sort.Ints(list)
	return list
}*/

func (s *sat) getAllSingleVarClauses() []SetVar {
	/*check := make(map[uint64]struct{})
	for _, clause := range s.Clauses {
		c := clause[0]
		length := c & lenMask
		if length == len1 {
			if _, ok := check[c]; !ok {
				continue
			}
			check[c] = struct{}{}
		}
	}*/

	//list := make([]SetVar, len(check))
	check := make(map[uint64]struct{})
	var list []SetVar
	var val uint64
	var on bool

	for _, clause := range s.Clauses {
		c := clause[0]
		length := c & lenMask
		if length == len1 {
			val = c >> 44
			on = (val&0x400 == 0)
			val = val & 0x3FF

			if _, ok := check[c]; !ok {
				/*delete(check, c)
				list = append(list, SetVar{VarNum: val, Value: on})*/
				continue
			}
			check[c] = struct{}{}

			list = append(list, SetVar{VarNum: val, Value: on})
		}
	}
	return list
}

func (s *sat) getNextSingleVar() (*uint64, *bool) {
	var val uint64

	// find a clause with a single variable
	for _, clause := range s.Clauses {
		length := clause[0] & lenMask
		if length == len1 {
			val = clause[0] >> 44
			on := (val&0x400 == 0)
			val = val & 0x3FF

			return &val, &on
		}
	}

	return nil, nil
}

func (s *sat) getNextVar() (*uint64, *bool) {
	singleVal, on := s.getNextSingleVar()
	if singleVal != nil {
		return singleVal, on
	}

	var val uint64

	for _, clause := range s.Clauses {
		length := clause[0] & lenMask
		if length == len2 {
			val = (clause[0] >> 44) & 0x3FF
			return &val, nil
		}
	}

	// let's try the first variable from the first clause
	val = (s.Clauses[0][0] >> 44) & 0x3FF
	return &val, nil
}

func (s *sat) Solve() []*sat {
	val, on := s.getNextVar()
	var s2, s3 []*sat
	if on != nil {
		s2 = set(s, *val, *on)
	} else {
		s2 = set(s, *val, false)
		if s2 == nil || s.FindMultipleSolutions {
			s3 = set(s, *val, true)
		}
	}

	var final []*sat
	if s2 != nil {
		final = append(final, s2...)
	}
	if s3 != nil {
		final = append(final, s3...)
	}

	if len(final) == 0 {
		return nil
	}

	return final
}

/*func (s2 *sat) solveSingleVarClauses() (*sat, []SetVar) {
	var bigList []SetVar

	val, on := s2.getNextSingleVar()
	for val != nil {
		var clauses [][2]uint64
		for _, clause := range s2.Clauses {
			newClause := up(&clause, *val, *on)
			if newClause != nil {
				if newClause != satisfied {
					clauses = append(clauses, *newClause)
				}
			} else {
				return nil, nil
			}
		}

		s2.Clauses = clauses
		bigList = append(bigList, SetVar{VarNum: *val, Value: *on})
		if len(clauses) == 0 {
			break
		}
		val, on = s2.getNextSingleVar()
	}

	return s2, bigList
}*/

func (s2 *sat) solveSingleVarClauses2() (*sat, []SetVar) {
	var bigList []SetVar

	list := s2.getAllSingleVarClauses()
	for len(list) != 0 {
		for _, item := range list {
			var clauses [][2]uint64
			for _, clause := range s2.Clauses {
				newClause := up(&clause, item.VarNum, item.Value)
				if newClause != nil {
					if newClause != satisfied {
						clauses = append(clauses, *newClause)
					}
				} else {
					return nil, nil
				}
			}

			s2.Clauses = clauses
			if len(s2.Clauses) == 0 {
				break
			}
		}
		bigList = append(bigList, list...)
		if len(s2.Clauses) == 0 {
			break
		}
		list = s2.getAllSingleVarClauses()
	}

	return s2, bigList
}

func set(s1 *sat, v uint64, isOn bool) []*sat {
	s2 := &sat{FindMultipleSolutions: s1.FindMultipleSolutions}

	for _, clause := range s1.Clauses {
		newClause := up(&clause, v, isOn)
		if newClause != nil {
			if newClause != satisfied {
				s2.Clauses = append(s2.Clauses, *newClause)
			}
		} else {
			return nil
		}
	}

	if len(s2.Clauses) == 0 {
		s2.SetVars = append(s2.SetVars, SetVar{VarNum: v, Value: isOn})
		return []*sat{s2}
	}

	var bigList []SetVar
	if s2.FindMultipleSolutions {
		s2, bigList = s2.solveSingleVarClauses2()
		if s2 == nil {
			return nil
		}

		if len(s2.Clauses) == 0 {
			s2.SetVars = append(s2.SetVars, bigList...)
			s2.SetVars = append(s2.SetVars, SetVar{VarNum: v, Value: isOn})
			return []*sat{s2}
		}
	}

	val, on := s2.getNextVar()

	var s3, s4 []*sat
	if on != nil {
		s3 = set(s2, *val, *on)
	} else {
		s3 = set(s2, *val, false)
		if s3 == nil || s2.FindMultipleSolutions {
			s4 = set(s2, *val, true)
		}
	}

	var final []*sat
	if s3 != nil {
		final = append(final, s3...)
	}
	if s4 != nil {
		final = append(final, s4...)
	}

	if len(final) == 0 {
		return nil
	}

	for _, item := range final {
		if s1.FindMultipleSolutions {
			item.SetVars = append(item.SetVars, bigList...)
		}
		item.SetVars = append(item.SetVars, SetVar{VarNum: v, Value: isOn})
	}

	return final
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
			return satisfied
		}
		return cut(clause, idx)
	} else if idx = indexOfValue(clause, v|0x400); idx != -1 {
		if !isOn {
			return satisfied
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
