package gerrorrs

type CustomNotSupported struct {
	s string
}

func (e CustomNotSupported) Error() string {
	return e.s
}

func NewCustomNotSupported() error {
	return CustomNotSupported{"custom not supported"}
}
