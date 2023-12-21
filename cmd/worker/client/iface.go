package client

type Client interface {
	Call([]byte) ([]byte, error)
	Close() error
}
