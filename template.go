// Use this template to jumpstart Go-based solutions for Advent of Code.

// This is a program that will compute the answers to both parts of the Day N
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/N)
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
	// Write the code to complete part one of the puzzle here.
	var result int
	return result, nil
}

func part2(input []string) (int, error) {
	// Write the code to complete part two of the puzzle here.
	var result int
	return result, nil
}
