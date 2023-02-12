package parser

import (
	"testing"

	"github.com/duckhue01/lexer"
	"github.com/google/go-cmp/cmp"
)

func TestLexFunctionsWithTrueCases(t *testing.T) {
	type args struct {
		l *lexer.L
	}

	tests := []struct {
		name string
		args args
		want []*lexer.Token
	}{
		{
			name: "one key value pair",
			args: args{
				l: lexer.New(`a="b"`, LexOpenBracket, func(e string) {}),
			},
			want: []*lexer.Token{
				{
					Typ: KeyTok,
					Val: "a",
				},
				{
					Typ: ComparisonTok,
					Val: "=",
				},
				{
					Typ: StringValueTok,
					Val: "b",
				},
			},
		},
		{
			name: "one key value pair with white space in front of key token",
			args: args{
				l: lexer.New(`   a="b"`, LexOpenBracket, func(e string) {
				}),
			},
			want: []*lexer.Token{
				{
					Typ: KeyTok,
					Val: "a",
				},
				{
					Typ: ComparisonTok,
					Val: "=",
				},
				{
					Typ: StringValueTok,
					Val: "b",
				},
			},
		},
		{
			name: "one key value pair with white space in front of comparison token",
			args: args{
				l: lexer.New(`a  ="b"`, LexOpenBracket, func(e string) {
				}),
			},
			want: []*lexer.Token{
				{
					Typ: KeyTok,
					Val: "a",
				},
				{
					Typ: ComparisonTok,
					Val: "=",
				},
				{
					Typ: StringValueTok,
					Val: "b",
				},
			},
		},
		{
			name: "one key value pair with white space in front of value token",
			args: args{
				l: lexer.New(`a=  "b"`, LexOpenBracket, func(e string) {
				}),
			},
			want: []*lexer.Token{
				{
					Typ: KeyTok,
					Val: "a",
				},
				{
					Typ: ComparisonTok,
					Val: "=",
				},
				{
					Typ: StringValueTok,
					Val: "b",
				},
			},
		},
		{
			name: "one key value pair with bracket pair",
			args: args{
				l: lexer.New(`(a="b")`, LexOpenBracket, func(e string) {
				}),
			},
			want: []*lexer.Token{
				{
					Typ: OpenBracketTok,
					Val: "(",
				},
				{
					Typ: KeyTok,
					Val: "a",
				},
				{
					Typ: ComparisonTok,
					Val: "=",
				},
				{
					Typ: StringValueTok,
					Val: "b",
				},
				{
					Typ: CLoseBracketTok,
					Val: ")",
				},
			},
		},
	}
	for _, tt := range tests {

		tokens := make([]*lexer.Token, 0, 10)
		tt.args.l.Lex()
		tok, done := tt.args.l.NextToken()
		for !done {

			tokens = append(tokens, tok)
			tok, done = tt.args.l.NextToken()
		}
		t.Run(tt.name, func(t *testing.T) {
			if !cmp.Equal(tokens, tt.want) {

				t.Errorf(cmp.Diff(tokens, tt.want))
			}
		})
	}
}

func TestLexFunctionsWithFalseCases(t *testing.T) {
	type args struct {
		l *lexer.L
	}

	tests := []struct {
		name      string
		args      args
		wantError bool
		errorM    string
	}{
		{
			name: "with just open bracket token",
			args: args{
				l: lexer.New(`(`, LexOpenBracket, func(e string) {}),
			},
			errorM:    `Expected KeyTok token, got '�'`,
			wantError: true,
		},
		{
			name: "with just open bracket token",
			args: args{
				l: lexer.New(`(key`, LexOpenBracket, func(e string) {}),
			},
			errorM:    `Expected ComparisonTok token, got '�'`,
			wantError: true,
		},
		{
			name: "with just open bracket token",
			args: args{
				l: lexer.New(`(key=`, LexOpenBracket, func(e string) {}),
			},
			errorM:    `Expect Value token, supported data structures [number array string], got '�'`,
			wantError: true,
		},
		{
			name: "with just open bracket token",
			args: args{
				l: lexer.New(`(key=#`, LexOpenBracket, func(e string) {}),
			},
			errorM:    `Expect Value token, supported data structures [number array string], got '#'`,
			wantError: true,
		},
		{
			name: "with just open bracket token",
			args: args{
				l: lexer.New(`(key==`, LexOpenBracket, func(e string) {}),
			},
			errorM:    `Expect Value token, supported data structures [number array string], got '='`,
			wantError: true,
		},
		{
			name: "with just open bracket token",
			args: args{
				l: lexer.New(`(key="value"`, LexOpenBracket, func(e string) {}),
			},
			errorM:    `Expected CLoseBracketTok token, got '�'`,
			wantError: false,
		},
		{
			name: "with just open bracket token",
			args: args{
				l: lexer.New(`(key="value"AND`, LexOpenBracket, func(e string) {}),
			},
			errorM:    `Expected KeyTok token, got '�'`,
			wantError: false,
		},
		{
			name: "with just open bracket token",
			args: args{
				l: lexer.New(`(key="value"ANDkey`, LexOpenBracket, func(e string) {}),
			},
			errorM:    `Expected ComparisonTok token, got '�'`,
			wantError: false,
		},
		{
			name: "with just open bracket token",
			args: args{
				l: lexer.New(`(key="value"ANDkey=[123,"12312"]`, LexOpenBracket, func(e string) {}),
			},
			errorM:    `Expected CLoseBracketTok token, got '�'`,
			wantError: false,
		},
	}
	for _, tt := range tests {
		tt.args.l.Lex()
		_, done := tt.args.l.NextToken()
		for !done {
			_, done = tt.args.l.NextToken()
		}

		if tt.args.l.Err == nil {
			t.Errorf("Expected error, got <nil>")
			return
		}
		if !cmp.Equal(tt.errorM, tt.args.l.Err.Error()) {
			t.Errorf("Expect error %q, but got %q", tt.errorM, tt.args.l.Err.Error())
		}

	}
}
