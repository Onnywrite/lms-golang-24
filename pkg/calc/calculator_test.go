package calc_test

import (
	"math"
	"testing"

	"github.com/Onnywrite/lms-golang-24/pkg/calc"
	"github.com/stretchr/testify/assert"
)

func TestCalc(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected float64
		err      error
	}{
		{"simple addition", "1 +	1", 2, nil},
		{"simple subtraction", "1.78-1", 0.78, nil},
		{"simple multiplication", "2*22.1", 44.2, nil},
		{"simple division", "4/2.25", 1.77778, nil},
		{"simple power", "2^3.5", 11.313708, nil},
		{"complex expression", "1+2*3-4/2^2", 6, nil},
		{"parentheses", "(1+2)*3", 9, nil},
		{"nested parentheses", "((1+2)*3)", 9, nil},
		{"very nested parentheses", "2.2*(52*(1+2^5)/26)", 145.2, nil},
		{"e pow pi", "e^p", math.Pow(math.E, math.Pi), nil},
		{"negative number", "-1*2", -2, nil},
		{"negative number in parentheses", "2+(-1)", 1, nil},
		{"negative number with parentheses", "(-1)", -1, nil},
		{"negative number with division", "-4/2", -2, nil},
		{"negative number with power", "-2^2", 4, nil},
		{"positive number with negative power", "-(2)^(-2)", -0.25, nil},
		{"negative number with negative power", "-2^(-2)", 0.25, nil},
		{"negative number with complex expression", "1+-2*3-4/2^2", -6, nil},
		{"empty", "", 0, calc.ErrEmptyExpression},
		{"division by zero", "1/0", 0, calc.ErrDivisionByZero},
		{"zero to negative power", "0^-1", 0, calc.ErrZeroBase},
		{"negative base to non-integer power", "-2^1.5", 0, calc.ErrNegativeBase},
		{"unclosed parentheses", "(1+2", 0, calc.ErrUnclosedParentheses},
		{"invalid token", "1+a", 0, calc.ErrInvalidToken},
		{"not enough operands", "+1", 0, calc.ErrNotEnoughOperands},
		{"unknown operator", "1&1", 0, calc.ErrInvalidToken},
	}

	t.Parallel()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calc.Calc(tc.input)
			if assert.ErrorIs(t, err, tc.err) {
				assert.InDelta(t, tc.expected, result, 0.00001)
			}
		})
	}
}
