package work

type Supplier[In any, Out any] interface {
	Supply(In) (Out, error)
}
