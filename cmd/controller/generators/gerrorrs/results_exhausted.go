package gerrorrs

type PotentialResultsExhausted struct {
	s string
}

func (e PotentialResultsExhausted) Error() string {
	return e.s
}

func NewPotentialResultsExhausted() error {
	return PotentialResultsExhausted{"potential results exhausted"}
}
