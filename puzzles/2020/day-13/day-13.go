// This is a program that will compute the answers to both parts of the Day 13
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/13)
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

	earliest, buses, err := getEarliestDepartureAndBuses(input)
	if err != nil {
		return fmt.Errorf("parsing input: %w", err)
	}

	var part1Result int
	{
		var err error
		part1Result = part1(earliest, buses)
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

func getEarliestDepartureAndBuses(input []string) (int, []int, error) {
	earliest, err := strconv.Atoi(input[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid earliest departure time: %w", err)
	}

	busesStr := strings.Split(input[1], ",")
	buses := make([]int, 0, len(busesStr))
	for _, busIDStr := range busesStr {
		if busIDStr == "x" {
			continue
		}
		busID, err := strconv.Atoi(busIDStr)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid bus ID: %w", err)
		}
		buses = append(buses, busID)
	}

	return earliest, buses, nil
}

func part1(earliest int, buses []int) int {
	minWait := earliest
	bestBus := 0
	for _, bus := range buses {
		mod := earliest % bus
		arrival := earliest - mod + bus
		wait := arrival - earliest
		trace("bus %d arrives %d minutes after %d", bus, mod, earliest)
		if wait < minWait {
			bestBus = bus
			minWait = wait
		}
	}
	trace("bus: %d wait: %d", bestBus, minWait)
	return minWait * bestBus
}

func part2(input []string) (int, error) {
	busesStr := strings.Split(input[1], ",")
	buses := make([]int, 0, len(busesStr))
	maxBus := 0
	maxBusIndex := 0
	for i, busStr := range busesStr {
		if busStr == "x" {
			buses = append(buses, -1)
			continue
		}
		bus, err := strconv.Atoi(busStr)
		if err != nil {
			return 0, fmt.Errorf("invalid bus ID: %w", err)
		}
		buses = append(buses, bus)
		if bus > maxBus {
			maxBus = bus
			maxBusIndex = i
		}
	}

	fmt.Printf("%d bus routes\n", len(buses))
	var result int

LOOP:
	// My attempt. Works, but runs forever with actual input.
	for t := 0; ; t += maxBus {
		trace("t: %d", t)
		for offset, b := range buses {
			if b == -1 {
				continue
			}
			arrival := t + offset - maxBusIndex
			offBy := arrival % b
			if offBy != 0 {
				// trace("t: %d bus %d does not arrive at %d (off by %d): TRY AGAIN", t, b, arrival, offBy)
				continue LOOP
			}
			trace("bus %d arrives at %d", b, arrival)
		}
		result = t - maxBusIndex
		break
	}

	return result, nil

	/*
		Chinese remainder theorem solution that doesn't work.
		nums := make([]int, 0, len(buses))
		rem := make([]int, 0, len(buses))
		for i, b := range buses {
			if b == -1 {
				continue
			}
			nums = append(nums, b)
			rem = append(rem, i)
		}
		result := crt(nums, rem)
		return result, nil
	*/
}

// Absolutely no idea. Translated from C++ from
// https://www.geeksforgeeks.org/chinese-remainder-theorem-set-2-implementation/
func modInv(a, m int) int {
	if m == 1 {
		return 0
	}

	m0 := m
	x0, x1 := 0, 1
	for a > 1 {
		q := a / m
		t := m
		m = a % m
		a = t
		x0 = x1 - q*x0
		x1 = t
	}
	if x1 < 0 {
		x1 += m0
	}
	return x1
}

// Translated from this after finding pages of mathematical gobledygook
// describing the Chinese Remainder Theorem but nothing actionable (i.e., code).
//
// https://www.geeksforgeeks.org/chinese-remainder-theorem-set-2-implementation/
func crt(nums, rem []int) int {
	prod := 1
	for _, n := range nums {
		prod *= n
	}

	result := 0
	for i, bus := range nums {
		pp := prod / bus
		result += rem[i] * modInv(pp, bus) * pp
	}

	return result % prod
}
