package token

import (
	"errors"
	"fmt"
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
		row, col := reader.row, reader.col
		char, err := reader.Next()

		if err != nil || (string(char) != " " && !IsIllegalChar(reader.charInByte)) {
			reader.row, reader.col = row, col
			break
		}
	}
}

func (reader *Reader) SkipLine() {
	reader.row += 1
	reader.col = -1
}

func (reader *Reader) Next() (byte, error) {
	if reader.isOverflow() {
		return 0, errors.New("Overflow")
	}

	reader.col += 1

	for reader.col >= len(reader.lines[reader.row]) || len(reader.lines[reader.row]) == 0 {
		reader.row += 1
		reader.col = 0

		if reader.isOverflow() {
			return 0, errors.New("Overflow")
		}
	}

	if reader.isOverflow() {
		return 0, errors.New("Overflow")
	}

	char := reader.lines[reader.row][reader.col]

	reader.char = string(char)
	reader.charInByte = char

	return char, nil
}

func (reader *Reader) Back() {
	if reader.col == 0 {
		reader.row -= 1

		if reader.row >= 0 {
			reader.col = len(reader.lines[reader.row]) - 1
		}
	} else {
		reader.col -= 1
	}
}

func (reader *Reader) isOverflow() bool {
	return reader.row >= len(reader.lines) || (reader.row == len(reader.lines)-1 && reader.col >= len(reader.lines[reader.row]))
}

func (reader *Reader) ReportLineError() bool {
	panic(fmt.Sprintf("Unexpect token: %s", string(reader.lines[reader.row])))
}

func readLines(s []byte) [][]byte {
	lines := [][]byte{}

	for i := 0; i < len(s); i++ {
		line := []byte{}

		for ; i < len(s); i++ {
			line = append(line, s[i])

			if string(s[i]) == "\n" {
				break
			}
		}

		lines = append(lines, line)
	}

	return lines
}

func NewReader(fileOrContent string) *Reader {
	file, err := os.Open(fileOrContent)
	var content []byte

	if err == nil {
		defer file.Close()
		fileContent, _ := ioutil.ReadAll(file)
		content = fileContent
	} else {
		content = []byte(fileOrContent)
	}

	lines := readLines(content)
	return &Reader{row: 0, col: -1, lines: lines}
}
