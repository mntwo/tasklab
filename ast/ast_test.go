package ast

import (
	"fmt"
	"testing"
)

func TestAst(t *testing.T) {
	event := map[string]interface{}{
		"age":    20,
		"score":  50,
		"status": "active",
	}

	// t1 := `(and (== age 20) (>= score 51) (in status ("active" "pending")))`
	t2 := `(<= 1.1 1.1)`
	expr, err := ParseExpression(t2)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	result, err := Evaluate(expr, event)
	if err != nil {
		fmt.Println("Evaluation error:", err)
		return
	}

	fmt.Println("Result:", result) // output: true
}
