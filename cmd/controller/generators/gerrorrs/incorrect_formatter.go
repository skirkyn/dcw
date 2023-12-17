package gerrorrs

type IncorrectFormatter struct {
	s string
}

func (e IncorrectFormatter) Error() string {
	return e.s
}

func NewIncorrectFormatter() error {
	return IncorrectFormatter{"incorrect formatter"}
}
