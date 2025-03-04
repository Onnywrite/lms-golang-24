package calc_test

import (
	"database/sql"
	"testing"

	"github.com/Onnywrite/lms-golang-24/pkg/calc"
	"github.com/stretchr/testify/require"
)

func TestBuildAST(t *testing.T) {
	expr, err := calc.NewCalculator("2.2 * (25 ^ 0.5 * (1 + 2 ^ (-5)) / 26)")
	require.NoError(t, err)

	tree, err := expr.BuildAST()
	require.NoError(t, err)

	expectedTree := []calc.Node{
		{nil, "", val(2.2)},
		{nil, "", val(25)},
		{nil, "", val(0.5)},
		{[]int{1, 2}, calc.OpPower, null()},
		{nil, "", val(1)},
		{nil, "", val(2)},
		{nil, "", val(-5)},
		{[]int{5, 6}, calc.OpPower, null()},
		{[]int{4, 7}, calc.OpAdd, null()},
		{[]int{3, 8}, calc.OpMultiply, null()},
		{nil, "", val(26)},
		{[]int{9, 10}, calc.OpDivide, null()},
		{[]int{0, 11}, calc.OpMultiply, null()},
	}

	require.Equal(t, expectedTree, []calc.Node(tree))
}

func val(v float64) sql.Null[float64] {
	return sql.Null[float64]{
		V:     v,
		Valid: true,
	}
}

func null() sql.Null[float64] {
	return sql.Null[float64]{
		V:     0,
		Valid: false,
	}
}
