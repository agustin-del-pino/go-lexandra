package golexandra

import "fmt"

type LexRun[T any] func(*Cursor, *Lex[T]) Token[T]

type LexerExtension[T any] func(*Cursor, *Tokens[T]) bool

type Lex[T any] struct {
	AllowedBytes ByteContainer
	Run          LexRun[T]
}

type Lexer[T any] struct {
	numbs      *Lex[T]
	strings    *Lex[T]
	words      *Lex[T]
	delimiters *Lex[T]
	ignores    ByteContainer

	extension LexerExtension[T]
}

func (l *Lexer[T]) Numbs(x *Lex[T]) {
	l.numbs = x
}

func (l *Lexer[T]) Strings(x *Lex[T]) {
	l.strings = x
}

func (l *Lexer[T]) Words(x *Lex[T]) {
	l.words = x
}

func (l *Lexer[T]) Delimiters(x *Lex[T]) {
	l.delimiters = x
}

func (l *Lexer[T]) Extension(e LexerExtension[T]) {
	l.extension = e
}

func (l *Lexer[T]) Tokenize(b []byte) (*Tokens[T], error) {
	c := NewCursor(b)
	t := make(Tokens[T], 0)

	c.Advance()

	for c.HasChar() {
		if l.ignores(c.Char) {
			c.Advance()
			continue
		}

		if l.numbs.AllowedBytes(c.Char) {
			t = append(t, l.numbs.Run(c, l.numbs))
			continue
		}

		if l.strings.AllowedBytes(c.Char) {
			t = append(t, l.strings.Run(c, l.strings))
			continue
		}

		if l.words.AllowedBytes(c.Char) {
			t = append(t, l.words.Run(c, l.words))
			continue
		}

		if l.delimiters.AllowedBytes(c.Char) {
			t = append(t, l.delimiters.Run(c, l.delimiters))
			continue
		}

		if l.extension(c, &t) {
			c.Advance()
			continue
		}

		return nil, fmt.Errorf("unexcepted token %v", c.Char)
	}

	return &t, nil
}

func NewLexer[T any](i ...byte) *Lexer[T] {
	return &Lexer[T]{
		ignores: BytePoints(i...),
		numbs: &Lex[T]{
			AllowedBytes: func(byte) bool { return false },
		},
		strings: &Lex[T]{
			AllowedBytes: func(byte) bool { return false },
		},
		words: &Lex[T]{
			AllowedBytes: func(byte) bool { return false },
		},
		delimiters: &Lex[T]{
			AllowedBytes: func(byte) bool { return false },
		},
		extension: func(c *Cursor, t *Tokens[T]) bool { return true },
	}
}
