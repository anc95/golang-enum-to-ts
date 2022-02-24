package token

type Status int

type Token struct {
	Value string
	Type  TokenType
	Start [2]int
	End   [2]int
}

type Parser struct {
	CurrentToken       Token
	PrevToken          Token
	Tokens             []Token
	Reader             Reader
	inConstDeclaration bool
}

func (parser *Parser) appendToken() {
	parser.Tokens = append(parser.Tokens, parser.CurrentToken)
	parser.PrevToken = parser.CurrentToken
	parser.CurrentToken = Token{Type: Initial}
	parser.Reader.SkipSpace()
}

func (parser *Parser) setCurrentTokenType(t TokenType) {
	parser.CurrentToken.Type = t
	parser.CurrentToken.Start = [2]int{parser.Reader.row, parser.Reader.col}

	if t == Unknown {
		parser.CurrentToken.Start = [2]int{parser.Reader.row, 0}

		index := len(parser.Tokens) - 1

		for index >= 0 {
			if parser.Tokens[index].Start[0] != parser.Reader.row {
				break
			}

			index -= 1
		}

		parser.Tokens = parser.Tokens[0 : index+1]

		parser.CurrentToken.Value = parser.collectUnknown()
	} else if t == Assignment {
		parser.CurrentToken.Value = "="
	} else if t == LeftParentheses {
		parser.CurrentToken.Value = "("
	} else if t == RightParentheses {
		parser.CurrentToken.Value = ")"
	}

	parser.appendToken()
}

func (parser *Parser) collectInt() string {
	result := []byte{parser.Reader.charInByte}

	for {
		charInByte, err := parser.Reader.Next()

		if err != nil || !IsDigit(charInByte) {
			parser.Reader.Back()
			break
		}

		result = append(result, parser.Reader.charInByte)
	}

	return string(result)
}

func (parser *Parser) collectIdentifier() string {
	result := []byte{parser.Reader.charInByte}

	for {
		charInByte, err := parser.Reader.Next()

		if err != nil || !IsLetterOrSlash(charInByte) {
			parser.Reader.Back()
			break
		}

		result = append(result, parser.Reader.charInByte)
	}

	return string(result)
}

func (parser *Parser) collectString() string {
	result := []byte{parser.Reader.charInByte}

	for {
		charInByte, err := parser.Reader.Next()

		if IsIllegalChar(parser.Reader.charInByte) {
			parser.Reader.ReportLineError()
		}

		if err != nil || string(charInByte) != "\"" {
			parser.Reader.Back()
			break
		}

		result = append(result, parser.Reader.charInByte)
	}

	return string(result)
}

func (parser *Parser) collectLineComment() string {
	row := parser.Reader.lines[parser.Reader.row]
	result := string(row[parser.Reader.col+1:])

	parser.Reader.SkipLine()

	return result
}

func (parser *Parser) collectUnknown() string {
	parser.Reader.col = -1
	result := []byte{}
	firstFlag := true

	for {
		_, err := parser.Reader.Next()

		if err != nil || (!firstFlag && IsLetterOrSlash(parser.Reader.charInByte) && parser.Reader.col == 0) {
			parser.Reader.Back()
			break
		}

		firstFlag = false

		result = append(result, parser.Reader.charInByte)
	}

	return string(result)
}

func (parser *Parser) getIdentifierTokenType(id string) TokenType {
	switch id {
	case "const":
		parser.Reader.SkipSpace()
		_, err := parser.Reader.Next()

		if err != nil {
			parser.Reader.ReportLineError()
		}

		if parser.Reader.char != "(" {
			return Unknown
		}

		parser.Reader.Back()
		parser.inConstDeclaration = true
		return Const
	case "type":
		return Type
	case "string":
		return StringType
	case "int":
		return IntType
	case "iota":
		return IOTA
	default:
		return Indetifier
	}
}

func NewParser(s string) Parser {
	reader := NewReader(s)

	return Parser{
		Reader:       *reader,
		CurrentToken: Token{Type: Initial},
		Tokens:       []Token{},
	}
}

func (parser *Parser) Parse() []Token {
	for {
		parser.Reader.SkipSpace()
		charInByte, err := parser.Reader.Next()

		if err != nil {
			break
		}

		switch string(charInByte) {
		case "=":
			if parser.inConstDeclaration {
				parser.setCurrentTokenType(Assignment)
			} else {
				parser.setCurrentTokenType(Unknown)
			}
		case "(":
			if parser.PrevToken.Type == Const {
				parser.setCurrentTokenType(LeftParentheses)
			} else {
				parser.setCurrentTokenType(Unknown)
			}
		case ")":
			parser.setCurrentTokenType(RightParentheses)
			parser.inConstDeclaration = false
		case "/":
			nextCharInByte, err := parser.Reader.Next()

			if err != nil {
				parser.Reader.ReportLineError()
			}

			if string(nextCharInByte) == "/" {
				parser.CurrentToken.Value = parser.collectLineComment()
				parser.setCurrentTokenType(LineComment)
			} else if string(nextCharInByte) == "*" {
				parser.setCurrentTokenType(LeftParentheses)
			} else {
				parser.setCurrentTokenType(Unknown)
			}
		case "\"":
			parser.setCurrentTokenType(StringValue)
			parser.CurrentToken.Value = parser.collectString()
		default:
			if IsDigit(charInByte) {
				parser.CurrentToken.Value = parser.collectInt()
				parser.setCurrentTokenType(IntValue)
			} else if IsLetterOrSlash(charInByte) {
				parser.CurrentToken.Value = parser.collectIdentifier()
				parser.setCurrentTokenType(parser.getIdentifierTokenType(parser.CurrentToken.Value))
			} else {
				parser.setCurrentTokenType(Unknown)
			}
		}
	}

	return parser.Tokens
}
