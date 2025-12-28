package main

import (
	"reflect"
	"testing"
)

func TestFindIntersection_Basic(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{4, 5, 6, 7}
	ok, res := FindIntersection(slice1, slice2)

	if !ok {
		t.Fatalf("expected ok=true, got false")
	}

	expected := map[int]bool{4: true, 5: true}
	if len(res) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(res))
	}
	for _, v := range res {
		if !expected[v] {
			t.Errorf("unexpected value in result: %d", v)
		}
		delete(expected, v)
	}
	if len(expected) != 0 {
		t.Errorf("some expected values are missing: %v", expected)
	}
}

func TestFindIntersection_NoIntersection(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}
	ok, res := FindIntersection(slice1, slice2)

	if ok {
		t.Fatalf("expected ok=false, got true")
	}
	if len(res) != 0 {
		t.Errorf("expected empty slice, got %v", res)
	}
}

func TestFindIntersection_EmptyFirst(t *testing.T) {
	var slice1 []int
	slice2 := []int{1, 2, 3}
	ok, res := FindIntersection(slice1, slice2)

	if ok {
		t.Fatalf("expected ok=false, got true")
	}
	if len(res) != 0 {
		t.Errorf("expected empty slice, got %v", res)
	}
}

func TestFindIntersection_EmptySecond(t *testing.T) {
	slice1 := []int{1, 2, 3}
	var slice2 []int
	ok, res := FindIntersection(slice1, slice2)

	if ok {
		t.Fatalf("expected ok=false, got true")
	}
	if len(res) != 0 {
		t.Errorf("expected empty slice, got %v", res)
	}
}

func TestFindIntersection_BothEmpty(t *testing.T) {
	var slice1 []int
	var slice2 []int
	ok, res := FindIntersection(slice1, slice2)

	if ok {
		t.Fatalf("expected ok=false, got true")
	}
	if len(res) != 0 {
		t.Errorf("expected empty slice, got %v", res)
	}
}

func TestFindIntersection_WithDuplicates(t *testing.T) {
	slice1 := []int{1, 2, 2, 3, 3, 3}
	slice2 := []int{2, 2, 3}
	ok, res := FindIntersection(slice1, slice2)

	if !ok {
		t.Fatalf("expected ok=true, got false")
	}

	expectedSet := map[int]bool{2: true, 3: true}
	if len(res) != len(expectedSet) {
		t.Fatalf("expected length %d, got %d", len(expectedSet), len(res))
	}
	for _, v := range res {
		if !expectedSet[v] {
			t.Errorf("unexpected value in result: %d", v)
		}
		delete(expectedSet, v)
	}
	if len(expectedSet) != 0 {
		t.Errorf("some expected values are missing: %v", expectedSet)
	}
}

func TestFindIntersection_SingleElementIntersection(t *testing.T) {
	slice1 := []int{10, 20, 30}
	slice2 := []int{5, 20, 40}
	ok, res := FindIntersection(slice1, slice2)

	if !ok {
		t.Fatalf("expected ok=true, got false")
	}

	if len(res) != 1 {
		t.Fatalf("expected length 1, got %d (%v)", len(res), res)
	}
	if res[0] != 20 {
		t.Errorf("expected [20] (в любом месте), got %v", res)
	}
}

func TestFindIntersection_NegativeAndZero(t *testing.T) {
	slice1 := []int{-1, 0, 1}
	slice2 := []int{-1, 2, 3}
	ok, res := FindIntersection(slice1, slice2)

	if !ok {
		t.Fatalf("expected ok=true, got false")
	}

	expected := []int{-1}
	if !reflect.DeepEqual(res, expected) && !(len(res) == 1 && res[0] == -1) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}
