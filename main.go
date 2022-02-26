package main

import (
	"github.com/anc95/golang-enum-to-ts/src/ast"
	"github.com/anc95/golang-enum-to-ts/src/generator"
	"github.com/anc95/golang-enum-to-ts/src/token"
)

func Transform(s string) string {
	parser := token.NewParser(s)
	tokens := parser.Parse()

	astGenerator := ast.NewAstGenerator(tokens)
	return generator.GenerateTS(astGenerator.Gen())
}
