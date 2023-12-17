package gerrorrs

type IncorrectVocabularyLength struct {
	s string
}

func (e IncorrectVocabularyLength) Error() string {
	return e.s
}

func NewIncorrectVocabularyLength() error {
	return IncorrectVocabularyLength{"incorrect result length"}
}
