package ast

import (
	"errors"
	"fmt"

	"github.com/anc95/golang-enum-to-ts/src/token"
)

type AstGenerator struct {
	Tokens       []token.Token
	index        int
	currentToken token.Token
}

func (a *AstGenerator) nextToken(reportErrorWhenIsNull bool) (token.Token, error) {
	a.index += 1

	if a.index >= len(a.Tokens) {
		if reportErrorWhenIsNull {
			a.reportTokenError()
		}

		return token.Token{}, errors.New("Overflow")
	}

	a.currentToken = a.Tokens[a.index]

	return a.currentToken, nil
}

func (a *AstGenerator) backToken() {
	a.index -= 1

	a.currentToken = a.Tokens[a.index]
}

func (a *AstGenerator) reportTokenError() {
	panic(fmt.Sprintf("Unexpected token at:%s", a.currentToken.Value))
}

func (a *AstGenerator) initFile() File {
	file := File{Body: []interface{}{}}

	for {
		item, err := a.nextToken(false)

		if err != nil {
			panic("Cant find package declaration")
		}

		if item.Type == token.LineComment {
			continue
		}

		if item.Type == token.Package {
			nextItem, err := a.nextToken(false)

			if err != nil {
				panic("Cant find package declaration")
			}

			if nextItem.Type == token.Identifier {
				file.Name = nextItem.Value
			} else {
				a.reportTokenError()
			}

			break
		}

		panic("Cant find package declaration")
	}

	return file
}

func (a *AstGenerator) readTypeDeclaration() TypeDeclaration {
	d := TypeDeclaration{}

	next, _ := a.nextToken(true)
	d.Id = next.Value

	next, _ = a.nextToken(true)

	if next.Type != token.IntType {
		d.Kind = Int
	} else {
		d.Kind = String
	}

	return d
}

func (a *AstGenerator) readConstDeclaration() ConstDeclaration {
	declarators := []ConstDeclarator{}

	a.nextToken(true)
	a.match(token.LeftParentheses)

	for {
		decl := ConstDeclarator{}

		a.nextToken(true)

		if a.currentToken.Type == token.RightParentheses {
			break
		}

		a.match(token.Identifier)
		prev := a.currentToken
		decl.Id = a.currentToken.Value

		a.nextToken(true)

		if a.currentToken.Type != token.Semicolon && prev.Start[0] == a.currentToken.Start[0] {
			if a.currentToken.Type != token.Assignment {
				a.match(token.Identifier)
				decl.Kind = a.currentToken.Value
				a.nextToken(true)
			}

			a.match(token.Assignment)
			a.nextToken(true)

			if a.currentToken.Type == token.StringValue || a.currentToken.Type == token.IntValue || a.currentToken.Type == token.IOTA {
				decl.Value = a.currentToken.Value
				declarators = append(declarators, decl)
			} else {
				a.reportTokenError()
			}
		} else {
			declarators = append(declarators, decl)
			if a.currentToken.Type == token.RightParentheses {
				break
			}
		}

		a.matchNextLine()
	}

	return ConstDeclaration{declarators}
}

func (a *AstGenerator) match(t token.TokenType) {
	if t != a.currentToken.Type {
		a.reportTokenError()
	}
}

func (a *AstGenerator) matchNextLine() {
	next := a.Tokens[a.index+1]

	if next.Type == token.Semicolon {
		a.nextToken(true)
		return
	}

	if next.Start[0] > a.currentToken.Start[0] {
		return
	}

	a.reportTokenError()
}

func (a *AstGenerator) Gen() File {
	file := a.initFile()

	for {
		item, err := a.nextToken(false)

		if err != nil {
			break
		}

		switch item.Type {
		case token.LineComment:
		case token.Unknown:
			continue
		case token.Type:
			file.Body = append(file.Body, a.readTypeDeclaration())
		case token.Const:
			file.Body = append(file.Body, a.readConstDeclaration())
		default:
			a.reportTokenError()
		}
	}

	return file
}
