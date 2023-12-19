package generators

type Generator[T any] interface {
	Next(batchSize int) ([]T, error)
	CurrentState() ([]byte, error)
}
