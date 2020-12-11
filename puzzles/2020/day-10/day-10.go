// This is a program that will compute the answers to both parts of the Day 10
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/10)
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
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

func readInput(r io.Reader) ([]int, error) {
	var (
		s     = bufio.NewScanner(r)
		input []int
	)
	for s.Scan() {
		line := s.Text()
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid input %s: %w", line, err)
		}
		input = append(input, n)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return input, nil
}

// Add built-in adapter to list of adapters. Need to know max value first, so
// sort them here.
func sortAndAddDeviceAdapter(adapters []int) []int {
	sort.Ints(adapters)

	max := adapters[len(adapters)-1]
	builtInAdapter := max + 3
	adapters = append(adapters, builtInAdapter)

	trace("adapters: %v", adapters)
	return adapters
}

// Returns a histogram of differences, where index of the slice is equal to
// difference. An error is returned if the adapters can not be connected to
// reach final joltage.
func canConnect(adapters []int) ([]int, bool) {
	const maxDiff = 3
	var (
		// Histogram of differences, where index is equal to difference.
		diffs = make([]int, maxDiff+1)

		// Track the current joltage that we've reached.
		currentJoltage = 0
	)

	// Step through joltages, selecting adpaters as we're able to.
	for _, outputJolts := range adapters {
		diff := outputJolts - currentJoltage
		if diff > maxDiff {
			return nil, false
		}
		diffs[diff]++
		currentJoltage = outputJolts
		trace(
			"%3d current + %d diff => %3d output diffs=%v",
			currentJoltage, diff, outputJolts, diffs)
	}

	return diffs, true
}

func part1(adapters []int) (int, error) {
	adapters = sortAndAddDeviceAdapter(adapters)
	diffs, connects := canConnect(adapters)
	if !connects {
		return 0, errors.New("adapters cannot connect")
	}
	prod := diffs[1] * diffs[3]
	return prod, nil
}

func connected(a []int, finalStage int) bool {
	if len(a) == 0 {
		return false
	}
	const maxDiff = 3
	last := 0
	a = append(a, finalStage)
	for _, n := range a {
		if n-last > maxDiff {
			return false
		}
		last = n
	}
	return true
}

type Tracker map[int][][]int

func (t Tracker) Add(a []int) {
	if t.Exists(a) {
		trace("exists in tracker: not adding: %v", a)
		return
	}
	key := a[0]
	val := t[key]
	val = append(val, a)
	t[key] = val
	trace("t[%d] = %v", key, val)
}

func (t Tracker) Exists(a []int) bool {
	key := a[0]
	bucket := t[key]
LOOP:
	for _, x := range bucket {
		if len(x) != len(a) {
			continue
		}
		for i := range x {
			if a[i] != x[i] {
				continue LOOP
			}
		}
		return true
	}
	return false
}

func (t Tracker) CountAll() int {
	count := 0
	for _, v := range t {
		count += len(v)
	}
	return count
}

func countChainsIter(a []int) int {
	finalStage := a[len(a)-1]
	if !connected(a, finalStage) {
		return 0
	}
	const maxDiff = 3
	t := make(Tracker)
	for i := 0; i < len(a)-1; i++ {
		for j := i + 1; j < len(a); j++ {
			if a[j]-a[i] >= maxDiff {
				break
			}
			b := make([]int, j-1, len(a)-1)
			copy(b, a[:i])
			b = append(b, a[j+1:]...)
			if !connected(b, finalStage) {
				break
			}
			trace("valid: i=%d j=%d %v %v", a[i], a[j], a[:j], a[j+1:])
			trace("adding b: %v", b)
			t.Add(b)
		}
	}
	return t.CountAll()
}

func count(adapters map[int]struct{}, joules, finalStage int, memo map[int]int) int {
	if n, ok := memo[joules]; ok {
		return n
	}
	if joules == finalStage {
		return 1
	}
	sum := 0
	for i := 1; i <= 3; i++ {
		candidate := joules + i
		if _, ok := adapters[candidate]; ok {
			sum += count(adapters, candidate, finalStage, memo)
		}
	}
	memo[joules] = sum
	return sum
}

func countChains(a []int, finalStage int, tracker Tracker) int {
	if tracker.Exists(a) {
		trace("already seen: %v", a)
		return 0
	}
	trace("checking connection for: %v", a)
	if !connected(a, finalStage) {
		trace("terminal condition: %v does not connect", a)
		return 0
	}

	trace("ADDING ONE: %v", a)
	tracker.Add(a)
	count := 1
	for i := 0; i < len(a)-1; i++ {
		trace("about to remove element at index %d: %v", i, a)
		b := make([]int, len(a))
		copy(b, a)
		b = append(b[:i], b[i+1:]...)
		trace("removed element: len(a)=%d, len(b)=%d", len(a), len(b))
		count += countChains(b, finalStage, tracker)
	}
	return count
}

func part2(adapters []int) (int, error) {
	adapters = sortAndAddDeviceAdapter(adapters)
	finalStage := adapters[len(adapters)-1]

	adapterMap := make(map[int]struct{}, len(adapters))
	for _, v := range adapters {
		adapterMap[v] = struct{}{}
	}

	memo := make(map[int]int)
	numChains := count(adapterMap, 0, finalStage, memo)

	// tracker := make(Tracker)
	// trace("part 2: adapters: %v", adapters)
	// numChains := countChains(adapters, finalStage, tracker)

	return numChains, nil
}
