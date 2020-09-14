package notes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNameWithSharpFlatModifier(t *testing.T) {
	n := note(0, "c", false, false)
	assert.Equal(t, "C", n.NameWithSharpFlatModifier())

	n.Modifier = NoteModifierFlat
	assert.Equal(t, "C♭", n.NameWithSharpFlatModifier())

	n.Modifier = NoteModifierSharp
	assert.Equal(t, "C♯", n.NameWithSharpFlatModifier())
}
