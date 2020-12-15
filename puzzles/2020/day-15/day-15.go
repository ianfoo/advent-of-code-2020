// This is a program that will compute the answers to both parts of the Day 15
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/15)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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
	const part1LastTurn = 2020
	{
		var err error
		part1Result, err = NthRoundNumber(input, part1LastTurn)
		if err != nil {
			return fmt.Errorf("part 1: %w", err)
		}
	}
	fmt.Printf("Part 1: Number spoken during turn %d: %d\n", part1LastTurn, part1Result)

	// NOTE: Part 2 ran for 6-7 seconds on my laptop, so there is surely a way
	// to optimize this. Perhaps a pattern that can be detected and used to pick
	// out the nth turn's number quickly.
	const part2LastTurn = 30000000
	var part2Result int
	{
		var err error
		part2Result, err = NthRoundNumber(input, part2LastTurn)
		if err != nil {
			return fmt.Errorf("part 2: %w", err)
		}
	}
	fmt.Printf("Part 2: Number spoken during turn %d: %d\n", part2LastTurn, part2Result)

	return nil
}

func readInput(r io.Reader) ([]int, error) {
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

	numStr := strings.Split(input[0], ",")
	nums := make([]int, 0, len(numStr))
	for _, s := range numStr {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}
	return nums, nil
}

func NthRoundNumber(startingNums []int, lastTurn int) (int, error) {
	m := make(map[int][2]int)

	for i, n := range startingNums {
		turn := i + 1
		trace("[turn %3d] starting number: %d", turn, n)
		m[n] = [2]int{turn, 0}
	}

	var mostRecent = startingNums[len(startingNums)-1]
	for turn := len(startingNums) + 1; turn <= lastTurn; turn++ {
		// Determine the next number we should say.
		lastMentions := m[mostRecent]
		age := lastMentions[0] - lastMentions[1]
		if lastMentions[1] == 0 {
			// This number hasn't been heard previously.
			age = 0
		}

		trace("[turn %3d] mostRecent: %d lastMentions: %v age: %d", turn, mostRecent, lastMentions, age)

		// Update history for number we're saying now.
		ageMostRecent := m[age]
		m[age] = [2]int{turn, ageMostRecent[0]}

		// Record that we most recently said this number.
		mostRecent = age
	}

	return mostRecent, nil
}
