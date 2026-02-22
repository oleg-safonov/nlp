package nlp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	var feats FEATS
	assert.Equal(t, "UNKNOWN", feats.String())

	feats = feats.SetPOS(DET)
	assert.Equal(t, "DET", feats.String())

	feats = feats.SetNumber(Sing)
	assert.Equal(t, "DET|Number=Sing", feats.String())

	feats = feats.SetPerson(Person2)
	assert.Equal(t, "DET|Person=Person2|Number=Sing", feats.String())

	feats = feats.SetVoice(Pass)
	assert.Equal(t, "DET|Person=Person2|Number=Sing|Voice=Pass", feats.String())

	feats = feats.SetGender(Masc)
	assert.Equal(t, "DET|Gender=Masc|Person=Person2|Number=Sing|Voice=Pass", feats.String())

	feats = feats.SetVerbForm(Conv)
	assert.Equal(t, "DET|VerbForm=Conv|Gender=Masc|Person=Person2|Number=Sing|Voice=Pass", feats.String())

	feats = feats.SetVariant(Short)
	assert.Equal(t, "DET|VerbForm=Conv|Variant=Short|Gender=Masc|Person=Person2|Number=Sing|Voice=Pass", feats.String())

	feats = feats.SetAspect(Imp)
	assert.Equal(t, "DET|VerbForm=Conv|Variant=Short|Gender=Masc|Person=Person2|Number=Sing|Aspect=Imp|Voice=Pass", feats.String())

	feats = feats.SetAnimacy(Anim)
	assert.Equal(t, "DET|VerbForm=Conv|Variant=Short|Gender=Masc|Person=Person2|Number=Sing|Animacy=Anim|Aspect=Imp|Voice=Pass", feats.String())

	feats = feats.SetDegree(Sup)
	assert.Equal(t, "DET|VerbForm=Conv|Variant=Short|Gender=Masc|Person=Person2|Number=Sing|Degree=Sup|Animacy=Anim|Aspect=Imp|Voice=Pass", feats.String())

	feats = feats.SetCase(Voc)
	assert.Equal(t, "DET|Case=Voc|VerbForm=Conv|Variant=Short|Gender=Masc|Person=Person2|Number=Sing|Degree=Sup|Animacy=Anim|Aspect=Imp|Voice=Pass", feats.String())
}
