package main

import (
	"reflect"
	"testing"
)

func TestSetsAreEqual(t *testing.T) {
	// arrange
	inputs := [][][]int{
		[][]int{[]int{1, 2, 3}, []int{4, 5, 6}},
		[][]int{[]int{1, 2, 3}, []int{3, 2, 1}},
		[][]int{[]int{1, 2, 3}, []int{2, 2, 4}},
		[][]int{[]int{}, []int{}},
	}

	expecteds := []bool{
		false,
		true,
		false,
		true,
	}

	// act
	var actuals []bool
	for _, input := range inputs {
		actual := setsAreEqual(input[0], input[1])
		actuals = append(actuals, actual)
	}

	// assert
	for i, expected := range expecteds {
		actual := actuals[i]
		if expected != actual {
			t.Fatalf("inputs: %v expected: %v actual: %v", inputs[i], expected, actual)
		}
	}
}

func TestUnion(t *testing.T) {
	// arrange
	inputs := [][][]int{
		[][]int{[]int{1, 2, 3}, []int{4, 5, 6}},
		[][]int{[]int{1, 2, 3}, []int{3, 2, 1}},
		[][]int{[]int{1, 2, 3}, []int{2, 2, 4}},
		[][]int{[]int{1, 2, 3}, []int{3, 4, 5}},
	}

	expecteds := [][]int{
		[]int{1, 2, 3, 4, 5, 6},
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		[]int{1, 2, 3, 4, 5},
	}

	// act
	var actuals [][]int
	for _, input := range inputs {
		actual := union(input[0], input[1])
		actuals = append(actuals, actual)
	}

	// assert
	for i, expected := range expecteds {
		actual := actuals[i]
		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("inputs: %v expected: %v actual: %v", inputs[i], expected, actual)
		}
	}
}

func TestIntersect(t *testing.T) {
	// arrange
	inputs := [][][]int{
		[][]int{[]int{1, 2, 3}, []int{4, 5, 6}},
		[][]int{[]int{1, 2, 3}, []int{3, 2, 1}},
		[][]int{[]int{1, 2, 3}, []int{2, 2, 4}},
		[][]int{[]int{1, 2, 3}, []int{3, 4, 5}},
	}

	expecteds := [][]int{
		[]int{},
		[]int{1, 2, 3},
		[]int{2},
		[]int{3},
	}

	// act
	var actuals [][]int
	for _, input := range inputs {
		actual := intersect(input[0], input[1])
		actuals = append(actuals, actual)
	}

	// assert
	for i, expected := range expecteds {
		actual := actuals[i]
		if len(expected) == 0 && len(actual) == 0 {
			continue
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("inputs: %v expected: %v actual: %v", inputs[i], expected, actual)
		}
	}
}

func TestSubtract(t *testing.T) {
	// arrange
	inputs := [][][]int{
		[][]int{[]int{1, 2, 3}, []int{4, 5, 6}},
		[][]int{[]int{1, 2, 3}, []int{3, 2, 1}},
		[][]int{[]int{1, 2, 3}, []int{2, 2, 4}},
		[][]int{[]int{1, 2, 3}, []int{3, 4, 5}},
	}

	expecteds := [][]int{
		[]int{1, 2, 3},
		[]int{},
		[]int{1, 3},
		[]int{1, 2},
	}

	// act
	var actuals [][]int
	for _, input := range inputs {
		actual := subtract(input[0], input[1])
		actuals = append(actuals, actual)
	}

	// assert
	for i, expected := range expecteds {
		actual := actuals[i]
		if len(expected) == 0 && len(actual) == 0 {
			continue
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("inputs: %v expected: %v actual: %v", inputs[i], expected, actual)
		}
	}
}
