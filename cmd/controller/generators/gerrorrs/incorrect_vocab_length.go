package gerrorrs

type IncorrectVocabularyLength struct {
	s string
}

func (e IncorrectVocabularyLength) Error() string {
	return e.s
}

func newIncorrectVocabularyLength() error {
	return IncorrectVocabularyLength{"incorrect result length"}
}
