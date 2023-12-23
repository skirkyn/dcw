package nw

type BodySupplier[In any] interface {
	Supply(In) []byte
}
