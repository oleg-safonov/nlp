package nlp

import (
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	TokenUnknown TokenType = iota
	TokenWord
	TokenNumber
	TokenSym
	TokenPunct
	TokenSpace
	TokenOther
	TokenKeyword
)

type tokenPart struct {
	start int
	end   int
	Type  TokenType
}

type Token struct {
	tp      TokenType
	rawText string
	parts   []tokenPart

	partsBuf [3]tokenPart
}

func (t *Token) Text() string {
	return t.rawText[t.parts[0].start:t.parts[len(t.parts)-1].end]
}

type TokenPart struct {
	Text string
	Type TokenType
}

func (t *Token) Parts() []TokenPart {
	parts := make([]TokenPart, len(t.parts))
	for i, p := range t.parts {
		parts[i].Text = t.rawText[p.start:p.end]
		parts[i].Type = p.Type
	}
	return parts
}

func (t *Token) Type() TokenType {
	return t.tp
}

func Tokenize(text string, keywords *Keywords) []Token {
	tokens := split(Normalize(text), keywords)

	tokens = mergeNumbers(tokens)
	tokens = mergeHyphenatedWords(tokens)
	tokens = mergeAbbreviationsAdvanced(tokens)
	return filterWords(tokens)
}

func CreateTokens(words []string) []Token {
	result := make([]Token, len(words))
	for i, w := range words {
		punct := 0
		num := 0
		space := 0
		word := 0

		for _, l := range w {
			if isPunct(l) {
				punct++
			}
			if isWord(l) {
				word++
			}
			if isDigit(l) {
				num++
			}
			if isSpace(l) {
				space++
			}
		}
		tp := TokenOther
		if space > 0 {
			tp = TokenSpace
		}
		if punct > 0 {
			tp = TokenPunct
		}
		if num > 0 {
			tp = TokenNumber
		}
		if word > 0 {
			tp = TokenWord
		}

		result[i].rawText = Normalize(w)
		result[i].parts = result[i].partsBuf[:1]
		result[i].parts[0] = tokenPart{start: 0, end: len(w), Type: tp}
		result[i].tp = tp
	}

	return result
}

func split(text string, keywords *Keywords) []Token {
	tokens := make([]Token, 32)
	numTokens := 0
	currTokenType := TokenUnknown
	var currPunct rune
	currTokenStart := 0

	addToken := func(start, end int) {
		if len(tokens) == numTokens {
			newTokens := make([]Token, (2 * len(tokens)))
			copy(newTokens, tokens)
			tokens = newTokens
		}
		tokens[numTokens].rawText = text
		tokens[numTokens].tp = currTokenType
		tokens[numTokens].parts = tokens[numTokens].partsBuf[:0]
		tokens[numTokens].parts = append(tokens[numTokens].parts, tokenPart{start: start, end: end, Type: currTokenType})
		numTokens++
		currTokenStart = end
	}

LOOP:
	for i, r := range text {
		if currTokenType == TokenUnknown {
			currTokenStart = i
		}

		for j, c := range text[currTokenStart:] {
			if !keywords.IsKeywordPrefix(text[currTokenStart : currTokenStart+j+utf8.RuneLen(c)]) {
				break
			}
			if keywords.IsKeyword(text[currTokenStart : currTokenStart+j+utf8.RuneLen(c)]) {
				currTokenType = TokenKeyword
				addToken(currTokenStart, currTokenStart+j+utf8.RuneLen(c))
				continue LOOP
			}
		}

		if currTokenType == TokenWord && isWord(r) {
			if currTokenType != TokenWord && i > currTokenStart {
				addToken(currTokenStart, i)
			}
			currTokenType = TokenWord
		} else if isLetter(r) {
			if currTokenType != TokenWord && i > currTokenStart {
				addToken(currTokenStart, i)
			}
			currTokenType = TokenWord
		} else if isDigit(r) {
			if currTokenType != TokenNumber && i > currTokenStart {
				addToken(currTokenStart, i)
			}
			currTokenType = TokenNumber
		} else if isSpace(r) {
			if currTokenType != TokenSpace && i > currTokenStart {
				addToken(currTokenStart, i)
			}
			currTokenType = TokenSpace
		} else if isSym(r) {
			if currTokenType != TokenSym && i > currTokenStart {
				addToken(currTokenStart, i)
			}
			currTokenType = TokenSym

		} else if isPunct(r) {
			if (currTokenType != TokenPunct || currPunct != r) && i > currTokenStart {
				addToken(currTokenStart, i)
			}

			currPunct = r
			currTokenType = TokenPunct
		} else {
			if currTokenType != TokenOther && i > currTokenStart {
				addToken(currTokenStart, i)
			}
			currTokenType = TokenOther
		}

	}
	if currTokenType != TokenUnknown && currTokenStart < len(text) {
		addToken(currTokenStart, len(text))
	}

	return tokens[:numTokens]
}

