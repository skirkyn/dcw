package gerrorrs

type CustomNotSupported struct {
	s string
}

func (e CustomNotSupported) Error() string {
	return e.s
}

func newCustomNotSupported() error {
	return CustomNotSupported{"custom not supported"}
}
