package calc

// TODO:
//  1. Validation
//  2. Negative numbers (2^(-1), 3*(-10), ...)
import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	ErrUnclosedParentheses = errors.New("unclosed parentheses")
	ErrNoOperand           = errors.New("no operand for operator")
	ErrInvalidOperator     = errors.New("invalid operator")
	ErrInvalidToken        = errors.New("invalid token")
	ErrEmptyExpression     = errors.New("empty expression")
)

func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	return parseExpression(tokens)
}

func tokenize(expression string) []string {
	var tokens []string
	var number strings.Builder
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

	return tokens
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
	var stack []float64
	var operators []string

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
				stack = applyOperator(stack, &operators)
			}
			operators = operators[:len(operators)-1]
		} else if _, ok := precedence[token]; ok {
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[token] {
				stack = applyOperator(stack, &operators)
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
		stack = applyOperator(stack, &operators)
	}

	return stack[0], nil
}

func applyOperator(stack []float64, operators *[]string) []float64 {
	if len(stack) < 2 {
		return stack
	}

	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]

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
		result = a / b
	case "^":
		result = math.Pow(a, b)
	}

	stack = append(stack, result)
	return stack
}
