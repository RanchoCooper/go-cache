package go_cache

// ByteView holds an immutable view of bytes.
type ByteView struct {
	data []byte // store the actual cache value which is read-only
}

// Len returns the view's length
func (b ByteView) Len() int {
	return len(b.data)
}

func byteClone(data []byte) []byte {
	tmp := make([]byte, len(data))
	copy(tmp, data)
	return tmp
}

// ByteSlice return a copy of the data as a byte slice.
func (b *ByteView) ByteSlice() []byte {
	return byteClone(b.data)
}

// String returns the data as a string, making a copy if necessary.
func (b *ByteView) String() string {
	return string(b.data)
}
