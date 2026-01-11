package nlp

type POS uint8

const (
	UNKNOWN POS = iota
	ADJ
	ADP
	ADV
	AUX
	CCONJ
	DET
	INTJ
	NOUN
	NUM
	PART
	PRON
	PROPN
	PUNCT
	SCONJ
	SYM
	VERB
	_END
)

const (
	shiftVerbForm = 0
	shiftVariant  = 2
	shiftDegree   = 3
	shiftPerson   = 5
	shiftNumber   = 7
	shiftGender   = 8
	shiftCase     = 10
	shiftPOS      = 27
)

type FEATS uint32

func (f FEATS) VerbForm() VerbForm {
	return VerbForm(f & 0b_00000000_00000000_00000000_00000011 >> shiftVerbForm)
}

func (f FEATS) Variant() Variant {
	return Variant(f & 0b_00000000_00000000_00000000_00000100 >> shiftVariant)
}

func (f FEATS) Degree() Degree {
	return Degree(f & 0b_00000000_00000000_00000000_00011000 >> shiftDegree)
}

func (f FEATS) Person() Person {
	return Person(f & 0b_00000000_00000000_00000000_01100000 >> shiftPerson)
}

func (f FEATS) Number() Number {
	return Number(f & 0b_00000000_00000000_00000000_10000000 >> shiftNumber)
}

func (f FEATS) Gender() Gender {
	return Gender(f & 0b_00000000_00000000_00000011_00000000 >> shiftGender)
}

func (f FEATS) Case() Case {
	return Case(f & 0b_00000000_00000000_01111100_00000000 >> shiftCase)
}

func (f FEATS) POS() POS {
	return POS(f & 0b_11111000_00000000_00000000_00000000 >> shiftPOS)
}

func setField(f FEATS, val uint32, shift uint, mask uint32) FEATS {
	return (f & ^(FEATS(mask) << shift)) | ((FEATS(val) & FEATS(mask)) << shift)
}

func (f FEATS) SetVerbForm(v VerbForm) FEATS { return setField(f, uint32(v), shiftVerbForm, 0b11) }
func (f FEATS) SetVariant(v Variant) FEATS   { return setField(f, uint32(v), shiftVariant, 0b1) }
func (f FEATS) SetDegree(d Degree) FEATS     { return setField(f, uint32(d), shiftDegree, 0b11) }
func (f FEATS) SetPerson(p Person) FEATS     { return setField(f, uint32(p), shiftPerson, 0b11) }
func (f FEATS) SetNumber(n Number) FEATS     { return setField(f, uint32(n), shiftNumber, 0b1) }
func (f FEATS) SetGender(g Gender) FEATS     { return setField(f, uint32(g), shiftGender, 0b11) }
func (f FEATS) SetCase(c Case) FEATS         { return setField(f, uint32(c), shiftCase, 0b11111) }
func (f FEATS) SetPOS(p POS) FEATS           { return setField(f, uint32(p), shiftPOS, 0b11111) }

type VerbForm uint8

const (
	Inf VerbForm = iota
	Fin
	Part
	Conv
)

type Variant uint8

const (
	Full Variant = iota
	Short
)

type Degree uint8

const (
	Pos Degree = iota
	Cmp
	Sup
)

type Person uint8

const (
	Person1 Person = iota
	Person2
	Person3
)

type Number uint8

const (
	Sing Number = iota
	Plur
)

type Gender uint8

const (
	Neut Gender = iota
	Fem
	Masc
)

type Case uint8

const (
	Nom Case = iota
	Gen
	Dat
	Acc
	Ins
	Loc
	Par
	Voc
)
