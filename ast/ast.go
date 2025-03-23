package ast

import (
	"fmt"
	"strconv"
	"strings"
)

// expression is a list of expressions, where the first element is the operator, and the rest are operands.
// For example, the following expression:
// (and (== age 20) (>= score 51) (in status ("active" "pending")))
// is represented as:
// ["and", ["==", "age", 20], [">=", "score", 51], ["in", "status", ["active", "pending"]]]
// The following operators are supported:
// - and: logical AND
// - or: logical OR
// - not: logical NOT
// - >, <, >=, <=, ==, !=, <>: comparison operators
// - in: membership operator
type Expression interface{}

// parse S expression
func ParseExpression(s string) (Expression, error) {
	tokens := tokenize(s)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty expression")
	}
	return parseTokens(&tokens)
}

func tokenize(s string) []string {
	s = strings.ReplaceAll(s, "(", " ( ")
	s = strings.ReplaceAll(s, ")", " ) ")
	fields := strings.Fields(s)
	return fields
}

func parseTokens(tokens *[]string) (Expression, error) {
	if len(*tokens) == 0 {
		return nil, fmt.Errorf("unexpected end of input")
	}

	token := (*tokens)[0]
	*tokens = (*tokens)[1:]

	if token == "(" {
		var list []interface{}
		for len(*tokens) > 0 && (*tokens)[0] != ")" {
			subExpr, err := parseTokens(tokens)
			if err != nil {
				return nil, err
			}
			list = append(list, subExpr)
		}
		if len(*tokens) == 0 {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		*tokens = (*tokens)[1:] // remove ")"
		return list, nil
	} else if token == ")" {
		return nil, fmt.Errorf("unexpected )")
	} else {
		return tryConvertToken(token), nil
	}
}

func tryConvertToken(token string) interface{} {
	if i, err := strconv.Atoi(token); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(token, 64); err == nil {
		return f
	}
	if strings.HasPrefix(token, `"`) && strings.HasSuffix(token, `"`) {
		return token[1 : len(token)-1]
	}
	return token
}

// evaluate expression
func Evaluate(expr Expression, event map[string]interface{}) (bool, error) {
	switch e := expr.(type) {
	case []interface{}:
		if len(e) == 0 {
			return false, fmt.Errorf("empty expression")
		}
		op, ok := e[0].(string)
		if !ok {
			return false, fmt.Errorf("invalid operator: %v", e[0])
		}
		switch op {
		case "and":
			for _, subExpr := range e[1:] {
				result, err := Evaluate(subExpr, event)
				if err != nil {
					return false, err
				}
				if !result {
					return false, nil
				}
			}
			return true, nil
		case "or":
			for _, subExpr := range e[1:] {
				result, err := Evaluate(subExpr, event)
				if err != nil {
					return false, err
				}
				if result {
					return true, nil
				}
			}
			return false, nil
		case "not":
			if len(e) != 2 {
				return false, fmt.Errorf("invalid not expression: %v", e)
			}
			result, err := Evaluate(e[1], event)
			if err != nil {
				return false, err
			}
			return !result, nil
		case ">", "<", ">=", "<=", "==", "!=", "<>":
			if len(e) != 3 {
				return false, fmt.Errorf("invalid comparison expression: %v", e)
			}
			a, err := getValue(e[1], event)
			if err != nil {
				return false, err
			}
			b, err := getValue(e[2], event)
			if err != nil {
				return false, err
			}
			cmp, err := compare(a, b)
			if err != nil {
				return false, err
			}
			switch op {
			case ">":
				return cmp > 0, nil
			case "<":
				return cmp < 0, nil
			case ">=":
				return cmp >= 0, nil
			case "<=":
				return cmp <= 0, nil
			case "==":
				return cmp == 0, nil
			case "!=", "<>":
				return cmp != 0, nil
			}
		case "in":
			if len(e) != 3 {
				return false, fmt.Errorf("invalid in expression: %v", e)
			}
			value, err := getValue(e[1], event)
			if err != nil {
				return false, err
			}
			list, ok := e[2].([]interface{})
			if !ok {
				return false, fmt.Errorf("invalid list for in operator: %v", e[2])
			}
			for _, item := range list {
				cmp, err := compare(value, item)
				if err != nil {
					return false, err
				}
				if cmp == 0 {
					return true, nil
				}
			}
			return false, nil
		default:
			return false, fmt.Errorf("unknown operator: %s", op)
		}
	case string:
		value, err := getValue(e, event)
		if err != nil {
			return false, err
		}
		return value != nil, nil
	default:
		return false, fmt.Errorf("invalid expression: %v", e)
	}
	return false, nil
}

func getValue(token interface{}, event map[string]interface{}) (interface{}, error) {
	switch t := token.(type) {
	case string:
		if strings.HasPrefix(t, `"`) && strings.HasSuffix(t, `"`) {
			return t[1 : len(t)-1], nil
		}
		if value, ok := event[t]; ok {
			return value, nil
		}
		return t, nil
	default:
		return t, nil
	}
}

func compare(a, b interface{}) (float64, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return float64(a) - float64(b), nil
		case float64:
			return float64(a) - b, nil
		}
	case float64:
		switch b := b.(type) {
		case int:
			return a - float64(b), nil
		case float64:
			return a - b, nil
		}
	case string:
		if b, ok := b.(string); ok {
			return float64(strings.Compare(a, b)), nil
		}
	}
	return 0, fmt.Errorf("unsupported types: %T and %T", a, b)
}
