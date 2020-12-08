package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(r io.Reader) error {
	var (
		buf   = make([]byte, 1024)
		floor int
		err   error

		count         int
		basementIndex int
	)
	for err == nil {
		var n int
		n, err = r.Read(buf)
		for i, ch := range buf[:n] {
			switch ch {
			case '(':
				floor++
			case ')':
				floor--
				if floor < 0 && basementIndex == 0 {
					basementIndex = count + i + 1
				}
			}
		}
		count += n
	}
	if err != nil && err != io.EOF {
		return err
	}
	fmt.Printf("Floor %d\n", floor)
	fmt.Printf("Entered basement at char %d (1-indexed)\n", basementIndex)
	return nil
}
