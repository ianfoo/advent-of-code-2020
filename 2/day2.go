// This is a program that will compute the answers to both parts of the Day 2
// Advent of Code 2020 puzzle. (adventofcode.com/2020/day/2)
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
	}
}

// Match password rule and password lines:
// [min]-[max] [char]: [password]
var ruleAndPasswordRegexp = regexp.MustCompile(`(?P<num1>\d+)-(?P<num2>\d+) (?P<char>\w): (?P<password>\w+)$`)

func run(r io.Reader) error {
	numOldValidPasswords, numNewValidPasswords, err := countValidPasswords(r)
	if err != nil {
		return err
	}
	fmt.Printf("Number of valid passwords by old rules: %d\n", numOldValidPasswords)
	fmt.Printf("Number of valid passwords by new rules: %d\n", numNewValidPasswords)
	return nil
}

// PasswordEntry represents a line in the input
type PasswordEntry struct {
	num1, num2 int
	char       string
	password   string
}

func countValidPasswords(r io.Reader) (int, int, error) {
	var (
		s                            = bufio.NewScanner(r)
		validCountOld, validCountNew int
	)
	for s.Scan() {
		t := s.Text()
		entry, err := getPasswordEntryForLine(t)
		if err != nil {
			return 0, 0, fmt.Errorf("extract password entry for line: %w", err)
		}

		if isValidForOldRules(entry) {
			validCountOld++
		}
		if isValidForNewRules(entry) {
			validCountNew++
		}
	}
	if err := s.Err(); err != nil {
		return 0, 0, err
	}
	return validCountOld, validCountNew, nil
}

func getPasswordEntryForLine(line string) (PasswordEntry, error) {
	m := ruleAndPasswordRegexp.FindStringSubmatch(line)
	if len(m) < 5 {
		return PasswordEntry{}, fmt.Errorf("invalid input: %q", line)
	}

	num1Str, num2Str, char, password := m[1], m[2], m[3], m[4]
	num1, err := strconv.Atoi(num1Str)
	if err != nil {
		return PasswordEntry{}, fmt.Errorf("invalid first number value: %w", err)
	}
	num2, err := strconv.Atoi(num2Str)
	if err != nil {
		return PasswordEntry{}, fmt.Errorf("invalid second number value: %w", err)
	}

	entry := PasswordEntry{
		num1:     num1,
		num2:     num2,
		char:     char,
		password: password,
	}
	return entry, nil
}

// isValidForOldRules determines password validity according to old rules,
// where the password must contain between num1 and num2 instances of the
// specified char.
func isValidForOldRules(pe PasswordEntry) bool {
	count := strings.Count(pe.password, pe.char)
	isValid := count >= pe.num1 && count <= pe.num2
	return isValid
}

// isValidForNewRules determines password validity according to new rules, where
// the specified character must be in exactly one of the two character locations
// provided.
func isValidForNewRules(pe PasswordEntry) bool {
	// Indices are 1-based in password file, so correct here!
	index1, index2 := pe.num1 -1, pe.num2 - 1

	inPos1 := pe.password[index1] == pe.char[0]
	inPos2 := pe.password[index2] == pe.char[0]

	// Char must be in either position 1 or 2, but not both.
	isValid := (inPos1 || inPos2) && !(inPos1 && inPos2)
	return isValid
}
