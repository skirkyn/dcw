package gerrorrs

type IncorrectFormatter struct {
	s string
}

func (e IncorrectFormatter) Error() string {
	return e.s
}

func newIncorrectFormatter() error {
	return IncorrectFormatter{"incorrect formatter"}
}
