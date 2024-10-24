package calc

// TODO: negative numbers (2^(-1), 3*(-10), ...)
import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	// Syntax errors
	ErrUnclosedParentheses = errors.New("unclosed parentheses")
	ErrNotEnoughOperands   = errors.New("not enough operands")
	ErrUnknownOperator     = errors.New("unknown operator")
	ErrInvalidToken        = errors.New("invalid token")
	ErrEmptyExpression     = errors.New("empty expression")
	// Math errors
	ErrDivisionByZero = errors.New("division by zero")
	ErrZeroBase       = errors.New("zero to a non-positive exponent")
	ErrNegativeBase   = errors.New("negative base to a non-integer exponent")
)

func Calc(expression string) (float64, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}

	return parseExpression(tokens)
}

func tokenize(expression string) ([]string, error) {
	var (
		tokens []string
		number strings.Builder
	)

	for _, char := range expression {
		switch char {
		case ' ', '\t':
			continue
		case '+', '-', '*', '/', '^', '(', ')':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			tokens = append(tokens, string(char))
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

func validate(tokens []string) error {
	// Check length of tokens
	if len(tokens) == 0 {
		return ErrEmptyExpression
	}

	// Check parentheses
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

func parseExpression(tokens []string) (float64, error) {
	stack := make([]float64, 0, len(tokens))
	operators := make([]string, 0, len(tokens))

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"^": 3,
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				if err := applyOperator(&stack, &operators); err != nil {
					return 0, err
				}
			}
			operators = operators[:len(operators)-1]
		} else if _, ok := precedence[token]; ok {
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[token] {
				if err := applyOperator(&stack, &operators); err != nil {
					return 0, err
				}
			}
			operators = append(operators, token)
		} else {
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("%w: %s", ErrInvalidToken, token)
			}
			stack = append(stack, value)
		}
	}

	for len(operators) > 0 {
		if err := applyOperator(&stack, &operators); err != nil {
			return 0, err
		}
	}

	return stack[0], nil
}

func applyOperator(stack *[]float64, operators *[]string) error {
	if len(*stack) < 2 {
		// If the last operator is a minus and we don't have
		// enough operands, we can assume that the last operator is a negative number
		// and we multiply the last operand by -1.
		if len(*stack) > 0 && (*operators)[len(*operators)-1] == "-" {
			(*stack)[len(*stack)-1] *= -1
			*operators = (*operators)[:len(*operators)-1]
			return nil
		}

		// But if the last operand in not minus, we lack of operands.
		return fmt.Errorf("%w for %s", ErrNotEnoughOperands, (*operators)[len(*operators)-1])
	}

	b := (*stack)[len(*stack)-1]
	a := (*stack)[len(*stack)-2]
	*stack = (*stack)[:len(*stack)-2]

	operator := (*operators)[len(*operators)-1]
	*operators = (*operators)[:len(*operators)-1]

	var result float64
	switch operator {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return fmt.Errorf("%w: %f/%f", ErrDivisionByZero, a, b)
		}

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

		result = math.Pow(a, b)
	default:
		return fmt.Errorf("%w: %s", ErrUnknownOperator, operator)
	}

	*stack = append(*stack, result)

	return nil
}
