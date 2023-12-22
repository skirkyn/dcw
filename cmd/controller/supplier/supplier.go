package supplier

type Supplier[In any, Out any] interface {
	Next(In) (Out, error)
}
