package golexandra

import "fmt"

type LexRun[T any] func(*Cursor, *Lex[T]) Token[T]

type LexerExtension[T any] func(*Lexer[T], *Cursor, *Tokens[T]) func() bool

type Lex[T any] struct {
	AllowedBytes ByteContainer
	Run          LexRun[T]
}

type Lexer[T any] struct {
	dlex      *Lex[T]
	lexs      map[string]*Lex[T]
	ignores   ByteContainer
	extension LexerExtension[T]
}

func (l *Lexer[T]) Get(n string, d bool) *Lex[T] {
	if lex, ok := l.lexs[n]; !ok {
		if d {
			return l.dlex
		} else {
			return nil
		}
	} else {
		return lex
	}
}

func (l *Lexer[T]) Register(n string, x *Lex[T]) {
	l.lexs[n] = x
}

func (l *Lexer[T]) Numbs(x *Lex[T]) {
	l.Register("numbs", x)
}

func (l *Lexer[T]) Strings(x *Lex[T]) {
	l.Register("strings", x)
}

func (l *Lexer[T]) Words(x *Lex[T]) {
	l.Register("words", x)
}

func (l *Lexer[T]) Delimiters(x *Lex[T]) {
	l.Register("delimiters", x)
}

func (l *Lexer[T]) Extension(e LexerExtension[T]) {
	l.extension = e
}

func (l *Lexer[T]) Tokenize(b []byte) (*Tokens[T], error) {
	c := NewCursor(b)
	t := make(Tokens[T], 0)

	c.Advance()

	n := l.Get("numbs", true)
	s := l.Get("strings", true)
	w := l.Get("words", true)
	d := l.Get("delimiters", true)

	e := l.extension(l, c, &t)

	for c.HasChar() {
		if l.ignores(c.Char) {
			c.Advance()
			continue
		}

		if n.AllowedBytes(c.Char) {
			t = append(t, n.Run(c, n))
			continue
		}

		if s.AllowedBytes(c.Char) {
			t = append(t, s.Run(c, s))
			continue
		}

		if w.AllowedBytes(c.Char) {
			t = append(t, w.Run(c, w))
			continue
		}

		if d.AllowedBytes(c.Char) {
			t = append(t, d.Run(c, d))
			continue
		}

		if e() {
			c.Advance()
			continue
		}

		return nil, fmt.Errorf("unexcepted token %v:%c", c.Char, c.Char)
	}

	return &t, nil
}

func NewLexer[T any](i ...byte) *Lexer[T] {
	return &Lexer[T]{
		ignores: BytePoints(i...),
		lexs:    make(map[string]*Lex[T]),
		dlex: &Lex[T]{
			AllowedBytes: func(byte) bool { return false },
		},
		extension: func(*Lexer[T], *Cursor, *Tokens[T]) func() bool { return func() bool { return true } },
	}
}
