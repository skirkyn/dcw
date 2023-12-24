package server

import "sync"

type Server interface {
	Start() (*sync.WaitGroup, error)
	Stop() error
}
