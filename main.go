package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: lox <file>")
		os.Exit(1)
	}
	body, _ := ioutil.ReadFile(os.Args[1])
	scanner := NewScanner(string(body))
	spew.Dump(scanner.ScanTokens())
}
