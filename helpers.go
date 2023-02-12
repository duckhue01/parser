package parser

import (
	"fmt"
	"unicode/utf8"

	"github.com/duckhue01/lexer"
)

func ignoreWhiteSpace(l *lexer.L) {
	for l.Next() == WHITE_SPACE {
		l.Ignore()
	}
	l.Rewind() // rewind last token
}

func isTokenInSlice[T string | rune](ss []T, s T) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}

	return false
}

func lexUTF8String(l *lexer.L) {
	l.Next()
	l.Skip()
	r := l.Next()
	escape := false
	for {
		isValidRune := utf8.ValidRune(r)
		if isValidRune {
			if r == DOUBLE_QUOTE && !escape {
				l.Skip()
				break
			}

			if escape {
				escape = false
			} else {
				if r == BACK_SLASH {
					l.Skip()
					escape = true
				}
			}

			r = l.Next()
		} else {
			l.Error(fmt.Sprintf("Invalid rune %q", r))
		}
	}
}

func lexUTF8Array(l *lexer.L) {
	l.Next()
	r := l.Next()
	for {
		if r == CLOSE_SQUARE_BRACKET {
			break
		}
		if r == DOUBLE_QUOTE {
			lexUTF8String(l)
		}

		if isTokenInSlice(NUMBER_DIGITS, r) {
			lexNumber(l)
		}
		r = l.Next()
	}
}

func lexNumber(l *lexer.L) {
	l.Take("0123456789.")
}
