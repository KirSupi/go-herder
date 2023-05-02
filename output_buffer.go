package herder

type OutputBuffer struct {
	maxLen int
	buffer []byte
}

func newOutputBuffer(maxLen int) *OutputBuffer {
	if maxLen < 0 {
		return &OutputBuffer{
			maxLen: -1,
			buffer: make([]byte, 0, 0),
		}
	}
	return &OutputBuffer{
		maxLen: maxLen,
		buffer: make([]byte, 0, 0),
	}
}

func (b *OutputBuffer) Write(p []byte) (n int, err error) {
	if b.maxLen == -1 { // if unlimited
		n = len(p)
		b.buffer = append(b.buffer, p...)
		return
	} else if b.maxLen == 0 { // if disabled
		return 0, nil
	}
	if len(b.buffer)+len(p) <= b.maxLen {
		b.buffer = append(b.buffer, p...)
		n = len(p)
	} else {
		n = b.maxLen
		if len(p) >= b.maxLen {
			if len(b.buffer) == b.maxLen {
				for i, j := 0, len(p)-b.maxLen; i < b.maxLen; i++ {
					b.buffer[i] = p[j]
					j++
				}
			} else {
				b.buffer = make([]byte, 0, b.maxLen)
				copy(p[len(p)-b.maxLen:], b.buffer)
			}
		} else {
			if len(b.buffer) != b.maxLen {
				b.buffer = make([]byte, b.maxLen, b.maxLen)
			}
			for i := len(b.buffer) + len(p) - b.maxLen; i < len(b.buffer); i++ {
				b.buffer[i] = p[i]
			}
			for i, j := len(p), 0; j < len(p); j++ {
				b.buffer[i] = p[j]
			}
		}
	}
	return
}
