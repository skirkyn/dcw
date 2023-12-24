package sfv

var (
	CustomNotSupportedError        = newCustomNotSupported()
	IncorrectFormatterError        = newIncorrectFormatter()
	IncorrectResultLengthError     = newIncorrectResultLength()
	IncorrectVocabularyLengthError = newIncorrectVocabularyLength()
	PotentialResultsExhaustedError = newPotentialResultsExhausted()
	InvalidStateFileError          = newInvalidStateFile()
)

type CustomNotSupported struct {
	s string
}

func (e CustomNotSupported) Error() string {
	return e.s
}

func newCustomNotSupported() error {
	return CustomNotSupported{"custom not supported"}
}

type IncorrectFormatter struct {
	s string
}

func (e IncorrectFormatter) Error() string {
	return e.s
}

func newIncorrectFormatter() error {
	return IncorrectFormatter{"incorrect formatter"}
}

type IncorrectResultLength struct {
	s string
}

func (e IncorrectResultLength) Error() string {
	return e.s
}

func newIncorrectResultLength() error {
	return IncorrectResultLength{"incorrect result length"}
}

type IncorrectVocabularyLength struct {
	s string
}

func (e IncorrectVocabularyLength) Error() string {
	return e.s
}

func newIncorrectVocabularyLength() error {
	return IncorrectVocabularyLength{"incorrect result length"}
}

type PotentialResultsExhausted struct {
	s string
}

func (e PotentialResultsExhausted) Error() string {
	return e.s
}

func newPotentialResultsExhausted() error {
	return PotentialResultsExhausted{"potential results exhausted"}
}

type InvalidStateFile struct {
	s string
}

func (e InvalidStateFile) Error() string {
	return e.s
}

func newInvalidStateFile() error {
	return CustomNotSupported{"state file doesn't exist or invalid"}
}
