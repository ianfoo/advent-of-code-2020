// This is a program that will compute the answers to both parts of the Day 12
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/12)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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

	var part1Result int
	{
		var err error
		part1Result, err = part1(input)
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

func part1(input []string) (int, error) {
	var (
		movement          = make(map[string]int)
		currentDir string = "E"
	)
	for i, line := range input {
		val := string(line[0])
		magStr := line[1:]
		mag, err := strconv.Atoi(magStr)
		if err != nil {
			return 0, fmt.Errorf("invalid input magnitude %q on line %d: %w", magStr, i+1, err)
		}

		switch val {
		case "R":
			times := mag / 90
			for n := 0; n < times; n++ {
				switch currentDir {
				case "E":
					currentDir = "S"
				case "S":
					currentDir = "W"
				case "W":
					currentDir = "N"
				case "N":
					currentDir = "E"
				}
			}
			continue
		case "L":
			times := (mag / 90) % 4
			for n := 0; n < times; n++ {
				switch currentDir {
				case "E":
					currentDir = "N"
				case "N":
					currentDir = "W"
				case "W":
					currentDir = "S"
				case "S":
					currentDir = "E"
				}
			}
			continue
		case "E", "W", "N", "S":
			movement[val] += mag
		case "F":
			movement[currentDir] += mag
		default:
			return 0, fmt.Errorf("invalid input direction %s on line %d", val, i+1)
		}
		trace("%v", movement)
	}
	eastWest := movement["E"] - movement["W"]
	if eastWest < 0 {
		eastWest *= -1
	}
	northSouth := movement["N"] - movement["S"]
	if northSouth < 0 {
		northSouth *= -1
	}
	fmt.Printf("movement matrix: %+v\n", movement)
	fmt.Printf("east-west: %d, north-south: %d\n", eastWest, northSouth)

	return eastWest + northSouth, nil
}

func part2(input []string) (int, error) {
	var (
		waypoint = map[string]int{
			"E": 10,
			"N": 1,
			"W": 0,
			"S": 0,
		}
		location = map[string]int{
			"E": 0,
			"N": 0,
			"W": 0,
			"S": 0,
		}
	)
	for i, line := range input {
		val := string(line[0])
		magStr := line[1:]
		mag, err := strconv.Atoi(magStr)
		if err != nil {
			return 0, fmt.Errorf("invalid input magnitude %q on line %d: %w", magStr, i+1, err)
		}

		switch val {
		case "E", "N", "W", "S":
			waypoint[val] += mag
		case "F":
			for k := range waypoint {
				location[k] += mag * waypoint[k]
			}
		case "R":
			times := (mag / 90) % 4
			for i := 0; i < times; i++ {
				south := waypoint["S"]
				waypoint["S"] = waypoint["E"]
				waypoint["E"] = waypoint["N"]
				waypoint["N"] = waypoint["W"]
				waypoint["W"] = south
			}
		case "L":
			times := mag / 90
			for i := 0; i < times; i++ {
				south := waypoint["S"]
				waypoint["S"] = waypoint["W"]
				waypoint["W"] = waypoint["N"]
				waypoint["N"] = waypoint["E"]
				waypoint["E"] = south
			}
		}
	}

	eastWest := location["E"] - location["W"]
	if eastWest < 0 {
		eastWest *= -1
	}
	northSouth := location["N"] - location["S"]
	if northSouth < 0 {
		northSouth *= -1
	}

	return eastWest + northSouth, nil
}
