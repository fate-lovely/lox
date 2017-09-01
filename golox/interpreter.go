package main

import "fmt"

type Val interface{}
type Number float64

type RuntimeError struct {
	token *Token
	msg   string
}

func NewRuntimeError(token *Token, msg string) *RuntimeError {
	return &RuntimeError{token, msg}
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintf("line %d, %s", re.token.line, re.msg)
}

/*----------  Literal  ----------*/

func (expr *ExprLiteral) Eval() Val {
	return expr.value
}

/*----------  Unary  ----------*/

func (expr *ExprUnary) Eval() Val {
	value := expr.operand.Eval()
	switch expr.operator.typ {
	case BANG:
		return !getTruthy(value)
	case MINUS:
		return -(value.(float64))
	}

	// unreachable
	panic("should neven reach here")
}

/*----------  Binary  ----------*/
func (expr *ExprBinary) Eval() Val {
	left := expr.left.Eval()
	right := expr.right.Eval()

	checkNumberOperands := func() {
		if isNumber(left) && isNumber(right) {
			return
		}
		panic(NewRuntimeError(expr.operator, "operands must be numbers"))
	}

	switch expr.operator.typ {
	case PLUS:
		if isNumber(left) && isNumber(right) {
			return toNumber(left) + toNumber(right)
		}
		if isString(left) && isString(right) {
			return toString(left) + toString(right)
		}
		panic(NewRuntimeError(expr.operator, "operands must be two numbers or two strings"))
	case MINUS:
		checkNumberOperands()
		return toNumber(left) - toNumber(right)
	case SLASH:
		checkNumberOperands()
		// catch divide by zero
		r := toNumber(right)
		if r == 0 {
			panic(NewRuntimeError(expr.operator, "divide by zero"))
		}
		return toNumber(left) / r
	case STAR:
		checkNumberOperands()
		return toNumber(left) * toNumber(right)
	case GREATER:
		checkNumberOperands()
		return toNumber(left) > toNumber(right)
	case GREATER_EQUAL:
		checkNumberOperands()
		return toNumber(left) >= toNumber(right)
	case LESS:
		checkNumberOperands()
		return toNumber(left) < toNumber(right)
	case LESS_EQUAL:
		checkNumberOperands()
		return toNumber(left) <= toNumber(right)
	case EQUAL_EQUAL:
		return left == right
	case BANG_EQUAL:
		return left != right
	}

	// unreachable
	panic("should neven reach here")
}

/*----------  Grouping  ----------*/

func (expr *ExprGrouping) Eval() Val {
	return expr.operand.Eval()
}

/*----------  Helper Methods  ----------*/

// `false` and `nil` is false
// everything else is true
func getTruthy(val Val) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return true
}

func isNumber(val Val) bool {
	_, ok := val.(Number)
	return ok
}

func isString(val Val) bool {
	_, ok := val.(string)
	return ok
}

func toNumber(val Val) Number {
	if n, ok := val.(Number); ok {
		return n
	}
	// should never happen
	panic("toNumber should always be called with a number")
}

func toString(val Val) string {
	if s, ok := val.(string); ok {
		return s
	}
	// should never happen
	panic("toString should always be called with a string")
}