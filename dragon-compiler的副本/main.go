package main

import (
	"fmt"
	"io"
	"lexer"
	"simple_parser"
)

func main() {
	source := "9-5+2"
	my_lexer := lexer.NewLexer(source)
	parser := simple_parser.NewSimpleParser(my_lexer)
	err := parser.Parse()
	if err == io.EOF {
		fmt.Println("parsing success")
	} else {
		fmt.Println("source is ilegal : ", err)
	}
}
