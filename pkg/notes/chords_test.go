package notes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChordName(t *testing.T) {
	assert.Equal(t, "C maj", Chord{RootNote: AllNotes[0], Type: ChordTypeMajorTriad}.Name())
	assert.Equal(t, "D min", Chord{RootNote: AllNotes[1], Type: ChordTypeMinorTriad}.Name())
	assert.Equal(t, "D maj7", Chord{RootNote: AllNotes[1], Type: ChordTypeDominantSeventh}.Name())
}
