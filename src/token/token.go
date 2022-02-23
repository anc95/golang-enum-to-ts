package token

type TokenType int

const (
	Initial TokenType = iota
	// === keyword ===
	Type
	Package
	Const

	// === type keyword ===
	IntType
	StringType

	// === const value ===
	IOTA

	// vairable
	IntValue
	StringValue

	// ==== punctuation ===
	Assignment        // '='
	LineComment       // '//'
	BlockCommentStart // '/*'
	BlockCommentEnd   // '*/'
	Div               // '/'
	EndOfFile
	EndOfLine
	QutationMark // '"'
	LineFeed

	Indetifier
	Unknown
)
