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

// we use exception as control flow
type FunctionReturn struct {
	value Val
}

func NewFunctionReturn(value Val) *FunctionReturn {
	return &FunctionReturn{value}
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintf("line %d, %s", re.token.line, re.msg)
}

/*----------  Stmt: Print  ----------*/

func (s *StmtPrint) Run(env *Env) {
	val := s.expr.Eval(env)
	fmt.Println(val)
}

/*----------  Stmt: Expression  ----------*/

func (s *StmtExpression) Run(env *Env) {
	s.expr.Eval(env)
}

/*----------  Stmt: Variable Declaration  ----------*/

func (s *StmtVarDecl) Run(env *Env) {
	var val Val
	if s.value != nil {
		val = s.value.Eval(env)
	}
	env.Define(s.name.lexeme, val)
}

/*----------  Stmt: Block  ----------*/

func (s *StmtBlock) Run(env *Env) {
	newEnv := NewEnv(env)
	for _, stmt := range s.stmts {
		stmt.Run(newEnv)
	}
}

/*----------  Stmt: If  ----------*/

func (s *StmtIf) Run(env *Env) {
	val := s.condition.Eval(env)
	if getTruthy(val) {
		s.trueBranch.Run(env)
	} else {
		if s.falseBranch != nil {
			s.falseBranch.Run(env)
		}
	}
}

/*----------  Stmt: While  ----------*/

func (s *StmtWhile) Run(env *Env) {
	for getTruthy(s.condition.Eval(env)) {
		s.body.Run(env)
	}
}

/*----------  Stmt: Function Declaration  ----------*/

func (s *StmtFuncDecl) Run(env *Env) {
	s.closure = env
	env.Define(s.name.lexeme, s)
}

/*----------  Stmt: Return  ----------*/

func (s *StmtReturn) Run(env *Env) {
	var value Val
	if s.value != nil {
		value = s.value.Eval(env)
	}
	panic(NewFunctionReturn(value))
}

/*----------  Expr: Assignment  ----------*/

func (expr *ExprAssignment) Eval(env *Env) Val {
	val := expr.val.Eval(env)
	env.Set(expr.name, val)
	return val
}

/*----------  Expr: Literal  ----------*/

func (expr *ExprLiteral) Eval(env *Env) Val {
	return expr.value
}

/*----------  Expr: Unary  ----------*/

func (expr *ExprUnary) Eval(env *Env) Val {
	value := expr.operand.Eval(env)
	switch expr.operator.typ {
	case BANG:
		return !getTruthy(value)
	case MINUS:
		return -(value.(Number))
	}

	// unreachable
	panic("should neven reach here")
}

/*----------  Expr: Binary  ----------*/
func (expr *ExprBinary) Eval(env *Env) Val {
	left := expr.left.Eval(env)
	right := expr.right.Eval(env)

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

/*----------  Expr: Grouping  ----------*/

func (expr *ExprGrouping) Eval(env *Env) Val {
	return expr.operand.Eval(env)
}

/*----------  Expr: Variable  ----------*/

func (expr *ExprVariable) Eval(env *Env) Val {
	return env.Get(expr.name)
}

/*----------  Expr: Logical  ----------*/

func (expr *ExprLogical) Eval(env *Env) Val {
	val := expr.left.Eval(env)
	if expr.operator.typ == OR {
		if getTruthy(val) {
			return val
		}
	} else {
		if !getTruthy(val) {
			return val
		}
	}
	return expr.right.Eval(env)
}

/*----------  Expr: Function Call  ----------*/

func (expr *ExprCall) Eval(env *Env) Val {
	callee := expr.callee.Eval(env)
	var arguments []Val
	for _, arg := range expr.arguments {
		arguments = append(arguments, arg.Eval(env))
	}
	if function, ok := callee.(Callable); ok {
		expected := function.Arity()
		got := len(arguments)
		if expected != got {
			panic(NewRuntimeError(expr.paren, fmt.Sprintf("expect %d arguments but got %d", expected, got)))
		}
		return function.Call(env, arguments)
	} else {
		panic(NewRuntimeError(expr.paren, "can only call functions and classes"))
	}
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
