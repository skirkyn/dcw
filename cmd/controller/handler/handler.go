package handler

type Handler interface {
	Handle([]byte, chan []byte, chan error)
}
