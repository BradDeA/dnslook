package main

import (
	"slices"
	"strconv"
	"testing"
)

func TestToString(t *testing.T) {
	input := []int{1, 2, 3}
	got := ToString(input, strconv.Itoa)
	want := []string{"1", "2", "3"}

	emptyInt := []int{}
	i := ToString(emptyInt, strconv.Itoa)
	o := []string{}

	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if !slices.Equal(i, o) {
		t.Errorf("got %v, want %v", i, o)
	}
}
