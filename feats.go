package nlp

import "strings"

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

func (p POS) String() string {
	switch p {
	case UNKNOWN:
		return "UNKNOWN"
	case ADJ:
		return "ADJ"
	case ADP:
		return "ADP"
	case ADV:
		return "ADV"
	case AUX:
		return "AUX"
	case CCONJ:
		return "CCONJ"
	case DET:
		return "DET"
	case INTJ:
		return "INTJ"
	case NOUN:
		return "NOUN"
	case NUM:
		return "NUM"
	case PART:
		return "PART"
	case PRON:
		return "PRON"
	case PROPN:
		return "PROPN"
	case PUNCT:
		return "PUNCT"
	case SCONJ:
		return "SCONJ"
	case SYM:
		return "SYM"
	case VERB:
		return "VERB"
	}
	return "ERROR"
}

const (
	shiftVerbForm = 0
	shiftVariant  = 3
	shiftDegree   = 5
	shiftPerson   = 7
	shiftNumber   = 9
	shiftGender   = 11
	shiftCase     = 13
	shiftAnimacy  = 17
	shiftAspect   = 19
	shiftVoice    = 21
	shiftPOS      = 27
)

const (
	shortVerbFormMask FEATS = 0b111
	shortVariantMask  FEATS = 0b11
	shortDegreeMask   FEATS = 0b11
	shortPersonMask   FEATS = 0b11
	shortNumberMask   FEATS = 0b11
	shortGenderMask   FEATS = 0b11
	shortCaseMask     FEATS = 0b1111
	shortAnimacyMask  FEATS = 0b11
	shortAspectMask   FEATS = 0b11
	shortVoiceMask    FEATS = 0b11
	shortPOSMask      FEATS = 0b11111
)

const (
	VerbFormMask FEATS = shortVerbFormMask << shiftVerbForm
	VariantMask  FEATS = shortVariantMask << shiftVariant
	DegreeMask   FEATS = shortDegreeMask << shiftDegree
	PersonMask   FEATS = shortPersonMask << shiftPerson
	NumberMask   FEATS = shortNumberMask << shiftNumber
	GenderMask   FEATS = shortGenderMask << shiftGender
	CaseMask     FEATS = shortCaseMask << shiftCase
	AnimacyMask  FEATS = shortAnimacyMask << shiftAnimacy
	AspectMask   FEATS = shortAspectMask << shiftAspect
	VoiceMask    FEATS = shortVoiceMask << shiftVoice
	POSMask      FEATS = shortPOSMask << shiftPOS
)

type FEATS uint32

func (f FEATS) String() string {
	strs := []string{f.POS().String(), f.Case().String(), f.VerbForm().String(), f.Variant().String(),
		f.Gender().String(), f.Person().String(), f.Number().String(), f.Degree().String(),
		f.Animacy().String(), f.Aspect().String(), f.Voice().String()}
	filtered := make([]string, 0, len(strs))
	for _, s := range strs {
		if len(s) > 0 {
			filtered = append(filtered, s)
		}
	}
	return strings.Join(filtered, "|")
}

func (f FEATS) VerbForm() VerbForm {
	return VerbForm(f & VerbFormMask >> shiftVerbForm)
}

func (f FEATS) Variant() Variant {
	return Variant(f & VariantMask >> shiftVariant)
}

func (f FEATS) Degree() Degree {
	return Degree(f & DegreeMask >> shiftDegree)
}

func (f FEATS) Person() Person {
	return Person(f & PersonMask >> shiftPerson)
}

func (f FEATS) Number() Number {
	return Number(f & NumberMask >> shiftNumber)
}

func (f FEATS) Gender() Gender {
	return Gender(f & GenderMask >> shiftGender)
}

func (f FEATS) Case() Case {
	return Case(f & CaseMask >> shiftCase)
}

func (f FEATS) Animacy() Animacy {
	return Animacy(f & AnimacyMask >> shiftAnimacy)
}

func (f FEATS) Aspect() Aspect {
	return Aspect(f & AspectMask >> shiftAspect)
}

func (f FEATS) Voice() Voice {
	return Voice(f & VoiceMask >> shiftVoice)
}

func (f FEATS) POS() POS {
	return POS(f & POSMask >> shiftPOS)
}

func setField(f FEATS, val FEATS, shift uint, mask FEATS) FEATS {
	return (f & ^(mask << shift)) | ((val & mask) << shift)
}

