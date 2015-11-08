package main

import (
	"fmt"
	"testing"
)

func TestCheckHiddenNPermutations(t *testing.T) {
	// TODO
	b := &board{}
	perms := b.getPermutations(3, []int{2, 3, 4, 5, 6}, []int{1})
	fmt.Printf("%v\n", perms)
}
