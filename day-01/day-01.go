// This is a program that will compute the answers to both parts of the Day 1
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/1)
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

const target = 2020

func main() {
	if err := run(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}

func run(r io.Reader) error {
	ints, err := readInput(r)
	if err != nil {
		return err
	}

	twoSumTerms, err := findTwoSumTerms(target, ints)
	if err != nil {
		return fmt.Errorf("finding two sum terms: %v", err)
	}
	{
		var (
			x, y = twoSumTerms[0], twoSumTerms[1]
			sum  = x + y
			prod = x * y
		)
		fmt.Printf("%d + %d = %d; %d x %d = %d\n", x, y, sum, x, y, prod)
	}

	threeSumTerms, err := findThreeSumTerms(target, ints)
	if err != nil {
		return fmt.Errorf("finding three sum terms: %v", err)
	}
	{
		var (
			x, y, z = threeSumTerms[0], threeSumTerms[1], threeSumTerms[2]
			sum     = x + y + z
			prod    = x * y * z
		)
		fmt.Printf("%d + %d + %d = %d; %d x %d x %d = %d\n", x, y, z, sum, x, y, z, prod)
	}
	return nil
}

func readInput(r io.Reader) ([]int, error) {
	var (
		s     = bufio.NewScanner(r)
		input = make([]int, 0)
	)
	for s.Scan() {
		t := s.Text()
		n, err := strconv.Atoi(t)
		if err != nil {
			return nil, fmt.Errorf("could not convert input to number: %v", err)
		}
		input = append(input, n)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return input, nil
}

// findTwoSumTerms scans the list of input for two terms that add up to the
// target, shortening its search each time since elements already checked as the
// first term will be disqualified from consideration.
func findTwoSumTerms(target int, ints []int) ([2]int, error) {
	for term1Idx := 0; term1Idx < len(ints)-2; term1Idx++ {
		for term2Idx := term1Idx + 1; term2Idx < len(ints)-1; term2Idx++ {
			if x, y := ints[term1Idx], ints[term2Idx]; x+y == target {
				return [2]int{x, y}, nil
			}
		}
	}
	return [2]int{}, fmt.Errorf("no terms sum to %d", target)
}

// findThreeSumTerms scans the list of input for three terms that add up to the
// target, shortening its search each time since elements already checked as the
// first terms will be disqualified from consideration.
func findThreeSumTerms(target int, ints []int) ([3]int, error) {
	for term1Idx := 0; term1Idx < len(ints)-3; term1Idx++ {
		for term2Idx := term1Idx + 1; term2Idx < len(ints)-2; term2Idx++ {
			for term3Idx := term2Idx + 1; term3Idx < len(ints)-1; term3Idx++ {
				if x, y, z := ints[term1Idx], ints[term2Idx], ints[term3Idx]; x+y+z == target {
					return [3]int{x, y, z}, nil
				}
			}
		}
	}
	return [3]int{}, fmt.Errorf("no terms sum to %d", target)
}
