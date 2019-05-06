package gwriter

// Writer is the interface of writer implementations
type Writer interface {
	WriteString(s string)
	WriteBytes(bs []byte)
	Start()
	Stop()
}
