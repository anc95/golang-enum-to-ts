package token

type TokenType string

const (
	Initial TokenType = "Initial"
	// === keyword ===
	Type    = "Type"
	Package = "Package"
	Const   = "Const"

	// === type keyword ===
	IntType    = "IntType"
	StringType = "StringType"

	// === const value ===
	IOTA = "IOTA"

	// vairable
	IntValue    = "IntValue"
	StringValue = "StringValue"

	// ==== punctuation ===
	Assignment       = "Assignment"       // '='
	Semicolon        = "Semicolon"        // ";"
	LineComment      = "LineComment"      // '//'
	BlockComment     = "BlockComment"     // '/* */'
	QutationMark     = "QutationMark"     // '"'
	LeftParentheses  = "LeftParentheses"  // '('
	RightParentheses = "RightParentheses" // ')'

	Identifier = "Identifier"
	Unknown    = "Unknown"
)
