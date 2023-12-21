package generators

type Generator[In, Out any] interface {
	Next(In) (Out, error)
	CurrentState() ([]byte, error)
}
