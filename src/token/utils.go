package token

func IsDigit(b byte) bool {
	return b >= 48 && b <= 57
}

func IsLetterOrSlash(b byte) bool {
	return (b >= 65 && b <= 90) || (b >= 97 && b <= 122) || b == 95
}

func IsIllegalChar(b byte) bool {
	// reference: https://zh.wikipedia.org/wiki/ASCII
	return b <= 31
}