func (f FEATS) SetVerbForm(v VerbForm) FEATS {
	return setField(f, FEATS(v), shiftVerbForm, shortVerbFormMask)
}
func (f FEATS) SetVariant(v Variant) FEATS {
	return setField(f, FEATS(v), shiftVariant, shortVariantMask)
}
func (f FEATS) SetDegree(d Degree) FEATS { return setField(f, FEATS(d), shiftDegree, shortDegreeMask) }
func (f FEATS) SetPerson(p Person) FEATS { return setField(f, FEATS(p), shiftPerson, shortPersonMask) }
func (f FEATS) SetNumber(n Number) FEATS { return setField(f, FEATS(n), shiftNumber, shortNumberMask) }
func (f FEATS) SetGender(g Gender) FEATS { return setField(f, FEATS(g), shiftGender, shortGenderMask) }
func (f FEATS) SetCase(c Case) FEATS     { return setField(f, FEATS(c), shiftCase, shortCaseMask) }
func (f FEATS) SetAnimacy(a Animacy) FEATS {
	return setField(f, FEATS(a), shiftAnimacy, shortAnimacyMask)
}
func (f FEATS) SetAspect(a Aspect) FEATS { return setField(f, FEATS(a), shiftAspect, shortAspectMask) }
func (f FEATS) SetVoice(v Voice) FEATS   { return setField(f, FEATS(v), shiftVoice, shortVoiceMask) }
func (f FEATS) SetPOS(p POS) FEATS       { return setField(f, FEATS(p), shiftPOS, shortPOSMask) }

const START_TAG FEATS = 0

const SuperMask = POSMask | CaseMask | NumberMask | GenderMask | VerbFormMask | PersonMask | VoiceMask | AnimacyMask | AspectMask // | DegreeMask | VariantMask

const TrigramMask = POSMask | CaseMask | NumberMask | GenderMask | VerbFormMask | PersonMask

const BigramMask = POSMask | CaseMask | NumberMask | GenderMask | VerbFormMask | PersonMask // | VoiceMask | AnimacyMask | AspectMask | DegreeMask | VariantMask

type VerbForm uint8

const (
	Inf VerbForm = iota + 1
	Fin
	Part
	Conv
)

func (vf VerbForm) String() string {
	switch vf {
	case 0:
		return ""
	case Inf:
		return "VerbForm=Inf"
	case Fin:
		return "VerbForm=Fin"
	case Part:
		return "VerbForm=Part"
	case Conv:
		return "VerbForm=Conv"
	}
	return "VerbForm=err"
}

type Variant uint8

const (
	Full Variant = iota + 1
	Short
)

func (v Variant) String() string {
	switch v {
	case 0:
		return ""
	case Full:
		return "Variant=Full"
	case Short:
		return "Variant=Short"
	}
	return "Variant=err"
}

type Degree uint8

const (
	Pos Degree = iota + 1
	Cmp
	Sup
)

func (d Degree) String() string {
	switch d {
	case 0:
		return ""
	case Pos:
		return "Degree=Pos"
	case Cmp:
		return "Degree=Cmp"
	case Sup:
		return "Degree=Sup"
	}
	return "Degree=err"
}

type Person uint8

const (
	Person1 Person = iota + 1
	Person2
	Person3
)

func (p Person) String() string {
	switch p {
	case 0:
		return ""
	case Person1:
		return "Person=Person1"
	case Person2:
		return "Person=Person2"
	case Person3:
		return "Person=Person3"
	}
	return "Person=err"
}

type Number uint8

const (
	Sing Number = iota + 1
	Plur
)

func (n Number) String() string {
	switch n {
	case 0:
		return ""
	case Sing:
		return "Number=Sing"
	case Plur:
		return "Number=Plur"
	}
	return "Number=err"
}

type Gender uint8

const (
	Neut Gender = iota + 1
	Fem
	Masc
)

func (g Gender) String() string {
	switch g {
	case 0:
		return ""
	case Neut:
		return "Gender=Neut"
	case Fem:
		return "Gender=Fem"
	case Masc:
		return "Gender=Masc"
	}
	return "Gender=err"
}

type Animacy uint8

const (
	Inan Animacy = iota + 1
	Anim
)

func (a Animacy) String() string {
	switch a {
	case 0:
		return ""
	case Inan:
		return "Animacy=Inan"
	case Anim:
		return "Animacy=Anim"
	}
	return "Animacy=err"
}

type Aspect uint8

const (
	Perf Aspect = iota + 1
	Imp
)

func (a Aspect) String() string {
	switch a {
	case 0:
		return ""
	case Perf:
		return "Aspect=Perf"
	case Imp:
		return "Aspect=Imp"
	}
	return "Aspect=err"
}

type Voice uint8

const (
	Act Voice = iota + 1
	Mid
	Pass
)

func (v Voice) String() string {
	switch v {
	case 0:
		return ""
	case Act:
		return "Voice=Act"
	case Mid:
		return "Voice=Mid"
	case Pass:
		return "Voice=Pass"
	}
	return "Voice=err"
}

type Case uint8

const (
	Nom Case = iota + 1
	Gen
	Dat
	Acc
	Ins
	Loc
	Par
	Voc
)

func (c Case) String() string {
	switch c {
	case 0:
		return ""
	case Nom:
		return "Case=Nom"
	case Gen:
		return "Case=Gen"
	case Dat:
		return "Case=Dat"
	case Acc:
		return "Case=Acc"
	case Ins:
		return "Case=Ins"
	case Loc:
		return "Case=Loc"
	case Par:
		return "Case=Par"
	case Voc:
		return "Case=Voc"
	}
	return "Case=err"
}
