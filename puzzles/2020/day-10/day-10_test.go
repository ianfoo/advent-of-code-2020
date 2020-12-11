package main

import "testing"

func TestTracker(t *testing.T) {
	a := []int{4, 6, 9, 10}
	tr := make(Tracker)
	tr.Add(a)
	if !tr.Exists(a) {
		t.Error("expected slice to be reported as existing in tracker")
	}

	b := []int{4, 7, 12}
	if tr.Exists(b) {
		t.Error("expected non-added slice to not be reported as existing in tracker")
	}

	tr.Add(b)
	if !tr.Exists(a) {
		t.Error("expected initially-added slice to still exist in tracker")
	}

	if c := tr.CountAll(); c != 2 {
		t.Errorf("expected count of 2, but got %d", c)
	}
}
