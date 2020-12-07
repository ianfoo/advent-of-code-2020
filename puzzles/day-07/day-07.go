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

	fmt.Println("verbose:", verbose)
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

type CapacityForColor struct {
	Number int
	Color  string
}

var containmentRules map[string][]CapacityForColor

var (
	rulePat        = regexp.MustCompile(`^(.+) contain (.+)$`)
	containablePat = regexp.MustCompile(`^(\d+) (.+)$`)
)

func parseRules(input []string) error {
	containmentRules = make(map[string][]CapacityForColor)
	for i, line := range input {
		match := rulePat.FindStringSubmatch(line)
		if len(match) != 3 {
			return fmt.Errorf("unexpected container rule format: %q", line)
		}

		color := match[1]
		color = strings.TrimSuffix(color, " bags")

		if strings.HasSuffix(line, "no other bags.") {
			// Bags that cannot contain any other bags are 'leaf nodes" in the
			// bag capacity tree. When they're encountered programmatically, it
			// indicates the end of a recursion descention.
			containmentRules[color] = nil
			continue
		}

		containableStr := match[2]
		containableStr = strings.Trim(containableStr, " .")
		containableSet := strings.Split(containableStr, ", ")
		trace("[%4d] %s contains %s\n", i+i, color, containableSet)

		for _, containable := range containableSet {
			c, err := parseCapacity(containable)
			if err != nil {
				return err
			}
			containmentRules[color] = append(containmentRules[color], c)
		}
	}
	return nil
}

func parseCapacity(caStr string) (CapacityForColor, error) {
	ruleParams := containablePat.FindStringSubmatch(caStr)
	if len(ruleParams) != 3 {
		return CapacityForColor{}, fmt.Errorf("unexpected capacity format: %q", caStr)
	}

	numStr := ruleParams[1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return CapacityForColor{}, fmt.Errorf("bad number %q in capacity: %w", numStr, err)
	}

	color := ruleParams[2]
	color = strings.TrimSuffix(color, " bag")
	color = strings.TrimSuffix(color, " bags")

	bagRule := CapacityForColor{
		Number: num,
		Color:  color,
	}
	return bagRule, nil
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
			// Current bag color can contain the target color bag, so add it to
			// count.
			count++

			trace("[%4d] %q can contain %q", count, color, targetColor)
		}
	}
	return count
}

func canContain(
	targetColor string,
	containableBags []CapacityForColor,
) bool {
	isContainable := false
	for _, c := range containableBags {
		color := c.Color
		if color == targetColor {
			return true
		}
		subRules, ok := containmentRules[color]
		if !ok {
			// I would normally send an error back up instead of panic, but in
			// this case cutting out excess error-handling simplifies the code
			// enough--especially in Go, in these reecursive functions--to be
			// worth it.
			panic(fmt.Sprintf("unknown color %q", color))
		}
		isContainable = isContainable || canContain(targetColor, subRules)
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
	for _, c := range containedBags {
		numBagsInOneBag := numberOfBagsForBagColor(c.Color)
		count += c.Number * numBagsInOneBag
	}
	trace("One %q bag contains %d bags inside it", color, count-1)
	return count
}
