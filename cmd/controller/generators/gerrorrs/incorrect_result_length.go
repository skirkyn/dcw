package gerrorrs

type IncorrectResultLength struct {
	s string
}

func (e IncorrectResultLength) Error() string {
	return e.s
}

func NewIncorrectResultLength() error {
	return IncorrectResultLength{"incorrect result length"}
}
