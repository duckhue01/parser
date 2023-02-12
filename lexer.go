package parser

import (
	"fmt"

	"github.com/duckhue01/lexer"
)

const (
	KeyTok = iota

	ComparisonTok

	StringValueTok
	NumberValueTok
	ListValueTok

	LogicalToK

	OpenBracketTok
	CLoseBracketTok
)

const (
	WHITE_SPACE                       = ' '
	DOUBLE_QUOTE                      = '"'
	OPEN_SQUARE_BRACKET               = '['
	CLOSE_SQUARE_BRACKET              = ']'
	EQUAL                             = '='
	BACK_SLASH                        = '\\'
	OPEN_BRACKET                      = '('
	CLOSE_BRACKET                     = ')'
	ALPHANUMERIC_AND_UNDERSCORE_CHARS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
	COMPARISON_OPS                    = ">=<#%"
)

var (
	AllowedComparisonOps = []string{">=", ">", "<", "<=", "!=", "=", "#", "%"}
	AllowedDSFirstChars  = []rune{'[', '"', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	AllowedDSNames       = []string{"number", "array", "string"}
	NUMBER_DIGITS        = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
)

var (
	DoubleComparisonFirstChars = []rune{'>', '<', '!'}
	ComparisonFirstChars       = []rune{'>', '<', '!', '=', '#', '%'}
)

func LexOpenBracket(l *lexer.L) lexer.StateFunc {
	ignoreWhiteSpace(l)
	l.Take("(")
	if len(l.Current()) == 1 {
		l.Inc(OPEN_BRACKET)
		l.Emit(OpenBracketTok)
	}

	return LexKey
}

func LexCloseBracket(l *lexer.L) lexer.StateFunc {
	ignoreWhiteSpace(l)
	l.Take(")")
	if len(l.Current()) == 1 {
		l.Inc(OPEN_BRACKET)
		l.Emit(CLoseBracketTok)
	}

	p := l.Peek()
	if p == lexer.EOFRune {

		if l.Count(OPEN_BRACKET) == l.Count(CLOSE_BRACKET) {
			return nil
		} else {
			l.Error(fmt.Sprintf("Expected CLoseBracketTok token, got %q", p))
		}

	}

	return LexLogicalOp
}

func LexKey(l *lexer.L) lexer.StateFunc {
	ignoreWhiteSpace(l)
	l.Take(ALPHANUMERIC_AND_UNDERSCORE_CHARS)
	if len(l.Current()) == 0 {
		l.Error(fmt.Sprintf("Expected KeyTok token, got %q", l.Peek()))
		return nil
	}

	l.Emit(KeyTok)

	return LexComparisonOp
}

func LexComparisonOp(l *lexer.L) lexer.StateFunc {
	ignoreWhiteSpace(l)
	firstChar := l.Next()

	// allowed tokens
	if !isTokenInSlice(ComparisonFirstChars, firstChar) {
		l.Error(fmt.Sprintf("Expected ComparisonTok token, got %q", firstChar))

		return nil
	}

	// single character comparison token
	if !isTokenInSlice(DoubleComparisonFirstChars, firstChar) {
		l.Emit(ComparisonTok)

		return LexValue
	}

	secondChar := l.Peek()

	if secondChar != EQUAL {
		l.Emit(ComparisonTok)

		return LexValue
	}

	l.Next()
	l.Emit(ComparisonTok)

	return LexValue
}

func LexValue(l *lexer.L) lexer.StateFunc {
	ignoreWhiteSpace(l)

	nextTok := l.Peek()

	if !isTokenInSlice(AllowedDSFirstChars, nextTok) {
		l.Error(fmt.Sprintf("Expect Value token, supported data structures %+v, got %q", AllowedDSNames, nextTok))
	}

	switch nextTok {
	//string
	case DOUBLE_QUOTE:
		lexUTF8String(l)
		l.Emit(StringValueTok)

	// array
	case OPEN_SQUARE_BRACKET:
		lexUTF8Array(l)
		l.Emit(StringValueTok)

	// number
	default:
		lexNumber(l)
		l.Emit(StringValueTok)
	}

	return LexCloseBracket
}

func LexLogicalOp(l *lexer.L) lexer.StateFunc {
	ignoreWhiteSpace(l)
	l.Take("ANDOR")

	if l.Current() == "OR" || l.Current() == "AND" {
		l.Emit(LogicalToK)
	}

	return LexOpenBracket
}
