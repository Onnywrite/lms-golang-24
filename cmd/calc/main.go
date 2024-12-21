package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Onnywrite/lms-golang-24/pkg/calc"
)

const verySusExitCode = 52

func main() {
	_, _ = fmt.Print("Enter an expression: ") //nolint: forbidigo

	buf := bufio.NewReader(os.Stdin)

	expression, _ := buf.ReadString('\n')
	expression = strings.TrimSpace(expression)

	result, err := calc.Calc(expression)
	if err != nil {
		_, _ = fmt.Println(err) //nolint: forbidigo

		os.Exit(verySusExitCode)
	}

	_, _ = fmt.Printf("Result: %v\n", result) //nolint: forbidigo
}
