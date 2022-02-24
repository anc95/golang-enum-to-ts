package token

import (
	"os"
	"path"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestNormal(t *testing.T) {
	wd, _ := os.Getwd()

	a := path.Join(wd, "../test-cases/normal.go")
	parser := NewParser(a)
	result := parser.Parse()

	snaps.MatchSnapshot(t, result)
}
