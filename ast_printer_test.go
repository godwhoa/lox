package main

import "testing"

func TestASTPrinter(t *testing.T) {
	expr := &BinaryExpr{
		Left: &UnaryExpr{
			Op:    Token{MINUS, "-", nil, 1},
			Right: &LiteralExpr{Lit: Token{NUMBER, "123", 123, 1}},
		},
		Right: &GroupingExpr{
			Expression: &LiteralExpr{Lit: Token{NUMBER, "45.67", 45.67, 1}},
		},
		Op: Token{STAR, "*", nil, 1},
	}
	printer := &ASTPrinter{}
	got := printer.Print(expr)
	want := "(* (- 123) (group 45.67))"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
