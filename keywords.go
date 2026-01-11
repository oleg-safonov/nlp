package nlp

import "unicode/utf8"

type KeywordSet []string

var DefaultKeywords = KeywordSet{
	"?!", ":)", ";)",
	"г.", "ул.", "д.",
	"н.э.", "н. э.", "т.е.", "т. е.", "т.д.", "т. д.", "т.п.", "т. п."}

type Keywords struct {
	keywords        map[string]struct{}
	keywordPrefixes map[string]struct{}
}

func NewKeywords(keywordSets ...KeywordSet) *Keywords {
	k := Keywords{
		keywords:        map[string]struct{}{},
		keywordPrefixes: map[string]struct{}{},
	}

	for _, set := range keywordSets {
		for _, kw := range set {
			k.keywords[kw] = struct{}{}
			for i, c := range kw {
				k.keywordPrefixes[kw[0:i+utf8.RuneLen(c)]] = struct{}{}
			}
		}
	}

	return &k
}

func (k *Keywords) IsKeyword(word string) bool {
	_, ok := k.keywords[word]
	return ok
}

func (k *Keywords) IsKeywordPrefix(prefix string) bool {
	_, ok := k.keywordPrefixes[prefix]
	return ok
}
