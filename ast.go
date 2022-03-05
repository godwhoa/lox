package main

type Expr interface{}

type BinaryExpr struct {
	Left  Expr
	Right Expr
	Op    Token
}

type LiteralExpr struct {
	Lit Token
}

type UnaryExpr struct {
	Op    Token
	Right Expr
}

type GroupingExpr struct {
	Expression Expr
}
