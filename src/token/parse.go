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

	if len(parser.CurrentToken.Value) == 0 {
		parser.CurrentToken.Start = [2]int{parser.Reader.row, parser.Reader.col}
	}

	switch t {
	case Unknown:
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
	case Assignment:
		parser.CurrentToken.Value = "="
	case LeftParentheses:
		parser.CurrentToken.Value = "("
	case RightParentheses:
		parser.CurrentToken.Value = ")"
	case Semicolon:
		parser.CurrentToken.Value = ";"
	case LineComment:
		parser.CurrentToken.Value = parser.collectLineComment()
	case StringValue:
		parser.CurrentToken.Value = parser.collectString()
	case IntValue:
		parser.CurrentToken.Value = parser.collectInt()
	case Identifier:
		if len(parser.CurrentToken.Value) > 0 {
			break
		}

		parser.CurrentToken.Value = parser.collectIdentifier()
		parser.setCurrentTokenType(parser.getIdentifierTokenType(parser.CurrentToken.Value))
		return
	}

	parser.CurrentToken.End = [2]int{parser.Reader.row, parser.Reader.col}
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
	result := []byte{}

	for {
		charInByte, err := parser.Reader.Next()

		if IsIllegalChar(parser.Reader.charInByte) {
			parser.Reader.ReportLineError()
		}

		if err != nil || string(charInByte) == "\"" {
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
	parser.Reader.Back()

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
	case "package":
		return Package
	default:
		return Identifier
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
		case ";":
			parser.setCurrentTokenType(Semicolon)
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
				parser.setCurrentTokenType(LineComment)
			} else if string(nextCharInByte) == "*" {
				parser.setCurrentTokenType(LeftParentheses)
			} else {
				parser.setCurrentTokenType(Unknown)
			}
		case "\"":
			parser.setCurrentTokenType(StringValue)
		default:
			if IsDigit(charInByte) {
				parser.setCurrentTokenType(IntValue)
			} else if IsLetterOrSlash(charInByte) {
				parser.setCurrentTokenType(Identifier)
			} else {
				parser.setCurrentTokenType(Unknown)
			}
		}
	}

	return parser.Tokens
}
