package main

import (
	"fmt"

	golex "github.com/agustin-del-pino/go-lexandra"
)

type TokenType int

const (
	TNumb TokenType = iota
	TWord
	TComment
)

var lexNumbs = &golex.Lex[TokenType]{
	AllowedBytes: golex.ByteRange(0x30, 0x39),
	Run: func(c *golex.Cursor, l *golex.Lex[TokenType]) golex.Token[TokenType] {
		t := golex.NewToken(TNumb, c.Char)

		c.Advance()

		for c.HasChar() && l.AllowedBytes(c.Char) {
			t.Value = append(t.Value, c.Char)
			c.Advance()
		}

		return t
	},
}

var lexWords = &golex.Lex[TokenType]{
	AllowedBytes: func(b byte) bool {
		return golex.ByteRange(0x41, 0x5A)(b) || golex.ByteRange(0x61, 0x7A)(b) || golex.ByteSingle(0x5F)(b)
	},
	Run: func(c *golex.Cursor, l *golex.Lex[TokenType]) golex.Token[TokenType] {
		t := golex.NewToken(TWord, c.Char)

		c.Advance()

		for c.HasChar() && (l.AllowedBytes(c.Char) || golex.ByteRange(0x30, 0x39)(c.Char)) {
			t.Value = append(t.Value, c.Char)
			c.Advance()
		}

		return t
	},
}

var lexComment = &golex.Lex[TokenType]{
	AllowedBytes: golex.ByteSingle(0x23),
	Run: func(c *golex.Cursor, l *golex.Lex[TokenType]) golex.Token[TokenType] {
		t := golex.NewToken(TComment)
		c.Advance()
		for c.HasChar() && c.Char != 0x0A {
			t.Value = append(t.Value, c.Char)
			c.Advance()
		}
		return t
	},
}

func lexExtend(l *golex.Lexer[TokenType], c *golex.Cursor, t *golex.Tokens[TokenType]) func() bool {
	m := l.Get("comments", true)

	return func() bool {
		if m.AllowedBytes(c.Char) {
			*t = append(*t, m.Run(c, m))
			return true
		}

		return false
	}
}

func main() {
	lexer := golex.NewLexer[TokenType](0x20, 0x0D, 0x0A)
	lexer.Numbs(lexNumbs)
	lexer.Words(lexWords)
	lexer.Register("comments", lexComment)
	lexer.Extension(lexExtend)

	if t, terr := lexer.Tokenize([]byte(`
    # this is a comment
    12345 Hola Mundo _342DDD`)); terr != nil {
		panic(terr)
	} else {
		fmt.Println(t)
	}
}
