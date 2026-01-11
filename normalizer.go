package nlp

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func Normalize(word string) string {
	word = strings.ToLower(strings.TrimSpace(word))

	word = strings.ReplaceAll(word, "ё", "е")

	removeMarksButKeepBreve := runes.Remove(runes.Predicate(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) && r != '\u0306' // й
	}))

	t := transform.Chain(norm.NFD, removeMarksButKeepBreve, norm.NFC)
	result, _, _ := transform.String(t, word)
	return result
}
