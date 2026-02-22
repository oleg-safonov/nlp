package nlp

import (
	"fmt"
	"math"

	"github.com/cespare/xxhash/v2"
)

type LinkType uint8

type FormText struct {
	TextStart uint32
	FormIdx   uint32
	TextLen   uint8
	FormLen   uint8
}

type Form struct {
	LemmaIdx   uint32
	FEATS      FEATS
	CountTotal uint16
	CountDocs  uint16
}

type Lemma struct {
	TextStart  uint32
	LinkIdx    uint32
	FEATS      FEATS
	CountTotal uint16
	CountDocs  uint16
	TextLen    uint8
	LinkLen    uint8
}

type Link struct {
	FromLemmaIdx uint32
	Type         LinkType
}

type StatisticalTagger struct {
	TransitionCounts map[FEATS]map[FEATS]int
	TagTotalCounts   map[FEATS]int
	UniqueWords      int
	UniqueTags       int
	Alpha            float64
}

type DictionaryBase struct {
	LinkTypes map[string]LinkType

	Texts string

	FormTexts []FormText
	Forms     []Form
	Lemmas    []Lemma
	Links     []Link

	FormTextIndex map[uint64]uint32

	Tagger StatisticalTagger

	importantLinks map[LinkType]bool
}

type LemmaRule struct {
	Cut    uint8
	Append string
	POS    string
}

type LemmatizerData struct {
	Dictionary      DictionaryBase
	SuffixPredictor SuffixPredictorBase
}

type Lemmatizer struct {
	base     LemmatizerData
	keywords *Keywords
}

func NewLemmatizer(data LemmatizerData) (*Lemmatizer, error) {
	l := Lemmatizer{
		base:     data,
		keywords: NewKeywords(DefaultKeywords),
	}

	l.base.Dictionary.importantLinks = map[LinkType]bool{}
	for _, typeText := range []string{"ADJF-ADJS", "ADJF-COMP", "INFN-VERB", "INFN-PRTF", "INFN-GRND", "PRTF-PRTS",
		"ADJF-SUPR_ejsh", "ADJF-SUPR_ajsh", "ADJF-SUPR_suppl", "ADJF-SUPR_nai", "ADJF-SUPR_slng", "NORM-ORPHOVAR",
		"SBST_MASC-SBST_FEMN", "SBST_MASC-SBST_PLUR", "ADVB-COMP"} {
		if id, ok := l.base.Dictionary.LinkTypes[typeText]; ok {
			l.base.Dictionary.importantLinks[id] = true
		} else {
			panic(fmt.Errorf("not found link type %s", typeText))
		}
	}

	return &l, nil
}

type Word struct {
	Text    string
	TokenID int
	Options []Form
	POS     POS
}

func (l *Lemmatizer) GetLogScore(prevTag, currentTag FEATS, currentWord Word) float64 {
	tagger := l.base.Dictionary.Tagger
	tagger.Alpha = 0.25
	transCount := tagger.TransitionCounts[prevTag&BigramMask][currentTag&BigramMask]
	transDenom := tagger.TagTotalCounts[prevTag&BigramMask] + int(tagger.Alpha*float64(tagger.UniqueTags))
	probTrans := (float64(transCount) + tagger.Alpha) / float64(transDenom)

	wordCount := 0
	for _, f := range currentWord.Options {
		if f.FEATS&BigramMask == currentTag&BigramMask {
			wordCount += int(f.CountTotal)
		}
	}

	coeff := 1.0
	if prevTag.POS() == VERB && currentTag.POS() == NOUN {
		if currentTag.Case() == Par || currentTag.Case() == Acc {
			coeff = 0.85
		}
	}

	wordDenom := tagger.TagTotalCounts[currentTag&BigramMask] + int(tagger.Alpha*float64(tagger.UniqueWords))
	probEmission := (float64(wordCount) + tagger.Alpha) / float64(wordDenom)

	return math.Log(probTrans)*coeff + math.Log(probEmission)
}

type ViterbiStep struct {
	LogProb float64
	BackPtr Form

	Text string
}

func (l *Lemmatizer) Viterbi(sentence []Word) []Form {
	n := len(sentence)
	if n == 0 {
		return nil
	}

	dp := make([]map[Form]ViterbiStep, n)
	for i := range dp {
		dp[i] = make(map[Form]ViterbiStep)
	}

	dict := l.base.Dictionary
	firstWord := sentence[0]
	for _, form := range firstWord.Options {
		score := l.GetLogScore(FEATS(math.MaxInt32), form.FEATS, firstWord)
		dp[0][form] = ViterbiStep{LogProb: score, BackPtr: Form{FEATS: FEATS(math.MaxInt32)},
			Text: dict.Texts[dict.Lemmas[form.LemmaIdx].TextStart : dict.Lemmas[form.LemmaIdx].TextStart+uint32(dict.Lemmas[form.LemmaIdx].TextLen)]}
	}

	for i := 1; i < n; i++ {
		currWord := sentence[i]
		prevWordForm := sentence[i-1].Options

		for _, currForm := range currWord.Options {
			bestLogProb := -math.MaxFloat64
			bestPrevForm := Form{}

			for _, prevForm := range prevWordForm {
				prevStep, ok := dp[i-1][prevForm]
				if !ok {
					continue
				}

				score := prevStep.LogProb + l.GetLogScore(prevForm.FEATS, currForm.FEATS, currWord)

				if score > bestLogProb {
					bestLogProb = score
					bestPrevForm = prevForm
				}
			}
			dp[i][currForm] = ViterbiStep{LogProb: bestLogProb, BackPtr: bestPrevForm,
				Text: dict.Texts[dict.Lemmas[currForm.LemmaIdx].TextStart : dict.Lemmas[currForm.LemmaIdx].TextStart+uint32(dict.Lemmas[currForm.LemmaIdx].TextLen)]}
		}
	}

	result := make([]Form, n)

	lastBestForm := Form{}
	lastMaxLogProb := -math.MaxFloat64
	for form, step := range dp[n-1] {
		if step.LogProb > lastMaxLogProb {
			lastMaxLogProb = step.LogProb
			lastBestForm = form
		}
	}

	currForm := lastBestForm
	for i := n - 1; i >= 0; i-- {
		result[i] = currForm
		currForm = dp[i][currForm].BackPtr
	}

	return result
}

