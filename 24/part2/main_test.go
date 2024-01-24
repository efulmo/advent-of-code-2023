package main

import (
	"fmt"
	"slices"
	"testing"
)

func TestPrimeFactors(t *testing.T) {
	inputs := []struct {
		num     uint
		factors []uint
	}{
		{0, []uint{}},
		{1, []uint{}},
		{2, []uint{2}},
		{3, []uint{3}},
		{4, []uint{2, 2}},
		{5, []uint{5}},
		{6, []uint{2, 3}},
		{7, []uint{7}},
		{9, []uint{3, 3}},
		{10, []uint{2, 5}},
		{12, []uint{2, 2, 3}},
		{15, []uint{3, 5}},
		{100, []uint{2, 2, 5, 5}},
	}

	for _, input := range inputs {
		t.Run(fmt.Sprintf("%d", input.num), func(t *testing.T) {
			got := primeFactors(input.num)
			if !slices.Equal(got, input.factors) {
				t.Errorf("input: %d. wanted %v, got %v", input.num, input.factors, got)
			}
		})
	}
}

func TestGetProducst(t *testing.T) {
	inputs := []struct {
		nums     []uint
		products []uint
	}{
		{[]uint{}, []uint{}},
		{[]uint{1}, []uint{1}},
		{[]uint{3}, []uint{3}},
		{[]uint{1, 2}, []uint{2}},
		{[]uint{2, 2, 3}, []uint{4, 6, 12}},
		{[]uint{1, 2, 3, 5}, []uint{2, 3, 5, 6, 10, 15, 30}},
	}

	for _, input := range inputs {
		testName := fmt.Sprintf("%v", input.nums)
		t.Run(testName, func(t *testing.T) {
			got := getProducts(input.nums)
			gotSl := make([]uint, 0, len(got))
			for v := range got {
				gotSl = append(gotSl, v)
			}
			slices.Sort(gotSl)

			if !slices.Equal(gotSl, input.products) {
				t.Errorf("input: %d. wanted %v, got %v", input.nums, input.products, gotSl)
			}
		})
	}
}
