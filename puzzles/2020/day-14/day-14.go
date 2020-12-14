// This is a program that will compute the answers to both parts of the Day 14
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/14)
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

	var part1Result int64
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

type Mask struct {
	And int64
	Or  int64
}

func (m Mask) Apply(n int64) int64 {
	n &= m.And
	n |= m.Or
	return n
}

func NewMask(s string) (Mask, error) {
	const (
		bits        = 36
		binaryOnes  = int64(2<<(bits-1) - 1)
		binaryZeros = int64(0)
	)
	var (
		and = binaryOnes
		or  = binaryZeros
	)
	for n, ch := range s {
		switch ch {
		case '1':
			or |= 1 << (bits - n - 1)
		case '0':
			and ^= 1 << (bits - n - 1)
		case 'X':
			// no-op
		default:
			return Mask{}, fmt.Errorf("invalid mask %q", s)
		}
	}
	m := Mask{
		And: and,
		Or:  or,
	}
	return m, nil
}

func runProgram(input []string) (map[int64]int64, error) {
	var (
		mask Mask
		mem  = make(map[int64]int64)
	)
	for i, line := range input {
		tokens := strings.Split(line, " = ")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("invalid input: %q", line)
		}
		if tokens[0] == "mask" {
			var err error
			mask, err = NewMask(tokens[1])
			if err != nil {
				return nil, err
			}
			continue
		}
		if !strings.HasPrefix(tokens[0], "mem") {
			return nil, fmt.Errorf("invalid input: %q", line)
		}

		var addr int64
		{
			instr := tokens[0]
			instr = strings.TrimPrefix(instr, "mem[")
			instr = strings.TrimSuffix(instr, "]")

			var err error
			addr, err = strconv.ParseInt(instr, 0, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid memory address %q on line %d: %w", instr, i+1, err)
			}
		}

		var value int64
		{
			valueStr := tokens[1]
			var err error
			value, err = strconv.ParseInt(valueStr, 0, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value %q on line %d: %w", valueStr, i+1, err)
			}
		}

		value = mask.Apply(value)
		mem[addr] = value
	}

	return mem, nil
}

func part1(input []string) (int64, error) {
	mem, err := runProgram(input)
	if err != nil {
		return 0, err
	}
	var sum int64
	for _, v := range mem {
		sum += v
	}
	return sum, nil
}

func part2(input []string) (int, error) {
	var result int

	// Write the code to complete part two of the puzzle here.

	return result, nil
}
