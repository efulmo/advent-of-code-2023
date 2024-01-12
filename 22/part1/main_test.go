package main

import (
	"testing"
)

func TestToStringId(t *testing.T) {
	inputs := []struct {
		num uint16
		str string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
		{26, "BA"},
		{27, "BB"},
	}

	for _, input := range inputs {
		got := toStringId(input.num)
		if got != input.str {
			t.Errorf("input: %d. wanted %s, got %s", input.num, input.str, got)
		}
	}
}

func TestOverlapsWith(t *testing.T) {
	basePrj := Projection{1, 1, 1, 1}
	inputs := []struct {
		other Projection
		want  bool
	}{
		{Projection{2, 2, 1, 1}, false},
		{Projection{1, 1, 2, 2}, false},
		{Projection{2, 2, 2, 2}, false},
		{Projection{1, 2, 1, 2}, true},
	}

	for i, input := range inputs {
		got := basePrj.overlapsWith(input.other)
		if got != input.want {
			t.Errorf("prj #%d. wanted %t, got %t", i, input.want, got)
		}
	}
}
