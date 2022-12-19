package golexandra

import "fmt"

type Token[T any] struct {
	Type  T
	Value []byte
}

func (t Token[T]) String() string {
	return fmt.Sprintf("[%v: %s]", t.Type, t.Value)
}

type Tokens[T any] []Token[T]

func (t Tokens[T]) String() string {
	var s string

	for _, _t := range t {
		s = fmt.Sprintf("%s%s", s, _t)
	}

	return s
}

func NewToken[T any](t T, v ...byte) Token[T] {
	return Token[T]{
		Type:  t,
		Value: v,
	}
}
