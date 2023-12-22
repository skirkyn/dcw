package verifier

type Verifier[In any, Out any] interface {
	Verify(In) (Out, error)
}
