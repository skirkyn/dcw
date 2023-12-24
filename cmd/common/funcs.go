package common

type Function[In any, Out any] interface {
	Apply(In) (Out, error)
}

type Consumer[In any] interface {
	Consume(In) error
}

type Supplier[Out any] interface {
	Supply() (Out, error)
}

type Predicate[In any] interface {
	Test(In) bool
}
