// This is a program that will compute the answers to both parts of the Day 8
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/8)
package main

import (
	"bufio"
	"errors"
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
		log.Printf(format, params...)
	}
}

func run(r io.Reader) error {
	input, err := readInput(r)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	instr := make([]Instruction, 0, len(input))
	for i, line := range input {
		current, err := parseInstruction(line)
		if err != nil {
			return err
		}
		instr = append(instr, current)
		trace("[%4d] instruction: %v", i, current)
	}

	{
		result, err := part1(instr)
		if err != nil {
			return fmt.Errorf("part 1: %w", err)
		}
		fmt.Printf("Part 1: %d\n", result)
	}

	{
		result, err := part2(instr)
		if err != nil {
			return fmt.Errorf("part 2: %w", err)
		}
		fmt.Printf("Part 2: %d\n", result)
	}

	return nil
}

func part1(instr []Instruction) (int, error) {
	acc, err := executeInstructions(instr)
	return acc, err
}

func part2(instr []Instruction) (int, error) {
	return executeAndModifyInstructions(instr)
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
	Op  string
	Arg int
}

func parseInstruction(instrStr string) (Instruction, error) {
	tokens := strings.Split(instrStr, " ")
	if len(tokens) != 2 {
		return Instruction{}, fmt.Errorf("invalid instruction %q", instrStr)
	}

	op := tokens[0]
	arg, err := strconv.Atoi(tokens[1])
	if err != nil {
		return Instruction{}, fmt.Errorf("invalid instruction argument %q: %w", instrStr, err)
	}

	instr := Instruction{
		Op:  op,
		Arg: arg,
	}
	return instr, nil
}

func exec(instr Instruction, prgCtr, acc int) (int, int) {
	trace("running instruction %v", instr)
	switch instr.Op {
	case "nop":
		prgCtr++
	case "acc":
		acc += instr.Arg
		prgCtr++
	case "jmp":
		prgCtr += instr.Arg
	}
	return prgCtr, acc
}

func executeInstructions(instr []Instruction) (int, error) {
	var (
		acc        = 0
		prgCtr     = 0
		runHistory = make([]bool, len(instr))
	)
	for !runHistory[prgCtr] {
		inst := instr[prgCtr]
		runHistory[prgCtr] = true
		prgCtr, acc = exec(inst, prgCtr, acc)
	}
	return acc, nil
}

func executeAndModifyInstructions(instr []Instruction) (int, error) {
	var (
		acc    = 0
		prgCtr = 0

		runHistory    = make([]bool, len(instr))
		revisionPoint = 0
	)

	findAndChangeNextCandidate := func(index int) ([]Instruction, int) {
		workingCopy := make([]Instruction, len(instr))
		copy(workingCopy, instr)

		for ; index < len(workingCopy); index++ {
			if next := workingCopy[index]; next.Op == "jmp" || next.Op == "nop" {
				trace("found candidate %v at index %d", next, index)
				break
			}
		}

		// We have reached the end of the instructions.
		if index == len(workingCopy) {
			return workingCopy, index
		}

		if workingCopy[index].Op == "nop" {
			workingCopy[index].Op = "jmp"
		} else if workingCopy[index].Op == "jmp" {
			workingCopy[index].Op = "nop"
		}

		// Move the index past the point that was just changed, so we don't
		// repeat this forever.
		index++

		trace("%v", workingCopy)
		return workingCopy, index
	}

	resetProgram := func() {
		prgCtr = 0
		acc = 0
		runHistory = make([]bool, len(instr))
	}

	var workingCopy = instr
	for prgCtr != len(instr) {
		if runHistory[prgCtr] {
			trace("cycle dectected on instruction %d: %v", prgCtr, workingCopy[prgCtr])

			if revisionPoint >= len(instr) {
				return 0, errors.New("no more candidate instructions to try")
			}
			workingCopy, revisionPoint = findAndChangeNextCandidate(revisionPoint)

			resetProgram()
			continue
		}

		inst := workingCopy[prgCtr]
		runHistory[prgCtr] = true

		prgCtr, acc = exec(inst, prgCtr, acc)
	}
	return acc, nil
}
