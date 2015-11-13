package sat

import (
	"fmt"
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
