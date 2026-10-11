package main

import (
	"os"
	"testing"
)

func TestParseInt(t *testing.T) {
	p, err := parsePacket("42")
	if err != nil {
		t.Fatal(err)
	}
	if val, ok := p.(float64); !ok || val != 42 {
		t.Errorf("expected 42, got %v", p)
	}
}

func TestParseList(t *testing.T) {
	p, err := parsePacket("[]")
	if err != nil {
		t.Fatal(err)
	}
	if list, ok := p.([]interface{}); !ok || len(list) != 0 {
		t.Errorf("expected empty list, got %v", p)
	}
}

func TestParseComplex(t *testing.T) {
	p, err := parsePacket("[1,[2,3],4]")
	if err != nil {
		t.Fatal(err)
	}
	if list, ok := p.([]interface{}); !ok || len(list) != 3 {
		t.Errorf("expected list of length 3, got %v", p)
	}
}

func TestCompareInts(t *testing.T) {
	a := float64(1)
	b := float64(2)
	if cmp := compare(a, b); cmp >= 0 {
		t.Errorf("expected 1 < 2, got cmp=%d", cmp)
	}
}

func TestCompareIntVsList(t *testing.T) {
	a := float64(2)
	b := []interface{}{float64(3)}
	if cmp := compare(a, b); cmp >= 0 {
		t.Errorf("expected [2] < [3], got cmp=%d", cmp)
	}
}

func TestPart1Example(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	result, err := solvePart1(string(input))
	if err != nil {
		t.Fatal(err)
	}
	if result != 13 {
		t.Errorf("expected 13, got %d", result)
	}
}

func TestPart2Example(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	result, err := solvePart2(string(input))
	if err != nil {
		t.Fatal(err)
	}
	if result != 140 {
		t.Errorf("expected 140, got %d", result)
	}
}

func TestMalformedEmpty(t *testing.T) {
	_, err := parsePacket("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestMalformedInvalid(t *testing.T) {
	_, err := parsePacket("abc")
	if err == nil {
		t.Error("expected error for invalid input")
	}
}
