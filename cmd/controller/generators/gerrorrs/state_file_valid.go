package gerrorrs

type InvalidStateFile struct {
	s string
}

func (e InvalidStateFile) Error() string {
	return e.s
}

func NewInvalidStateFile() error {
	return CustomNotSupported{"state file doesn't exist or invalid"}
}
