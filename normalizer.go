package nlp

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Normalize(word string) string {
	word = strings.ToLower(strings.TrimSpace(word))

	word = strings.ReplaceAll(word, "ё", "е")

	if !needsTransformation(word) {
		return word
	}

	removeMarksButKeepBreve := runes.Remove(runes.Predicate(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) && r != '\u0306' // й
	}))

	t := transform.Chain(norm.NFD, removeMarksButKeepBreve, norm.NFC)
	result, _, _ := transform.String(t, word)
	return result
}

func needsTransformation(s string) bool {
	if !norm.NFC.IsNormalString(s) {
		return true
	}

	for _, r := range s {
		if unicode.Is(unicode.Mn, r) {
			return true
		}
	}

	return false
}
