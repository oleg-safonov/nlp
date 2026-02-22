package nlp

import (
	"encoding/gob"
	"os"
	"sort"
)

type PredictionRule struct {
	Tag         FEATS
	Counter     uint32
	AppendStart uint32
	AppendLen   uint8
	Cut         uint8
}

type SuffixPredictorBase struct {
	RulePool    []PredictionRule
	NodePool    []SuffixNode
	EdgesPool   []Edge
	AppendTexts string
}

type SuffixNode struct {
	ChildrenIdx uint32
	RulesIdx    uint32
	Counter     uint32
	ChildrenLen uint8
	RulesLen    uint8
}

type Edge struct {
	Char    rune
	NodeIdx int
}

func (p *SuffixPredictorBase) getChild(n *SuffixNode, c rune) *SuffixNode {
	for i := range n.ChildrenLen {
		edge := p.EdgesPool[n.ChildrenIdx+uint32(i)]
		if edge.Char == c {
			return &p.NodePool[edge.NodeIdx]
		}
	}

	return nil
}

type Prediction struct {
	Lemma       string
	Tag         FEATS
	Score       float64
	RuleCounter uint32
	NodeCounter uint32
	MatchLen    int
}

func (p *SuffixPredictorBase) Predict(word string) []Prediction {
	runes := []rune(word)
	node := &p.NodePool[0]
	var results []Prediction

	for i := len(runes) - 1; i >= 0; i-- {
		node = p.getChild(node, runes[i])
		if node == nil {
			break
		}

		for ri := range node.RulesLen {
			rule := p.RulePool[node.RulesIdx+uint32(ri)]

			base := string(runes[:len(runes)-int(rule.Cut)])
			lemma := base + p.AppendTexts[rule.AppendStart:rule.AppendStart+uint32(rule.AppendLen)]

			results = append(results, Prediction{
				Lemma:       lemma,
				Tag:         rule.Tag,
				Score:       float64(rule.Counter) / float64(node.Counter), // TODO: поменять на посчитанную при обучении
				RuleCounter: rule.Counter,
				NodeCounter: node.Counter,
				MatchLen:    len(runes) - i,
			})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].MatchLen == results[j].MatchLen {
			return results[i].Score > results[j].Score
		}
		return results[i].MatchLen > results[j].MatchLen
	})

	return results
}

func LoadOOV(path string) (*SuffixPredictorBase, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var base SuffixPredictorBase
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&base)
	if err != nil {
		return nil, err
	}

	return &base, nil
}
