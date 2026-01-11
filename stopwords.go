package nlp

type StopwordSet []string

var DefaultStopwords = StopwordSet{}

type Stopwords struct {
	stopwords map[string]struct{}
}

func NewStopwords(stopwordSets ...StopwordSet) *Stopwords {
	s := Stopwords{
		stopwords: map[string]struct{}{},
	}

	for _, set := range stopwordSets {
		for _, sw := range set {
			s.stopwords[sw] = struct{}{}

		}
	}

	return &s
}

func (s *Stopwords) IsStopword(word string) bool {
	_, ok := s.stopwords[word]
	return ok
}
