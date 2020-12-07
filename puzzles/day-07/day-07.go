// This is a program that will compute the answers to both parts of the Day 7
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/7)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const myBagColor = "shiny gold"

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
		log.Printf(format, params...)
	}
}

func run(r io.Reader) error {
	input, err := readInput(r)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	if err := parseRules(input); err != nil {
		return err
	}

	part1Answer := Part1_HowManyColorsCanContain(myBagColor)
	part2Answer := Part2_NumberOfBagsContainedByBagColor(myBagColor)

	fmt.Printf("Part 1 Answer: %d\n", part1Answer)
	fmt.Printf("Part 2 Answer: %d\n", part2Answer)

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

var containmentRules map[string]map[string]int

var (
	rulePat        = regexp.MustCompile(`^(.+) bags contain (.+)$`)
	containablePat = regexp.MustCompile(`^(\d+) (.+) bags?\.?$`)
)

func parseRules(input []string) error {
	containmentRules = make(map[string]map[string]int)
	for i, line := range input {
		match := rulePat.FindStringSubmatch(line)
		if len(match) != 3 {
			return fmt.Errorf("unexpected container rule format: %q", line)
		}

		containerColor := match[1]
		sequence := i + 1

		if strings.HasSuffix(line, "no other bags.") {
			// Bags that cannot contain any other bags are 'leaf nodes" in the
			// bag capacity tree. When they're encountered programmatically, it
			// indicates the end of a recursion descention.
			containmentRules[containerColor] = nil
			trace("[%4d] rule: %s contains nothing\n", sequence, containerColor)
			continue
		}

		containableStr := match[2]
		containableStr = strings.Trim(containableStr, " ")
		containableList := strings.Split(containableStr, ", ")
		trace("[%4d] rule: %s contains %s\n", sequence, containerColor, containableStr)

		for _, containable := range containableList {
			color, capacity, err := parseCapacityForColor(containable)
			if err != nil {
				return err
			}
			if _, ok := containmentRules[containerColor]; !ok {
				containmentRules[containerColor] = make(map[string]int)
			}
			containmentRules[containerColor][color] = capacity
		}
	}
	return nil
}

func parseCapacityForColor(caStr string) (string, int, error) {
	ruleParams := containablePat.FindStringSubmatch(caStr)
	if len(ruleParams) != 3 {
		return "", 0, fmt.Errorf("unexpected capacity format: %q", caStr)
	}

	numStr := ruleParams[1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "", 0, fmt.Errorf("bad number %q in capacity: %w", numStr, err)
	}

	color := ruleParams[2]

	return color, num, nil
}

func Part1_HowManyColorsCanContain(targetColor string) int {
	count := 0
	for color, r := range containmentRules {
		if color == targetColor {
			// Bag of target color being able to hold bags of its own color
			// would lead to infinite bags, so we'll skip that to not be silly.
			continue
		}
		if canContain(targetColor, r) {
			trace("[%4d] %q can contain %q", count, color, targetColor)
			count++
		}
	}
	return count
}

func canContain(targetColor string, containableBags map[string]int) bool {
	isContainable := false
	for color := range containableBags {
		if color == targetColor {
			return true
		}
		subContainable, ok := containmentRules[color]
		if !ok {
			// I would normally send an error back up instead of panic, but in
			// this case cutting out excess error-handling simplifies the code
			// enough--especially in Go, in these reecursive functions--to be
			// worth it.
			panic(fmt.Sprintf("unknown color %q", color))
		}
		isContainable = isContainable || canContain(targetColor, subContainable)
	}
	return isContainable
}

func Part2_NumberOfBagsContainedByBagColor(color string) int {
	total := numberOfBagsForBagColor(color)

	// Do not count the outermost bag! We cannot do this inside the recursively
	// called function because it would misrepresent the number of bags nested.
	// We'll just adjust it here for the single case.
	numInside := total - 1

	return numInside
}

func numberOfBagsForBagColor(color string) int {
	containedBags, ok := containmentRules[color]
	if !ok {
		panic(fmt.Sprintf("unknown color %q", color))
	}
	count := 1
	for color, capacity := range containedBags {
		numBagsInOneBag := numberOfBagsForBagColor(color)
		count += capacity * numBagsInOneBag
	}
	trace("One %q bag contains %d bags inside it", color, count-1)
	return count
}
