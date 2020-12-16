package main

import (
	"fmt"
	"strings"
	"testing"
)

func makeMaskStr(forceZero, forceOne []int) string {
	const bits = 36
	var b strings.Builder
	contains := func(i []int, target int) bool {
		for _, n := range i {
			if n == target {
				return true
			}
		}
		return false
	}
	for i := bits; i > 0; i-- {
		switch {
		case contains(forceZero, i):
			b.WriteRune('0')
		case contains(forceOne, i):
			b.WriteRune('1')
		default:
			b.WriteRune('X')
		}
	}
	return b.String()
}

func TestMakeMaskStr(t *testing.T) {
	var (
		maskStr = makeMaskStr([]int{7}, []int{6})
		want    = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX01XXXXX"
	)
	if got := maskStr; got != want {
		t.Fatalf("wanted: %s, but got %s", want, got)
	}
}

func TestNewMask(t *testing.T) {
	tt := []struct {
		maskStr string
		in      []int64
		want    []int64
	}{
		{
			maskStr: makeMaskStr([]int{11, 1}, []int{4, 13}),
			in:      []int64{1024, 1, 0},
			want:    []int64{4104, 4104, 4104},
		},
		{
			maskStr: makeMaskStr([]int{7}, []int{6}),
			in:      []int64{64, 0, 128},
			want:    []int64{32, 32, 160},
		},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {
			t.Logf("maskStr: %s (len: %d)", tc.maskStr, len(tc.maskStr))
			m, err := NewMaskV1(tc.maskStr)
			t.Logf("And: %036b (%d)", m.And, m.And)
			t.Logf("Or:  %036b (%d)", m.Or, m.Or)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for i := range tc.in {
				t.Run(fmt.Sprintf("%d->%d", tc.in[i], tc.want[i]), func(t *testing.T) {
					got := m.Apply(tc.in[i])
					if want := tc.want[i]; got != want {
						t.Errorf("expected %d, but got %d", want, got)
					}
				})
			}
		})
	}
}
