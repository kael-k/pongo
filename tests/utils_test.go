package tests

import (
	"reflect"
	"testing"

	"github.com/kael-k/pongo/pongo"
)

type testListDiffCase[T comparable] struct {
	testA []T
	testB []T
	want  []T
}

var testListDiffCases = []testListDiffCase[string]{
	{
		testA: []string{"A", "B", "C"},
		testB: []string{"A", "C"},
		want:  []string{"B"},
	},
	{
		testA: []string{"A", "B", "C"},
		testB: []string{"B", "A"},
		want:  []string{"C"},
	},
	{
		testA: []string{"A", "B", "C"},
		testB: []string{"A", "C", "B"},
		want:  []string{},
	},
	{
		testA: []string{"A", "B", "C"},
		testB: []string{},
		want:  []string{"A", "B", "C"},
	},
}

func TestListDiff(t *testing.T) {
	var originalTestA *[]string

	for _, v := range testListDiffCases {
		originalTestA = &[]string{}
		*originalTestA = append(*originalTestA, v.testA...)

		if l := pongo.ListDiff[string](v.testA, v.testB); !reflect.DeepEqual(l, v.want) {
			t.Errorf("error in ListDiff: %v - %v expected %v, got %v", v.testA, v.testB, v.want, l)
		}
		if !reflect.DeepEqual(*originalTestA, v.testA) {
			t.Errorf("error in ListMapDiff: diff altered B dict, expected %v, got %v", *originalTestA, v.testA)
		}
	}
}

type testListMapDiffCase[T comparable] struct {
	testA []T
	testB map[T]interface{}
	want  []T
}

var testListMapDiffCases = []testListMapDiffCase[string]{
	{
		testA: []string{"A", "B", "C"},
		testB: map[string]interface{}{"A": "test", "C": "test"},
		want:  []string{"B"},
	},
	{
		testA: []string{"A", "B", "C"},
		testB: map[string]interface{}{"B": "test", "A": "test"},
		want:  []string{"C"},
	},
	{
		testA: []string{"A", "B", "C"},
		testB: map[string]interface{}{"A": "test", "B": "test", "C": "test"},
		want:  []string{},
	},
	{
		testA: []string{"A", "B", "C"},
		testB: map[string]interface{}{},
		want:  []string{"A", "B", "C"},
	},
}

func TestListMapDiff(t *testing.T) {
	var originalTestB *map[string]interface{}
	for _, v := range testListMapDiffCases {
		originalTestB = &map[string]interface{}{}
		for k, s := range v.testB {
			(*originalTestB)[k] = s
		}
		if l := pongo.ListMapDiff[string](v.testA, v.testB); !reflect.DeepEqual(l, v.want) {
			t.Errorf("error in ListMapDiff: %v - %v expected %v, got %v", v.testA, v.testB, v.want, l)
		}
		if !reflect.DeepEqual(*originalTestB, v.testB) {
			t.Errorf("error in ListMapDiff: diff altered B dict, expected %v, got %v", *originalTestB, v.testB)
		}
	}
}

func InArray[T comparable](value T, array []T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
