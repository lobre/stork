package stork

import "testing"

func TestCalculateLeash(t *testing.T) {
	lp := leashParams{0, 400, 0, 40}

	tests := map[string]struct {
		input, want int
	}{
		"50":   {50, 5},
		"-50":  {-50, 0},
		"2000": {2000, 40},
		"400":  {400, 40},
	}

	for name, tc := range tests {
		got := calculateLeash(lp, tc.input)
		if got != tc.want {
			t.Fatalf("%s: expected %d, got %d", name, tc.want, got)
		}
	}

}
