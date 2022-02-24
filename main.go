package main

import (
	"fmt"

	"github.com/anc95/golang-enum-to-ts/src/token"
)

func main() {
	parser := token.NewParser("func() {\n hell\n xxx\n dsdasdsa\n \n} \na=1\ntype C string //hello\nconst ( A C = 1 \n B \n D")

	a := parser.Parse()
	for _, v := range a {
		fmt.Printf("[type: %d, value: %s]\n", int(v.Type), v.Value)
	}
}
