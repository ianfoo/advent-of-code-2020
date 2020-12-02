# Advent of Code 2020

This repository tracks puzzles for [Advent of Code 2020](https://adventofcode.com/2020). Go play!

If you're really ambitious, try doing each puzzle [in a different 
language](https://www.benkraft.org/2017/12/26/advent-of-code/).

## Running

For solutions implemented in Go, which is all of them at the time of this
writing, you can run the program by changing to the directory of the day
you wish to run and running
 
```
go run . < input.txt
```
unless I've accidentally named the input file something else.

Alternatively, there's a naive Makefile that, as of this writing, assumes all
exercises are completed in Go and attempts to run them all. Just run `make`
in the project root.

## Caveats

This is slapdash code, with only as much effort put into it as required to
get the right answer according to Advent of Code. There are no tests unless
debugging is required. Advent of Code is as much an exercise for me to not
get hung up on always making pristine production quality code when the
situation does not call for it, instead optimizing for speed of
implementation here, without sacrificing too much readiability.