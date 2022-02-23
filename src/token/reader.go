package token

import (
	"errors"
	"io/ioutil"
	"os"
)

type Reader struct {
	row        int
	col        int
	char       string
	charInByte byte
	lines      [][]byte
}

func (reader *Reader) SkipSpace() {
	for {
		char, err := reader.Next()

		if err != nil || string(char) != " " {
			break
		}
	}
}

func (reader *Reader) SkipLine() {
	reader.row += 1
	reader.col = 0
}

func (reader *Reader) Next() (byte, error) {
	if reader.isOverflow() {
		return 0, errors.New("Overflow")
	}

	char := reader.lines[reader.row][reader.col]

	reader.char = string(char)
	reader.charInByte = char

	reader.col += 1

	if reader.col == len(reader.lines[reader.row]) {
		reader.row += 1
		reader.col = 0
	}

	return char, nil
}

func (reader *Reader) isOverflow() bool {
	return reader.row >= len(reader.lines)
}

func readLines(s []byte) [][]byte {
	lines := [][]byte{}

	for i := 0; i < len(s); i++ {
		line := []byte{}

		for {
			if s[i] == "\n"[0] {
				break
			}

			line = append(line, s[i])
		}

		lines = append(lines, line)
	}

	return lines
}

func NewReader(fileOrContent string) *Reader {
	file, err := os.Open("./main.go")
	content := []byte{}

	if err == nil {
		defer file.Close()
		fileContent, _ := ioutil.ReadAll(file)
		content = fileContent
	} else {
		content = []byte(fileOrContent)
	}

	lines := readLines(content)

	return &Reader{row: 0, col: 0, lines: lines}
}
