package main

import (
	"fmt"

	"github.com/anc95/golang-enum-to-ts/src/token"
)

func main() {
	// file, err := os.Open("./main.go")

	// if err != nil {
	// 	panic(err)
	// }

	// defer file.Close()
	// content, _ := ioutil.ReadAll(file)

	// fmt.Print((string(content)))

	a := token.Parse("1=a")

	for _, v := range a {
		fmt.Printf("[type: %d, value: %s]\n", int(v.Type), v.Value)
	}
}
