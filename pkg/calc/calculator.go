package calc

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	// Syntax errors.
	ErrUnclosedParentheses = errors.New("unclosed parentheses")
	ErrNotEnoughOperands   = errors.New("not enough operands")
	ErrUnknownOperator     = errors.New("unknown operator")
	ErrInvalidToken        = errors.New("invalid token")
	ErrEmptyExpression     = errors.New("empty expression")
	// Math errors.
	ErrDivisionByZero = errors.New("division by zero")
	ErrZeroBase       = errors.New("zero to a non-positive exponent")
	ErrNegativeBase   = errors.New("negative base to a non-integer exponent")
)

// Operator represents a mathematical operator.
type Operator string

const (
	OpAdd      Operator = "+"
	OpSubtract Operator = "-"
	OpMultiply Operator = "*"
	OpDivide   Operator = "/"
	OpPower    Operator = "^"
)

// Calculator represents a stateful calculator for long-running operations.
type Calculator struct {
	stack      []float64
	operators  []string
	tokens     []string
	valueStack []int
	tree       []Node
}

// NewCalculator initializes a new Calculator instance.
func NewCalculator(expression string) (*Calculator, error) {
	tokens, err := Tokenize(expression)
	if err != nil {
		return nil, err
	}

	return &Calculator{
		stack:      make([]float64, 0, len(tokens)),
		operators:  make([]string, 0, len(tokens)),
		tokens:     tokens,
		valueStack: make([]int, 0, len(tokens)),
		tree:       make([]Node, 0, len(tokens)),
	}, nil
}

// Calculate processes the tokens and computes the result.
func (c *Calculator) Calculate() (float64, error) {
	for _, token := range c.tokens {
		switch {
		case token == "(":
			c.operators = append(c.operators, token)

		case token == ")":
			for len(c.operators) > 0 && c.operators[len(c.operators)-1] != "(" {
				err := c.applyOperator()
				if err != nil {
					return 0, err
				}
			}

			c.operators = c.operators[:len(c.operators)-1]

		case PrecedenceOf(token) != 0:
			for len(c.operators) > 0 && PrecedenceOf(c.operators[len(c.operators)-1]) >= PrecedenceOf(token) {
				err := c.applyOperator()
				if err != nil {
					return 0, err
				}
			}

			c.operators = append(c.operators, token)

		default:
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("%w: %s", ErrInvalidToken, token)
			}

			c.stack = append(c.stack, value)
		}
	}

	for len(c.operators) > 0 {
		err := c.applyOperator()
		if err != nil {
			return 0, err
		}
	}

	if len(c.stack) == 0 {
		return 0, ErrEmptyExpression
	}

	return c.stack[0], nil
}

// applyOperator applies the top operator to the operands on the stack.
func (c *Calculator) applyOperator() error {
	if len(c.stack) < 2 {
		if len(c.stack) > 0 && c.operators[len(c.operators)-1] == "-" {
			c.stack[len(c.stack)-1] *= -1
			c.operators = c.operators[:len(c.operators)-1]

			return nil
		}

		// But if the last operand in not minus, we lack of operands.
		return fmt.Errorf("%w for %s", ErrNotEnoughOperands, c.operators[len(c.operators)-1])
	}

	b := c.stack[len(c.stack)-1]
	a := c.stack[len(c.stack)-2]
	c.stack = c.stack[:len(c.stack)-2]

	operator := c.operators[len(c.operators)-1]
	c.operators = c.operators[:len(c.operators)-1]

	var result float64

	switch operator {
	case "+":
		// TODO: add sleep
		result = a + b

	case "-":
		// TODO: add sleep
		result = a - b

	case "*":
		// TODO: add sleep
		result = a * b

	case "/":
		if b == 0 {
			return fmt.Errorf("%w: %f/%f", ErrDivisionByZero, a, b)
		}

		// TODO: add sleep
		result = a / b

	case "^":
		// We can't calc 0^I if I <= 0, because it's undefined.
		if a == 0 && b <= 0 {
			return fmt.Errorf("%w: %f^%f", ErrZeroBase, a, b)
		}

		// We can't raise a negative number to a non-integer power.
		//
		// Here is why:
		//  (-2)^(1.5) = (-2)^(3/2) = sqrt((-2)^3) = sqrt(-8)
		// the result is a complex number, which my calculator doesn't support.
		if a < 0 && math.Trunc(b) != b {
			return fmt.Errorf("%w: %f^%f", ErrNegativeBase, a, b)
		}

		// TODO: add sleep
		result = math.Pow(a, b)

	default:
		return fmt.Errorf("%w: %s", ErrUnknownOperator, operator)
	}

	c.stack = append(c.stack, result)

	return nil
}

// Calculate is a wrapper of [Calculator.Calculate] method.
func Calculate(expr string) (float64, error) {
	c, err := NewCalculator(expr)
	if err != nil {
		return 0, err
	}

	return c.Calculate()
}

// PrecedenceOf returns the precedence level of an operator.
func PrecedenceOf(op string) int {
	const (
		plusPredence = 1
		multPredence = 2
		powPredence  = 3
	)

	switch op {
	case "+", "-":
		return plusPredence

	case "*", "/":
		return multPredence

	case "^":
		return powPredence
	}

	return 0
}

func Tokenize(expression string) ([]string, error) {
	number := strings.Builder{}
	tokens := make([]string, 0, 16)

	for i, char := range expression {
		switch char {
		case ' ', '\t':
			continue

		case '+', '*', '/', '^', '(', ')':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}

			tokens = append(tokens, string(char))

		case '-':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}

			if i == 0 || expression[i-1] == '(' || IsOperator(string(expression[i-1])) {
				_, _ = number.WriteRune(char)
			} else {
				tokens = append(tokens, string(char))
			}

		case 'e':
			number.WriteString(strconv.FormatFloat(math.E, 'f', -1, 64))

		case 'p':
			number.WriteString(strconv.FormatFloat(math.Pi, 'f', -1, 64))

		default:
			number.WriteRune(char)
		}
	}

	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	if err := validate(tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func IsOperator(token string) bool {
	return token == string(OpAdd) ||
		token == string(OpSubtract) ||
		token == string(OpMultiply) ||
		token == string(OpDivide) ||
		token == string(OpPower)
}

func validate(tokens []string) error {
	if len(tokens) == 0 {
		return ErrEmptyExpression
	}

	openParentheses := 0

	for _, token := range tokens {
		if token == "(" {
			openParentheses++
		} else if token == ")" {
			openParentheses--
		}
	}

	if openParentheses != 0 {
		return ErrUnclosedParentheses
	}

	return nil
}
