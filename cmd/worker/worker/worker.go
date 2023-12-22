package worker

type Worker[Out any] interface {
	Process([]byte) (Out, error)
}
