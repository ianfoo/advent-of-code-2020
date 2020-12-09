// This is a program that will compute the answers to both parts of the Day 6
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/6)
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

	{
		result, err := part1(input)
		if err != nil {
			return fmt.Errorf("part 1: %w", err)
		}
		fmt.Printf("Part 1: %d\n", result)
	}

	{
		result, err := part2(input)
		if err != nil {
			return fmt.Errorf("part 2: %w", err)
		}
		fmt.Printf("Part 2: %d\n", result)
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

func part1(input []string) (int, error) {
	var (
		sum    int
		yesses = make(map[rune]struct{})
	)
	for _, line := range input {
		if line == "" {
			sum += len(yesses)
			yesses = make(map[rune]struct{})
			continue
		}
		for _, question := range line {
			yesses[question] = struct{}{}
		}
	}

	// Don't forget the last group.
	sum += len(yesses)

	return sum, nil
}

func part2(input []string) (int, error) {
	var (
		sum       int
		yesses    = make(map[rune]int)
		groupSize int
	)
	countUnanimous := func(yesses map[rune]int, groupSize int) int {
		unanimousYes := make([]rune, 0, len(yesses))
		for question, count := range yesses {
			if count == groupSize {
				unanimousYes = append(unanimousYes, question)
			}
		}
		return len(unanimousYes)
	}
	for _, line := range input {
		if line == "" {
			numUnanimousYes := countUnanimous(yesses, groupSize)
			sum += numUnanimousYes

			// Reset the tracking state.
			yesses = make(map[rune]int)
			groupSize = 0
			continue
		}

		groupSize++
		for _, question := range line {
			yesses[question] += 1
		}
	}

	// Don't forget the last group.
	sum += countUnanimous(yesses, groupSize)

	return sum, nil
}
