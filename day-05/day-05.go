// This is a program that will compute the answers to both parts of the Day 5
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/5)
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(r io.Reader) error {
	input, err := readInput(r)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	var (
		minSeatID int
		maxSeatID int
	)
	{
		var err error
		minSeatID, maxSeatID, err = findMinAndMaxSeatID(input)
		if err != nil {
			return fmt.Errorf("part 1: %w", err)
		}
		fmt.Printf("Maximum Seat ID: %d\n", maxSeatID)
	}

	{
		missingSeatID, err := findMissingSeatID(input, minSeatID, maxSeatID)
		if err != nil {
			return fmt.Errorf("part 2: %w", err)
		}
		fmt.Printf("Missing Seat ID: %d\n", missingSeatID)
	}

	return nil
}

func readInput(r io.Reader) ([]string, error) {
	var (
		s     = bufio.NewScanner(r)
		input []string
	)
	for s.Scan() {
		line := s.Text()
		input = append(input, line)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return input, nil
}

// decodeSeat walks through the boarding pass string and decodes the
// row and seat, and calculates the seat ID.
func decodeSeat(s string) int {
	const (
		rowIDLength = 7
		colIDLength = 3
	)
	var (
		row int
		col int
	)
	for i := 0; i < rowIDLength; i++ {
		if s[i] == 'B' {
			pow := (rowIDLength - 1) - i
			row += 1 << pow
		}
	}
	for i := 0; i < colIDLength; i++ {
		offset := rowIDLength + i
		if s[offset] == 'R' {
			pow := (colIDLength - 1) - i
			col += 1 << pow
		}
	}
	seatID := row*8 + col
	return seatID

}

func findMinAndMaxSeatID(passes []string) (int, int, error) {
	const maxPossibleSeatID = 127*8 + 8
	var (
		min = maxPossibleSeatID
		max = 0
	)
	for _, pass := range passes {
		seatID := decodeSeat(pass)
		if seatID > max {
			max = seatID
			continue
		}
		if seatID < min {
			min = seatID
		}
	}
	return min, max, nil
}

// findMissingSeatID finds which boarding pass is missing by using the knowledge
// that seat IDs are consecutive from minimum to maximum, and comparing the sum
// of those values with the sum of the values represented in the acutal list of
// boarding passes, which is missing one value. The difference between the sums
// will indicate which seat is missing.
func findMissingSeatID(passes []string, minSeatID, maxSeatID int) (int, error) {
	var (
		seatIDsPerfectSum int
		seatIDsActualSum  int
	)
	for i := minSeatID; i <= maxSeatID; i++ {
		seatIDsPerfectSum += i
	}
	for _, pass := range passes {
		seatID := decodeSeat(pass)
		seatIDsActualSum += seatID
	}

	missingSeatID := seatIDsPerfectSum - seatIDsActualSum
	return missingSeatID, nil
}
