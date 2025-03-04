package calc

import (
	"database/sql"
	"fmt"
	"strconv"
)

type AST []Node

// Node represents a node in the AST.
type Node struct {
	SubnodesIdxs []int
	Operator     Operator
	Value        sql.Null[float64]
}

// BuildAST builds an AST ([]Node) the expression.
func (c *Calculator) BuildAST() (AST, error) {
	if len(c.tree) != 0 {
		return AST(c.tree), nil
	}

	for _, token := range c.tokens {
		switch {
		case token == "(":
			c.operators = append(c.operators, token)

		case token == ")":
			for len(c.operators) > 0 && c.operators[len(c.operators)-1] != "(" {
				err := c.buildNode()
				if err != nil {
					return nil, err
				}
			}

			c.operators = c.operators[:len(c.operators)-1]

		case IsOperator(token):
			for len(c.operators) > 0 && IsOperator(c.operators[len(c.operators)-1]) &&
				PrecedenceOf(c.operators[len(c.operators)-1]) >= PrecedenceOf(token) {
				err := c.buildNode()
				if err != nil {
					return nil, err
				}
			}

			c.operators = append(c.operators, token)

		default:
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: %s", ErrInvalidToken, token)
			}

			// Create a leaf task for the number
			valueNode := Node{
				SubnodesIdxs: nil,
				Operator:     "",
				Value:        sql.Null[float64]{V: value, Valid: true},
			}

			c.tree = append(c.tree, valueNode)
			c.valueStack = append(c.valueStack, len(c.tree)-1)
		}
	}

	// Process remaining c.operators
	for len(c.operators) > 0 {
		err := c.buildNode()
		if err != nil {
			return nil, err
		}
	}

	// The final result should be the only value left on the stack
	if len(c.valueStack) != 1 {
		return nil, ErrNotEnoughOperands
	}

	return AST(c.tree), nil
}

// buildNode creates a new node for an operator and links it to its operands.
func (c *Calculator) buildNode() error {
	if len(c.valueStack) < 2 {
		if len(c.valueStack) > 0 && c.operators[len(c.operators)-1] == "-" {
			c.tree[c.valueStack[len(c.stack)-1]].Value.V *= -1.0
			c.operators = c.operators[:len(c.operators)-1]

			return nil
		}

		return fmt.Errorf("%w for %s", ErrNotEnoughOperands, c.operators[len(c.operators)-1])
	}

	// Pop the operator
	operator := c.operators[len(c.operators)-1]
	c.operators = c.operators[:len(c.operators)-1]

	// Pop the two operands
	right := c.valueStack[len(c.valueStack)-1]
	left := c.valueStack[len(c.valueStack)-2]
	c.valueStack = c.valueStack[:len(c.valueStack)-2]

	// Create a new task for the operator
	node := Node{
		SubnodesIdxs: []int{left, right}, // Link to child tasks
		Operator:     Operator(operator),
		Value:        sql.Null[float64]{Valid: false, V: 0},
	}
	c.tree = append(c.tree, node)

	// Push the new task's ID onto the value stack
	c.valueStack = append(c.valueStack, len(c.tree)-1)

	return nil
}
