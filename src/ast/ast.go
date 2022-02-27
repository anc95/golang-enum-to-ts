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
	usedComments map[token.Token]bool
}

func (a *AstGenerator) nextToken(reportErrorWhenIsNull bool) (token.Token, error) {
	for i := a.index + 1; ; i++ {
		if i >= len(a.Tokens) {
			if reportErrorWhenIsNull {
				a.reportTokenError()
			}

			return token.Token{}, errors.New("Overflow")
		}

		tok := a.Tokens[i]

		if tok.Type == token.LineComment {
			continue
		}

		a.index = i
		a.currentToken = a.Tokens[a.index]
		break
	}

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

func (a *AstGenerator) resolveComments(node *BaseDeclaration, leading bool) {
	currentToken := a.currentToken
	comments := []Comment{}

	if leading {
		for i := a.index - 1; i >= 0; i-- {
			tok := a.Tokens[i]

			if tok.Type == token.LineComment && !a.usedComments[tok] {
				comment := Comment{}
				comment.Value = tok.Value
				comments = append([]Comment{comment}, comments...)
				a.usedComments[tok] = true
			} else {
				break
			}
		}

		node.LeadingComments = comments
	} else {
		for i := a.index + 1; i < len(a.Tokens); i++ {
			tok := a.Tokens[i]

			if tok.Start[0] != currentToken.Start[0] {
				break
			}

			if tok.Type == token.LineComment {
				comments = []Comment{{Value: tok.Value}}
				a.usedComments[tok] = true
				break
			}

			break
		}

		node.TrailingComments = comments
	}
}

func (a *AstGenerator) readTypeDeclaration() TypeDeclaration {
	d := TypeDeclaration{}

	a.resolveComments(&d.BaseDeclaration, true)

	next, _ := a.nextToken(true)
	d.Id = next.Value

	next, _ = a.nextToken(true)

	if next.Type == token.IntType {
		d.Kind = Int
	} else {
		d.Kind = String
	}

	a.resolveComments(&d.BaseDeclaration, false)

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

		a.resolveComments(&decl.BaseDeclaration, true)
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
				a.resolveComments(&decl.BaseDeclaration, false)
				declarators = append(declarators, decl)
			} else {
				a.reportTokenError()
			}
		} else {
			a.resolveComments(&decl.BaseDeclaration, false)
			declarators = append(declarators, decl)
			if a.currentToken.Type == token.RightParentheses {
				break
			}

			a.backToken()
		}

		a.matchNextLine()
	}

	return ConstDeclaration{Declarators: declarators}
}

func (a *AstGenerator) match(t token.TokenType) {
	if t != a.currentToken.Type {
		a.reportTokenError()
	}
}

func (a *AstGenerator) matchNextLine() {
	next := a.Tokens[a.index+1]

	if next.Type == token.LineComment {
		next = a.Tokens[a.index+2]
	}

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

func NewAstGenerator(tokens []token.Token) AstGenerator {
	return AstGenerator{Tokens: tokens, index: -1, usedComments: map[token.Token]bool{}}
}
