package generators

type Generator[T any] interface {
	Next(batchSize int, resultChannel *chan []string) error
	CurrentState() ([]byte, error)
}
