package report

import (
	"testing"
)

func compare(results []DiffReport, exp []string, t *testing.T) {
	if len(results) != len(exp) {
		t.Errorf("Got %v, expected %v", results, exp)
		return
	}

	for i := range results {
		if results[i].String() != exp[i] {
			t.Errorf("Got %v, expected %v", results, exp)
			return
		}
	}
}

func TestDiffResources(t *testing.T) {
	a := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}}
	b := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}}
	exp := []string{}
	compare(Diff(a, b), exp, t)
}

func TestDiffResourcesA(t *testing.T) {
	a := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}}
	b := []Keyer{}
	exp := []string{`Element: "  x1", A: exists, B: does not exist`, `Element: "  x2", A: exists, B: does not exist`}
	compare(Diff(a, b), exp, t)
}

func TestDiffResourcesB(t *testing.T) {
	a := []Keyer{}
	b := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}}
	exp := []string{`Element: "  x1", A: does not exist, B: exists`, `Element: "  x2", A: does not exist, B: exists`}
	compare(Diff(a, b), exp, t)
}

func TestDiffResourcesATail(t *testing.T) {
	a := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}}
	b := []Keyer{Resource{Name: "x1"}}
	exp := []string{`Element: "  x2", A: exists, B: does not exist`}
	compare(Diff(a, b), exp, t)
}

func TestDiffResourcesBTail(t *testing.T) {
	a := []Keyer{Resource{Name: "x1"}}
	b := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}}
	exp := []string{`Element: "  x2", A: does not exist, B: exists`}
	compare(Diff(a, b), exp, t)
}

func TestDiffResourcesAHead(t *testing.T) {
	a := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}}
	b := []Keyer{Resource{Name: "x2"}}
	exp := []string{`Element: "  x1", A: exists, B: does not exist`}
	compare(Diff(a, b), exp, t)
}

func TestDiffResourcesBHead(t *testing.T) {
	a := []Keyer{Resource{Name: "x2"}}
	b := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}}
	exp := []string{`Element: "  x1", A: does not exist, B: exists`}
	compare(Diff(a, b), exp, t)
}

func TestDiffResourcesAInsert(t *testing.T) {
	a := []Keyer{Resource{Name: "x1"}, Resource{Name: "x3"}}
	b := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}, Resource{Name: "x3"}}
	exp := []string{`Element: "  x2", A: does not exist, B: exists`}
	compare(Diff(a, b), exp, t)
}

func TestDiffResourcesBInsert(t *testing.T) {
	a := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}, Resource{Name: "x3"}}
	b := []Keyer{Resource{Name: "x1"}, Resource{Name: "x3"}}
	exp := []string{`Element: "  x2", A: exists, B: does not exist`}
	compare(Diff(a, b), exp, t)
}

func TestDiffResourcesAInsert2(t *testing.T) {
	a := []Keyer{Resource{Name: "x1"}, Resource{Name: "x4"}}
	b := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}, Resource{Name: "x3"}, Resource{Name: "x4"}}
	exp := []string{`Element: "  x2", A: does not exist, B: exists`, `Element: "  x3", A: does not exist, B: exists`}
	compare(Diff(a, b), exp, t)
}

func TestDiffResourcesBInsert2(t *testing.T) {
	a := []Keyer{Resource{Name: "x1"}, Resource{Name: "x2"}, Resource{Name: "x3"}, Resource{Name: "x4"}}
	b := []Keyer{Resource{Name: "x1"}, Resource{Name: "x4"}}
	exp := []string{`Element: "  x2", A: exists, B: does not exist`, `Element: "  x3", A: exists, B: does not exist`}
	compare(Diff(a, b), exp, t)
}
