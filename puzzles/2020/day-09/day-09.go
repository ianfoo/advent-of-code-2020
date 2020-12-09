// This is a program that will compute the answers to both parts of the Day 9
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/9)
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

// Sample data uses differnet value, so allow it to be changed.
const DefaultCypherSize = 25

var (
	cypherSize int
	verbose    bool
)

func main() {
	flag.BoolVar(&verbose, "verbose", false, "Print more detail.")
	flag.IntVar(&cypherSize, "cypher-size", DefaultCypherSize, "Change size of cypher.")
	flag.Parse()

	log.SetFlags(log.Ltime | log.Lmicroseconds)

	if err := run(os.Stdin, cypherSize); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func trace(format string, params ...interface{}) {
	if verbose {
		log.Printf(format, params...)
	}
}

func run(r io.Reader, cypherSize int) error {
	input, err := readInput(r)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	invalidSum, err := part1(input, cypherSize)
	if err != nil {
		return fmt.Errorf("part 1: %w", err)
	}
	fmt.Printf("Part 1: %d\n", invalidSum)

	result, err := part2(input, invalidSum)
	if err != nil {
		return fmt.Errorf("part 2: %w", err)
	}
	fmt.Printf("Part 2: %d\n", result)

	return nil
}

func readInput(r io.Reader) ([]int, error) {
	var (
		s     = bufio.NewScanner(r)
		input []int
	)
	for s.Scan() {
		line := s.Text()
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
		input = append(input, n)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return input, nil
}

func isValid(window []int, targetSum int) bool {
	for i, x := range window {
		for j := i + 1; j < len(window); j++ {
			y := window[j]
			if sum := x + y; sum == targetSum {
				return true
			}
		}
	}
	return false
}

func part1(input []int, cypherSize int) (int, error) {
	trace("input: %v", input)
	for i := cypherSize; i < len(input); i++ {

		// Define window of input for building sums.
		var (
			lower = i - cypherSize
			upper = i
		)
		window := input[lower : upper : cypherSize+lower]
		trace("window: elements %d-%d: %v", lower, upper, window)

		// Check validity of number following window.
		if n := input[i]; !isValid(window, n) {
			trace("%d on line %d is not a sum of two entries in window", n, i+1)
			return n, nil
		}
	}

	return 0, errors.New("no invalid data found")
}

func part2(input []int, targetSum int) (int, error) {
	sumRange := func(r []int) int {
		sum := 0
		for _, x := range r {
			sum += x
		}
		return sum
	}

	minAndMax := func(numbers []int) (int, int) {
		var (
			min = math.MaxInt64
			max = 0
		)
		for _, x := range numbers {
			if x < min {
				min = x
				continue
			}
			if x > max {
				max = x
			}
		}
		return min, max
	}

	// Find contiguous numbers that add up to sum.
LOOP:
	for lower := 0; lower < len(input)-1; lower++ {
		for upper := lower + 1; upper < len(input); upper++ {
			window := input[lower:upper]
			sum := sumRange(window)
			if sum > targetSum {
				// Overshot our target, so start next iteration.
				continue LOOP
			}
			if sum == targetSum {
				min, max := minAndMax(window)
				return min + max, nil
			}
		}
	}

	return 0, errors.New("no range matching sum found")
}
