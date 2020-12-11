// This is a program that will compute the answers to both parts of the Day 11
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/11)
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var verbose bool

func main() {
	flag.BoolVar(&verbose, "verbose", false, "Print more detail.")
	flag.Parse()

	log.SetFlags(log.Ltime | log.Lmicroseconds)

	if err := run(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func trace(format string, params ...interface{}) {
	if verbose {
		fmt.Printf(format+"\n", params...)
	}
}

func run(r io.Reader) error {
	input, err := readInput(r)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	var part1Result int
	{
		var err error
		part1Result, err = part1(input)
		if err != nil {
			return fmt.Errorf("part 1: %w", err)
		}
	}

	var part2Result int
	{
		var err error
		part2Result, err = part2(input)
		if err != nil {
			return fmt.Errorf("part 2: %w", err)
		}
	}

	fmt.Printf("Part 1: %d\n", part1Result)
	fmt.Printf("Part 2: %d\n", part2Result)
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

// x and y are seat coordinates
func occupiedCountAroundSeat(x, y int, seats []string) int {
	var occupied int
	for i := y - 1; i <= y+1; i++ {
		if i < 0 || i >= len(seats) {
			// Out of bounds.
			continue
		}
		for j := x - 1; j <= x+1; j++ {
			if j < 0 || j >= len(seats[i]) {
				// Out of bounds.
				continue
			}
			if i == y && j == x {
				// This is the seat we're summing around: skip it.
				continue
			}
			if seats[i][j] == '#' {
				occupied++
			}
		}
	}
	return occupied
}

func nextIteration(
	rows []string,
	occCalcFunc func(int, int, []string) int,
	tooCrowdedThreshold int,
) ([]string, int) {
	var (
		next          []string
		occupiedCount int
	)
	for y, row := range rows {
		var newRow strings.Builder
		for x, seat := range row {
			if seat == '.' {
				// Skip the floor!
				newRow.WriteRune('.')
				continue
			}
			occ := occCalcFunc(x, y, rows)
			var newSeat rune
			switch {
			case seat == 'L' && occ == 0:
				newSeat = '#'
			case seat == '#' && occ >= tooCrowdedThreshold:
				newSeat = 'L'
			default:
				newSeat = seat
			}
			newRow.WriteRune(newSeat)
			if newSeat == '#' {
				occupiedCount++
			}
		}
		next = append(next, newRow.String())
	}
	return next, occupiedCount
}

func part1(rows Rows) (int, error) {
	var occupiedCount int

	for numRounds := 0; ; numRounds++ {
		var newOccupiedCount int
		rows, newOccupiedCount = nextIteration(rows, occupiedCountAroundSeat, 4)

		trace("%v", rows)
		trace("after round %d: %d occupied seats\n", numRounds+1, occupiedCount)
		if newOccupiedCount == occupiedCount {
			break
		}
		occupiedCount = newOccupiedCount
	}

	return occupiedCount, nil
}

func occupiedCountRays(x, y int, rows []string) int {
	var occupied int
	for dirY := -1; dirY <= 1; dirY++ {
		for dirX := -1; dirX <= 1; dirX++ {
			nearestX, nearestY, err := findNearestSeatInDirection(x, y, dirX, dirY, rows)
			if err != nil {
				// No seat found in this direction, so continue.
				continue
			}
			if rows[nearestY][nearestX] == '#' {
				occupied++
			}
		}
	}
	return occupied
}

func findNearestSeatInDirection(x, y, dirX, dirY int, rows []string) (int, int, error) {
	var (
		foundX int
		foundY int
	)
	if dirX == 0 && dirY == 0 {
		return 0, 0, errors.New("no direction provided")
	}
	for n := 1; ; n++ {
		checkY := y + (n * dirY)
		if checkY < 0 || checkY >= len(rows) {
			return 0, 0, errors.New("out of bounds (y)")
		}

		checkX := x + (n * dirX)
		if checkX < 0 || checkX >= len(rows[checkY]) {
			return 0, 0, errors.New("out of bounds (x)")
		}

		checkSeat := rows[checkY][checkX]
		if checkSeat != '.' {
			foundX = checkX
			foundY = checkY
			break
		}
	}
	return foundX, foundY, nil
}

type Rows []string

func (r Rows) String() string {
	s := ""
	for _, row := range r {
		s += row + "\n"
	}
	return s[:len(s)-1]
}

func part2(rows Rows) (int, error) {
	var occupiedCount int

	for numRounds := 0; ; numRounds++ {
		var newOccupiedCount int
		rows, newOccupiedCount = nextIteration(rows, occupiedCountRays, 5)

		trace("%v", rows)
		trace("after round %d: %d occupied seats\n", numRounds+1, occupiedCount)

		if newOccupiedCount == occupiedCount {
			break
		}
		occupiedCount = newOccupiedCount
	}

	return occupiedCount, nil
}
