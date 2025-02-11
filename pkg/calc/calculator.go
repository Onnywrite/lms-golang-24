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

			if i == 0 || expression[i-1] == '(' || isOperator(string(expression[i-1])) {
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

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/" || token == "^"
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

const (
	plusPredence  = 1
	minusPredence = 1
	multPredence  = 2
	divPredence   = 2
	powPredence   = 3
)

var precedence = map[string]int{
	"+": plusPredence,
	"-": minusPredence,
	"*": multPredence,
	"/": divPredence,
	"^": powPredence,
}

func parseExpression(tokens []string) (float64, error) {
	stack := make([]float64, 0, len(tokens))
	operators := make([]string, 0, len(tokens))

	for i := range tokens {
		token := tokens[i]

		_, isOperator := precedence[token]

		switch {
		case token == "(":
			operators = append(operators, token)

		case token == ")":
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				if err := applyOperator(&stack, &operators); err != nil {
					return 0, err
				}
			}

			operators = operators[:len(operators)-1]

		case isOperator:
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[token] {
				if err := applyOperator(&stack, &operators); err != nil {
					return 0, err
				}
			}

			operators = append(operators, token)

		default:
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
