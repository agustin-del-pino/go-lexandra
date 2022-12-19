package golexandra

const (
	eof = 0x00
)

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
		c.Char = eof
	}
}

func (c *Cursor) HasChar() bool {
	return c.Char != eof
}

func NewCursor(b []byte) *Cursor {
	return &Cursor{
		Position: 0,
		Char:     eof,
		bytes:    b,
		length:   len(b),
	}
}
