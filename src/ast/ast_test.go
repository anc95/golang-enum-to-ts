package ast

import (
	"os"
	"path"
	"testing"

	"github.com/anc95/golang-enum-to-ts/src/token"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestNormal(t *testing.T) {
	wd, _ := os.Getwd()

	a := path.Join(wd, "../test-cases/normal.go")
	parser := token.NewParser(a)
	tokens := parser.Parse()
	ast := AstGenerator{Tokens: tokens, index: -1}

	result := ast.Gen()

	snaps.MatchSnapshot(t, result)
}
