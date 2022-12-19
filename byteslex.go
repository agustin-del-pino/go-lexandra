package golexandra

type ByteContainer func(byte) bool

func ByteSingle(b byte) ByteContainer {
	return func(_b byte) bool {
		return b == _b
	}
}

func ByteRange(f byte, t byte) ByteContainer {
	return func(b byte) bool {
		return b >= f && b <= t
	}
}

func BytePoints(bs ...byte) ByteContainer {
	return func(b byte) bool {
		for _, _b := range bs {
			if b == _b {
				return true
			}
		}
		return false
	}
}
