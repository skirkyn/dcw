package handler

type Handler interface {
	Handle([]byte) []byte
}