func (l *Lemmatizer) Disambiguate(tokens []Token) []Word {
	words := make([]Word, 0, len(tokens))

	for i, token := range tokens {
		if token.Type() == TokenNumber ||
			token.Type() == TokenWord ||
			token.Type() == TokenKeyword {

			forms := l.getForms(token.Text())
			if len(forms) == 0 {
				predictions := l.base.SuffixPredictor.Predict(token.Text())
				if len(predictions) > 0 {
					matchlen := predictions[0].MatchLen
					for _, pred := range predictions {
						if pred.MatchLen < matchlen-1 {
							break
						}
						forms = append(forms, Form{FEATS: pred.Tag & BigramMask, CountTotal: uint16(pred.RuleCounter)})
					}
				} else {
					forms = append(forms, Form{FEATS: FEATS(0).SetPOS(NOUN)},
						Form{FEATS: FEATS(0).SetPOS(VERB)},
						Form{FEATS: FEATS(0).SetPOS(ADJ)},
						Form{FEATS: FEATS(0).SetPOS(ADV)})
				}
			}
			words = append(words, Word{
				Text:    token.Text(),
				TokenID: i,
				Options: forms,
			})
		}
	}

	forms := l.Viterbi(words)
	for i := range words {
		words[i].Options = []Form{forms[i]}
		words[i].POS = forms[i].FEATS.POS()
	}

	return words
}

func (l *Lemmatizer) LemmatizeTokens(tokens []Token) []string {
	results := make([]string, 0, len(tokens))

	for _, t := range tokens {
		results = append(results, t.Text())
	}

	words := l.Disambiguate(tokens)
	for _, w := range words {
		if len(w.Options) > 0 {
			form := w.Options[0]
			if form.LemmaIdx == 0 {
				predictions := l.base.SuffixPredictor.Predict(w.Text)
				if len(predictions) > 0 {
					results[w.TokenID] = predictions[0].Lemma
				}
				continue
			}
			lemma := l.base.Dictionary.Lemmas[form.LemmaIdx]
			lemma, _ = l.followLinks(lemma)
			text := l.base.Dictionary.Texts[lemma.TextStart : lemma.TextStart+uint32(lemma.TextLen)]
			results[w.TokenID] = text
		}

	}
	return results
}

func (l *Lemmatizer) LemmatizeText(text string) []string {
	tokens := Tokenize(text, l.keywords)
	return l.LemmatizeTokens(tokens)
}

func (l *Lemmatizer) LemmatizeWord(word string) string {
	word = Normalize(word)

	if res, _, _, ok := l.lemmatizeByDict(word); ok {
		return res
	}

	predictions := l.base.SuffixPredictor.Predict(word)
	if len(predictions) > 0 {
		return predictions[0].Lemma
	}

	return word
}

func (l *Lemmatizer) lemmatizeByDict(word string) (string, POS, uint16, bool) {
	forms := l.getForms(word)

	maxScore := -2_000_000_000
	maxForm := Form{}
	resLemma := Lemma{}
	for _, f := range forms {
		lemma := l.base.Dictionary.Lemmas[f.LemmaIdx]
		fromLemma, lemmaScore := l.followLinks(lemma)
		score := 50*int(f.CountDocs) + int(f.CountTotal) + int(lemmaScore)
		if score > maxScore {
			maxScore = score
			maxForm = f
			resLemma = fromLemma
		}

	}

	if maxScore > 0 {
		return l.base.Dictionary.Texts[int(resLemma.TextStart) : int(resLemma.TextStart)+int(resLemma.TextLen)], maxForm.FEATS.POS(), 0, true
	}

	return word, POS(math.MaxUint8), 0, false
}

func (l *Lemmatizer) getForms(text string) []Form {
	digest := xxhash.New()
	digest.WriteString(text)
	hash := digest.Sum64()

	if idx, ok := l.base.Dictionary.FormTextIndex[hash]; ok {
		formText := l.base.Dictionary.FormTexts[idx]
		forms := make([]Form, 0, formText.FormLen)
		for i := range formText.FormLen {
			form := l.base.Dictionary.Forms[formText.FormIdx+uint32(i)]
			forms = append(forms, form)
		}

		return forms
	}
	return nil
}

// TODO

func (l *Lemmatizer) followLinks(lemma Lemma) (Lemma, int) {
	var maxLemma Lemma
	maxScore := -2_000_000_000
	for i := range lemma.LinkLen {
		link := l.base.Dictionary.Links[int(lemma.LinkIdx)+int(i)]
		if _, ok := l.base.Dictionary.importantLinks[link.Type]; !ok {
			continue
		}

		fromLemma, score := l.followLinks(l.base.Dictionary.Lemmas[link.FromLemmaIdx])
		if score > maxScore {
			maxScore = score
			maxLemma = fromLemma
		}
	}

	if maxScore >= 0 {
		return maxLemma, maxScore + int(maxLemma.CountDocs)
	}
	return lemma, int(lemma.CountDocs) // + int(lemma.CountTotal)
}
