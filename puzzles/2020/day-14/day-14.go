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

	instructions, err := parseProgram(input)
	if err != nil {
		return fmt.Errorf("parsing program: %w", err)
	}

	var part1Result int64
	{
		var err error
		part1Result, err = part1(instructions)
		if err != nil {
			return fmt.Errorf("part 1: %w", err)
		}
	}

	var part2Result int64
	{
		var err error
		part2Result, err = part2(instructions)
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

type Instruction struct {
	Type    int
	Payload string
	Addr    int64
}

const (
	SetMaskInstruction = iota + 1
	SetMemInstruction
)

func parseProgram(input []string) ([]Instruction, error) {
	instructions := make([]Instruction, 0, len(input))
	for i, line := range input {
		tokens := strings.Split(line, " = ")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("invalid input: %q", line)
		}

		op, payload := tokens[0], tokens[1]

		if op == "mask" {
			instr := Instruction{
				Type:    SetMaskInstruction,
				Payload: payload,
			}
			instructions = append(instructions, instr)
			continue
		}
		if !strings.HasPrefix(op, "mem") {
			return nil, fmt.Errorf("invalid input: %q", line)
		}

		var addr int64
		{
			addrStr := op
			addrStr = strings.TrimPrefix(addrStr, "mem[")
			addrStr = strings.TrimSuffix(addrStr, "]")

			var err error
			addr, err = strconv.ParseInt(addrStr, 0, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid memory address %q on line %d: %w", addrStr, i+1, err)
			}
		}
		instr := Instruction{
			Type:    SetMemInstruction,
			Addr:    addr,
			Payload: payload,
		}
		instructions = append(instructions, instr)
	}

	return instructions, nil
}

type MaskV1 struct {
	And int64
	Or  int64
}

func (m MaskV1) Apply(n int64) int64 {
	n &= m.And
	n |= m.Or
	return n
}

func NewMaskV1(s string) (MaskV1, error) {
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
			return MaskV1{}, fmt.Errorf("invalid mask %q", s)
		}
	}
	m := MaskV1{
		And: and,
		Or:  or,
	}
	return m, nil
}

func part1(instructions []Instruction) (int64, error) {
	var (
		mask MaskV1
		mem  = make(map[int64]int64)
	)
	for i, instr := range instructions {
		var err error
		if instr.Type == SetMaskInstruction {
			mask, err = NewMaskV1(instr.Payload)
			if err != nil {
				return 0, err
			}
			continue
		}

		var value int64
		{
			valueStr := instr.Payload
			var err error
			value, err = strconv.ParseInt(valueStr, 0, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid value %q on line %d: %w", valueStr, i+1, err)
			}
		}
		value = mask.Apply(value)
		mem[instr.Addr] = value
	}

	var sum int64
	for _, v := range mem {
		sum += v
	}
	return sum, nil
}

type MaskV2 struct {
	Or        int64
	FloatBits int64
}

func NewMaskV2(s string) (MaskV2, error) {
	const bits int = 36
	var (
		or        int64
		floatBits int64
	)
	for n, ch := range s {
		switch ch {
		case '1':
			or |= 1 << (bits - n - 1)
		case 'X':
			floatBits |= 1 << (bits - n - 1)
		case '0':
			// no-op
		default:
			return MaskV2{}, fmt.Errorf("invalid mask %q", s)
		}
	}
	m := MaskV2{
		Or:        or,
		FloatBits: floatBits,
	}
	return m, nil
}

func (m MaskV2) SetValueAtAddresses(addr, val int64, mem map[int64]int64) {
	addr |= m.Or
	addr &^= m.FloatBits
	m.setValueAtAddresses(0, 0, addr, val, mem)
}

func (m MaskV2) setValueAtAddresses(exp int, sum, addr, val int64, mem map[int64]int64) {
	if 1<<exp > m.FloatBits {
		newAddr := addr + sum
		mem[newAddr] = val
		trace("setting value %d at address %d", val, newAddr)
		return
	}
	m.setValueAtAddresses(exp+1, sum, addr, val, mem)
	if n := int64(1 << exp); m.FloatBits&n != 0 {
		sum += n
		m.setValueAtAddresses(exp+1, sum, addr, val, mem)
	}
}

func part2(instructions []Instruction) (int64, error) {
	var (
		mask MaskV2
		mem  = make(map[int64]int64)
	)
	for i, instr := range instructions {
		var err error
		if instr.Type == SetMaskInstruction {
			mask, err = NewMaskV2(instr.Payload)
			if err != nil {
				return 0, err
			}
			trace("new mask: %s", instr.Payload)
			continue
		}

		var value int64
		{
			valueStr := instr.Payload
			var err error
			value, err = strconv.ParseInt(valueStr, 0, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid value %q on line %d: %w", valueStr, i+1, err)
			}
		}
		mask.SetValueAtAddresses(instr.Addr, value, mem)
	}

	var sum int64
	for _, v := range mem {
		sum += v
	}
	return sum, nil
}
