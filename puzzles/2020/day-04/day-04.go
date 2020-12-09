// This is a program that will compute the answers to both parts of the Day 4
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/4)
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if err := run(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(r io.Reader) error {
	input, err := readInput(r)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	passports := buildPassports(input)

	{
		result, err := part1(passports)
		if err != nil {
			return fmt.Errorf("part 1: %w", err)
		}
		fmt.Printf("Part 1: %d\n", result)
	}

	{
		result, err := part2(passports)
		if err != nil {
			return fmt.Errorf("part 2: %w", err)
		}
		fmt.Printf("Part 2: %d\n", result)
	}

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

func buildPassports(input []string) []string {
	var (
		current   string
		passports []string
	)
	for _, line := range input {
		if line != "" {
			current += " " + line
			continue
		}
		current = strings.Trim(current, " ")
		passports = append(passports, current)
		current = ""
	}
	if current != "" {
		current = strings.Trim(current, " ")
		passports = append(passports, current)
	}
	return passports
}

func part1(input []string) (int, error) {
	var (
		requiredFields = []string{
			"byr",
			"iyr",
			"eyr",
			"hgt",
			"hcl",
			"ecl",
			"pid",
		}
		numValid int
	)

	for _, currentPassport := range input {
		isValid := true
		for _, tok := range requiredFields {
			if !strings.Contains(currentPassport, tok+":") {
				isValid = false
				break
			}
		}
		if isValid {
			numValid++
		}
	}
	return numValid, nil
}

func part2(input []string) (int, error) {
	var (
		validators = map[string]func(string) bool{
			"byr": validateByr,
			"iyr": validateIyr,
			"eyr": validateEyr,
			"hgt": validateHgt,
			"hcl": validateHcl,
			"ecl": validateEcl,
			"pid": validatePid,
		}

		numValid int
	)

	for _, currentPassport := range input {
		var (
			numValidFields = 0
			fields         = strings.Split(currentPassport, " ")
		)
		for _, field := range fields {
			splitField := strings.SplitN(field, ":", 2)
			var (
				fieldName  = splitField[0]
				fieldValue = splitField[1]
			)
			validator, ok := validators[fieldName]
			if !ok {
				// Ignore unknown fields.
				continue
			}
			if validField := validator(fieldValue); !validField {
				break
			}
			numValidFields++
		}

		if numValidFields == len(validators) {
			// If we've made it here, all required fields have been validated.
			numValid++
		}
	}
	return numValid, nil
}

func validateByr(s string) bool {
	byr, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return byr >= 1920 && byr <= 2002
}

func validateIyr(s string) bool {
	iyr, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return iyr >= 2010 && iyr <= 2020
}

func validateEyr(s string) bool {
	eyr, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return eyr >= 2020 && eyr <= 2030
}

func validateHgt(s string) bool {
	if strings.HasSuffix(s, "cm") {
		s = strings.TrimSuffix(s, "cm")
		hgt, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return hgt >= 150 && hgt <= 193
	}
	if strings.HasSuffix(s, "in") {
		s = strings.TrimSuffix(s, "in")
		hgt, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return hgt >= 59 && hgt <= 76
	}
	return false
}

func validateHcl(s string) bool {
	match, err := regexp.MatchString(`^#[0-9a-f]{6}$`, s)
	if err != nil {
		panic("bad regexp in validateHcl")
	}
	return match
}

func validateEcl(s string) bool {
	match, err := regexp.MatchString(`^(?:amb)|(?:blu)|(?:brn)|(?:gry)|(?:grn)|(?:hzl)|(?:oth)$`, s)
	if err != nil {
		panic("bad regexp in validateEcl")
	}
	return match
}

func validatePid(s string) bool {
	match, err := regexp.MatchString(`^\d{9}$`, s)
	if err != nil {
		panic("bad regexp in validatePid")
	}
	return match
}
