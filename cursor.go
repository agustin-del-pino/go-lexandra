package golexandra

type Cursor struct {
	Char     byte
	Position int
	length   int
	bytes    []byte
}

func (c *Cursor) Advance() {
	if c.Position < c.length {
		c.Char = c.bytes[c.Position]
		c.Position += 1
	} else {
		c.Char = EOF
	}
}

func (c *Cursor) HasChar() bool {
	return c.Char != EOF
}

func NewCursor(b []byte) *Cursor {
	return &Cursor{
		Position: 0,
		Char: EOF,
		bytes: b,
		length: len(b),
	}
}

