package main

import (
	"fmt"

	golexandra "github.com/agustin-del-pino/go-lexandra"
)

type TokenType int

const (
	TNumb TokenType = iota
	TWord
)

var lexNumbs = &golexandra.Lex[TokenType]{
	AllowedBytes: golexandra.ByteRange(0x30, 0x39),
	Run: func(c *golexandra.Cursor, l *golexandra.Lex[TokenType]) golexandra.Token[TokenType] {
		t := golexandra.NewToken(TNumb, c.Char)
		
		c.Advance()

		for c.HasChar() && l.AllowedBytes(c.Char) {
			t.Value = append(t.Value, c.Char)
			c.Advance()
		}

		return t
	},
}

var lexWords = &golexandra.Lex[TokenType]{
	AllowedBytes: func(b byte) bool {
		return golexandra.ByteRange(0x41, 0x5A)(b) || golexandra.ByteRange(0x61, 0x7A)(b) || golexandra.ByteSingle(0x5F)(b)
	},
	Run: func(c *golexandra.Cursor, l *golexandra.Lex[TokenType]) golexandra.Token[TokenType] {
		t := golexandra.NewToken(TWord, c.Char)
		
		c.Advance()

		for c.HasChar() && (l.AllowedBytes(c.Char) || golexandra.ByteRange(0x30, 0x39)(c.Char)){
			t.Value = append(t.Value, c.Char)
			c.Advance()
		}

		return t
	},
}

func main() {
	lexer := golexandra.NewLexer[TokenType](0x20)
	lexer.Numbs(lexNumbs)
	lexer.Words(lexWords)

	if t, terr := lexer.Tokenize([]byte("12345 ABCD 999 182647 _fadsa333")); terr != nil {
		panic(terr)
	} else {
		fmt.Println(t)
	}
}
