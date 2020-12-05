package main

import (
	"testing"
)

func TestDecodeSeat(t *testing.T) {
	tt := []struct {
		boardingPass   string
		expectedSeatID int
	}{
		// Examples provided in problem description.
		{"BFFFBBFRRR", 567}, // row 70, column 7, seat ID 567.
		{"FFFBBBFRRR", 119}, // row 14, column 7, seat ID 119.
		{"BBFFBBFRLL", 820}, // row 102, column 4, seat ID 820.
		{"FBFBBFFRLR", 357}, // row 44, col 5, seat ID 357.
	}
	for _, tc := range tt {
		t.Run(tc.boardingPass, func(t *testing.T) {
			if want, got := tc.expectedSeatID, decodeSeat(tc.boardingPass); want != got {
				t.Fatalf("expected seat ID %d, but got %d", want, got)
			}
		})
	}
}