func mergeTokens(tokens []Token, res *Token) {
	if len(tokens) == 0 {
		*res = Token{}
	}
	res.rawText = tokens[0].rawText
	res.partsBuf = tokens[0].partsBuf
	if len(tokens[0].parts) <= len(res.partsBuf) {
		res.parts = res.partsBuf[:len(tokens[0].parts)]
	} else {
		res.parts = tokens[0].parts
	}
	res.tp = tokens[0].tp

	for m := 1; m < len(tokens); m++ {
		res.parts = append(res.parts, tokens[m].parts...)
	}
}

func mergeNumbers(tokens []Token) []Token {
	currToken := 0

	i := 0
	for i < len(tokens) {
		if i+2 < len(tokens) &&
			tokens[i].tp == TokenNumber &&
			tokens[i+1].tp == TokenPunct &&
			(tokens[i+1].Text() == "." || tokens[i+1].Text() == "," || tokens[i+1].Text() == "/" || tokens[i+1].Text() == ":") &&
			tokens[i+2].tp == TokenNumber {

			mergeTokens(tokens[i:i+3], &tokens[currToken])
			tokens[currToken].tp = TokenWord
			currToken++

			i += 3
			continue
		}

		mergeTokens(tokens[i:i+1], &tokens[currToken])
		currToken++
		i++
	}

	return tokens[:currToken]
}

func mergeHyphenatedWords(tokens []Token) []Token {
	currToken := 0
	i := 0

	for i < len(tokens) {
		if tokens[i].tp == TokenWord || tokens[i].tp == TokenNumber {
			start_token := i

			j := i + 1
			for j < len(tokens) {
				midToken := tokens[j].Text()
				if tokens[j].tp != TokenWord && tokens[j].tp != TokenNumber && !(tokens[j].tp == TokenOther && utf8.RuneCountInString(midToken) == 1 && isHyphenRune([]rune(midToken)[0])) {
					break
				}
				j++
			}

			if start_token+1 < j {
				mergeTokens(tokens[start_token:j], &tokens[currToken])
				tokens[currToken].tp = TokenWord
				currToken++
				i += j - start_token
				continue
			}
		}

		mergeTokens(tokens[i:i+1], &tokens[currToken])
		currToken++
		i++
	}

	return tokens[:currToken]
}

func mergeAbbreviationsAdvanced(tokens []Token) []Token {
	currToken := 0
	i := 0

	for i < len(tokens) {
		if tokens[i].tp == TokenWord {
			if i+3 < len(tokens) && tokens[i+1].Text() == "." {
				if (tokens[i+2].tp == TokenPunct && tokens[i+2].Text() != ".") || unicode.IsUpper([]rune(tokens[i].Text())[0]) {
					mergeTokens(tokens[i:i+2], &tokens[currToken])
					tokens[currToken].tp = TokenWord
					currToken++
					i += 2
					continue
				}
				if tokens[i+2].tp == TokenSpace {
					if tokens[i+3].tp == TokenWord && unicode.IsLower([]rune(tokens[i+3].Text())[0]) {
						mergeTokens(tokens[i:i+2], &tokens[currToken])
						tokens[currToken].tp = TokenWord
						currToken++
						i += 2
						continue
					}

					if tokens[i+3].tp == TokenNumber {
						mergeTokens(tokens[i:i+2], &tokens[currToken])
						tokens[currToken].tp = TokenWord
						currToken++
						i += 2
						continue
					}
				}
			}
		}

		mergeTokens(tokens[i:i+1], &tokens[currToken])
		currToken++
		i++
	}

	return tokens[:currToken]
}

func filterWords(tokens []Token) []Token {
	currToken := 0

	for i, t := range tokens {
		//if t.Type == TokenWord || t.Type == TokenNumber || t.Type == TokenSym || t.Type == TokenPunct || t.Type == TokenOther {
		if t.tp != TokenSpace && t.tp != TokenUnknown {
			mergeTokens(tokens[i:i+1], &tokens[currToken])
			currToken++
		}
	}

	return tokens[:currToken]
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isSym(r rune) bool {
	return r == '%' || r == '+' || r == '°'
}

func isWord(r rune) bool {
	if isLetter(r) || isDigit(r) || unicode.IsMark(r) {
		return true
	}
	return false
}

func isHyphenRune(r rune) bool {
	switch r {
	case '-', '‐', '‑', '‒':
		return true
	default:
		return false
	}
}

func isPunct(r rune) bool {
	if isHyphenRune(r) {
		return false
	}
	return unicode.IsPunct(r)
}
