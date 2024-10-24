package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Onnywrite/lms-golang-24/pkg/calc"
)

func main() {
	fmt.Print("Enter an expression: ")

	buf := bufio.NewReader(os.Stdin)

	expression, _ := buf.ReadString('\n')
	expression = strings.TrimSpace(expression)

	result, err := calc.Calc(expression)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Result: %v\n", result)
}
