package gerrorrs

var (
	CustomNotSupportedError        = newCustomNotSupported()
	IncorrectFormatterError        = newIncorrectFormatter()
	IncorrectResultLengthError     = newIncorrectResultLength()
	IncorrectVocabularyLengthError = newIncorrectVocabularyLength()
	PotentialResultsExhaustedError = newPotentialResultsExhausted()
	InvalidStateFileError          = newInvalidStateFile()
)
