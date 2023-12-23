package verifier

type Verifier[In any] interface {
	Verify(In) bool
}

type SuccessPredicate interface {
	Test([]byte) bool
}
