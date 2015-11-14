package sat

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
	sat, err := NewSAT(input)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Printf("clauses:  %#v\n", sat.Clauses)
	sln := sat.Solve()
	if sln == nil {
		fmt.Printf("no solution\n")
	} else {
		fmt.Printf("set vars: %v\n", sln.SetVars)
		fmt.Printf("clauses:  %#v\n", sln.Clauses)
	}
}

func TestHasClause(t *testing.T) {
	clause := []int{0, 2, 6, 8, 12}
	for idx, val := range clause {
		actual := indexOfValue(clause, val)
		if actual != idx {
			t.Fatalf("%d not found in clause %v. idx:%d", val, clause, idx)
		}
	}

	clause = []int{0, 2, 6, 8, 12, 14}
	for idx, val := range clause {
		actual := indexOfValue(clause, val)
		if actual != idx {
			t.Fatalf("%d not found in clause %v. idx:%d", val, clause, idx)
		}
	}

	clause = []int{-937, -737}
	for idx, val := range clause {
		actual := indexOfValue(clause, val)
		if actual != idx {
			t.Fatalf("%d not found in clause %v. idx:%d", val, clause, idx)
		}
	}
}

func TestUnitPropogation(t *testing.T) {
	clause := []int{1, 2, 6, 8, 12}

	expected := []int{1, 2, 6, 8}
	actual := up(clause, 12, false)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v actual: %v", expected, actual)
	}

	expected = nil
	actual = up(clause, 12, true)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v actual: %v", expected, actual)
	}

	clause = []int{-1, 2, 6, 8, 12}

	expected = []int{2, 6, 8, 12}
	//fmt.Printf("clause before1: %v\n", clause)
	actual = up(clause, 1, true)
	//fmt.Printf("clause before2: %v\n", clause)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v actual: %v", expected, actual)
	}

	expected = nil
	//fmt.Printf("clause before3: %v\n", clause)
	actual = up(clause, 1, false)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v actual: %v", expected, actual)
	}
}
