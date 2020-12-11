package main

import (
	"fmt"
	"testing"
)

func TestOccupiedAroundSeat(t *testing.T) {
	testInput := []string{
		"######",
		"#.L..#",
		"#LL#LL",
		"..L.L#",
	}
	tt := []struct {
		x, y int
		want int
	}{
		{x: 0, y: 0, want: 2},
		{x: 1, y: 0, want: 3},
		{x: 1, y: 2, want: 2},
		{x: 4, y: 3, want: 2},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("x:%d y:%d", tc.x, tc.y), func(t *testing.T) {
			if got := occupiedCountAroundSeat(tc.x, tc.y, testInput); got != tc.want {
				t.Fatalf("expected %d occupied for %d, %d, but got %d", tc.want, tc.x, tc.y, got)
			}
		})
	}
}

type coordinates struct {
	x, y      int
	expectErr bool
}

var grids = []struct {
	pattern  []string
	x, y     int
	occupied int
	closest  map[int]map[int]coordinates
}{
	{
		pattern: []string{
			".......#.",
			"...#.....",
			".#.......",
			".........",
			"..#L....#",
			"....#....",
			".........",
			"#........",
			"...#.....",
		},
		x:        3,
		y:        4,
		occupied: 8,
		closest: map[int]map[int]coordinates{
			-1: {
				-1: {x: 1, y: 2},
				0:  {x: 2, y: 4},
				1:  {x: 0, y: 7},
			},
			0: {
				-1: {x: 3, y: 1},
				1:  {x: 3, y: 8},
			},
			1: {
				-1: {x: 7, y: 0},
				0:  {x: 8, y: 4},
				1:  {x: 4, y: 5},
			},
		},
	},
	{
		pattern: []string{
			".............",
			".L.L.#.#.#.#.",
			".............",
		},
		x:        0,
		y:        1,
		occupied: 0,
		closest: map[int]map[int]coordinates{
			-1: {
				-1: {expectErr: true},
				0:  {expectErr: true},
				1:  {expectErr: true},
			},
			0: {
				-1: {expectErr: true},
				1:  {expectErr: true},
			},
			1: {
				-1: {expectErr: true},
				0:  {x: 1, y: 1, expectErr: false},
				1:  {expectErr: true},
			},
		},
	},
	{
		pattern: []string{
			".##.##.",
			"#.#.#.#",
			"##...##",
			"...L...",
			"##...##",
			"#.#.#.#",
			".##.##.",
		},
		x:        3,
		y:        3,
		occupied: 0,
		closest: map[int]map[int]coordinates{
			-1: {
				-1: {expectErr: true},
				0:  {expectErr: true},
				1:  {expectErr: true},
			},
			0: {
				-1: {expectErr: true},
				1:  {expectErr: true},
			},
			1: {
				-1: {expectErr: true},
				0:  {expectErr: true},
				1:  {expectErr: true},
			},
		},
	},
}

func TestOccupiedCountRays(t *testing.T) {
	for _, tc := range grids {
		t.Run(tc.pattern[0], func(t *testing.T) {
			if got := occupiedCountRays(tc.x, tc.y, tc.pattern); got != tc.occupied {
				t.Fatalf("expected %d seats occupied, but got %d", tc.occupied, got)
			}
		})
	}
}

func TestFindNearestSeatInDirection(t *testing.T) {
	g := grids[0]
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				// Can't search in no direction.
				continue
			}
			t.Run(fmt.Sprintf("dirX:%d dirY:%d", x, y), func(t *testing.T) {
				foundX, foundY, err := findNearestSeatInDirection(g.x, g.y, x, y, g.pattern)
				c := g.closest[x][y]
				if err != nil {
					if c.expectErr {
						return
					}
					t.Fatalf("unexpected error: %v", err)
				}
				if c.x != foundX || c.y != foundY {
					t.Fatalf(
						"expected closest coordinates (%d,%d), but got (%d,%d)",
						c.x, c.y,
						foundX, foundY)
				}
			})
		}
	}
}
