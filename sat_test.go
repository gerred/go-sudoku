package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSolve(t *testing.T) {
	inputs := []string{
		// (a ∨ ¬b) ∧ (a ∨ b) should return a: True, b: anything
		`p cnf 2 2
		 1 -2 0
		 1 2 0`,
		// (a ∧ b) ∧ (a ∧ ¬b) should return: null
		`p cnf 2 4
		 1 0
		 2 0
		 1 0
		 -2 0`,
		// (a ∧ b) ∧ (¬b ∨ c) should return: a: True, b: True, c: True.
		`p cnf 3 3
		 1 0
		 2 0
		 -1 3 0`,
		// (x ∨ x ∨ y) ∧ (¬x ∨ ¬y ∨ ¬y) ∧ (¬x ∨ y ∨ y) x: True, y: False
		`p cnf 2 3
		 1 1 2 0
		 -1 -2 -2 0
		 -1 2 2 0`,
	}

	for _, input := range inputs {
		testInput(t, input)
	}
}

func testInput(t *testing.T, input string) {
	fmt.Println()
	if len(input) < 1000 {
		fmt.Println(input)
	}
	sat, err := NewSAT(input, false, 0)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Printf("clauses:  %#v\n", sat.Clauses)
	sln := sat.Solve()
	if sln == nil || len(sln) == 0 {
		fmt.Printf("no solution\n")
	} else {
		for _, item := range sln {
			fmt.Printf("set vars: %v\n", item.SetVars)
			fmt.Printf("clauses:  %#v\n", item.Clauses)
		}
	}
}

func TestHasClause(t *testing.T) {
	clauseIntArray := []int{0, 2, 6, 8, 12}
	clause := intArrayToBin(clauseIntArray)
	//fmt.Printf("%b %b\n", clause[0], clause[1])
	for idx, val := range clauseIntArray {
		//fmt.Printf("%d\n", idx)
		actual := indexOfValue(&clause, uint64(val))
		if actual != idx {
			t.Fatalf("%d not found in clause %b %b. idx:%d", val, clause[0], clause[1], idx)
		}
	}

	clauseIntArray = []int{0, 2, 6, 8, 12, 14}
	clause = intArrayToBin(clauseIntArray)
	for idx, val := range clauseIntArray {
		actual := indexOfValue(&clause, uint64(val))
		if actual != idx {
			// 110 0000 = 6
			// 00000000000 = 0
			// 00000000010 = 2
			// 00000000110 = 6
			// 00000001000 = 8
			// 00000001100 = 12

			// 111
			// 00000000000
			// 00000000000
			// 00000000000
			// 00000000000
			//         111000000000000000000000000000000000000000000000
			// 1234567890x1234567890x1234567890x12345678901x12345678901
			t.Fatalf("%d not found in clause %b %b. idx:%d", val, clause[0], clause[1], idx)
		}
	}

	clauseIntArray = []int{-937, -737}
	clause = intArrayToBin(clauseIntArray)
	for idx, val := range clauseIntArray {
		uintVal := uint64(abs(val))
		if val < 0 {
			uintVal |= 0x400
		}

		//  100000
		// 11110101001
		// 11011100001
		// 11011100001
		// 000000000000000000000000000000000
		// 111xxxx1234567890x1234567890x1234567890x12345678901x12345678901

		actual := indexOfValue(&clause, uintVal)
		if actual != idx {
			t.Fatalf("%d (%b) not found in clause %b %b. idx:%d", val, uintVal, clause[0], clause[1], idx)
		}
	}
}

func TestUnitPropogation(t *testing.T) {
	clause := intArrayToBin([]int{1, 2, 6, 8, 12})

	expected_tmp := intArrayToBin([]int{1, 2, 6, 8})
	expected := &expected_tmp
	actual := up(&clause, 12, false)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v actual: %v", expected, actual)
	}

	expected = satisfied
	actual = up(&clause, 12, true)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v actual: %v", expected, actual)
	}

	clause = intArrayToBin([]int{-1, 2, 6, 8, 12})

	expected_tmp = intArrayToBin([]int{2, 6, 8, 12})
	expected = &expected_tmp
	//fmt.Printf("clause before1: %v\n", clause)
	actual = up(&clause, 1, true)
	//fmt.Printf("clause before2: %v\n", clause)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v actual: %v", expected, actual)
	}

	expected = satisfied
	//fmt.Printf("clause before3: %v\n", clause)
	actual = up(&clause, 1, false)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v actual: %v", expected, actual)
	}
}

func TestIntArrayToBin(t *testing.T) {
	// TODO
}
