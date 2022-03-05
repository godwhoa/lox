package main

import "strings"

type ASTPrinter struct {
}

func (a *ASTPrinter) parenthesize(name string, exprs ...Expr) string {
	var s strings.Builder
	s.WriteString("(")
	s.WriteString(name)
	for _, expr := range exprs {
		s.WriteString(" ")
		s.WriteString(a.Print(expr))
	}
	s.WriteString(")")
	return s.String()
}

func (a *ASTPrinter) Print(expr Expr) string {
	switch expr := expr.(type) {
	case *BinaryExpr:
		return a.parenthesize(expr.Op.Lexeme, expr.Left, expr.Right)
	case *GroupingExpr:
		return a.parenthesize("group", expr.Expression)
	case *LiteralExpr:
		return expr.Lit.Lexeme
	case *UnaryExpr:
		return a.parenthesize(expr.Op.Lexeme, expr.Right)
	default:
		return "?"
	}
}
