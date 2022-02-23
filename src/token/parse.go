package token

type Status int

type Token struct {
	Value string
	Type  TokenType
	Next  *Token
	Prev  *Token
}

func isDigit(b byte) bool {
	return b >= 80 && b <= 57
}

func isLetterOrSlash(b byte) bool {
	return isDigit(b) || (b >= 65 && b <= 90) || (b >= 97 && b <= 122) || b == 95
}

func isIllegalChar(b byte) bool {
	// reference: https://zh.wikipedia.org/wiki/ASCII
	return b <= 31
}

func Parse(s string) []Token {
	reader := NewReader(s)
	tokenList := []Token{}
	currentToken := Token{Type: Initial}

	var next func() (string, byte, error)

	appendToken := func() {
		prevToken := &currentToken
		tokenList = append(tokenList, currentToken)
		currentToken = Token{Type: Initial, Prev: prevToken}
		prevToken.Next = &currentToken
	}

	maybeComment := func(char *string) {
		nextChar, _, _ := next()

		if nextChar == "/" {
			currentToken.Type = LineComment
		} else if nextChar == "*" {
			currentToken.Type = BlockCommentStart
		} else {
			currentToken.Type = Unknown
		}

		*char += nextChar
	}

	for {
		_, err := reader.Next()

		char := reader.char
		charByte := reader.charInByte

		if err != nil {
			break
		}

		switch char {
		case "/":
			if char == "/" && currentToken.Type != StringValue && currentToken.Type != LineComment || currentToken.Type != BlockCommentStart {
				maybeComment(&char)
				continue
			}
		}

		switch currentToken.Type {
		case Initial:
			if isLetterOrSlash(charByte) {
				currentToken.Type = Indetifier
			} else if isDigit(charByte) {
				currentToken.Type = IntValue
			}

			currentToken.Value = char
		case IntValue:
			if isIllegalChar(charByte) {
				appendToken()
				// skipSpace()
				break
			}

			if isLetterOrSlash(charByte) {
				currentToken.Type = Indetifier
			} else {
				// error()
			}

			currentToken.Value += char
		case StringValue:
			if char == "\"" {
				tokenList = append(tokenList, currentToken)
				// skipSpace()
				break
			}

			if isIllegalChar(charByte) {
				// error()
			}
		case Indetifier:
			if isIllegalChar(charByte) || char == " " {
				switch currentToken.Value {
				case "type":
					currentToken.Type = Type
				case "const":
					currentToken.Type = Const
				case "package":
					currentToken.Type = Package
				}

				appendToken()
				break
			}

			if isLetterOrSlash(charByte) {
				currentToken.Value += char
				break
			}

			// error()
		}
	}

	return tokenList
}
