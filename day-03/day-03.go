// This is a program that will compute the answers to both parts of the Day 3
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/3)
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
	const step = 3
	var numTrees int

	for n, line := range input {
		offset := (n * step) % len(line)
		if line[offset] == '#' {
			numTrees++
		}
	}
	return numTrees, nil
}

func part2(input []string) (int, error) {
	var patterns = []struct {
		over, down int
	}{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}

	var product = 1

	for _, pattern := range patterns {
		numTrees := 0
		for n, line := range input {
			// Skip the current line if its index doesn't divide cleanly into
			// the number of lines we should move down each step.
			if n%pattern.down != 0 {
				continue
			}
			offset := (n * pattern.over) % len(line)
			if line[offset] == '#' {
				numTrees++
			}
		}
		product *= numTrees
	}
	return product, nil
}
