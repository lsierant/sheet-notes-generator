package notes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDegreeOfNoteInScale(t *testing.T) {
	assert.Equal(t, 0, degreeOfNoteInScale("c", CMajorScale))
	assert.Equal(t, 1, degreeOfNoteInScale("d", CMajorScale))
	assert.Equal(t, 2, degreeOfNoteInScale("e", CMajorScale))
	assert.Equal(t, 3, degreeOfNoteInScale("f", CMajorScale))
	assert.Equal(t, 4, degreeOfNoteInScale("g", CMajorScale))
	assert.Equal(t, 5, degreeOfNoteInScale("a", CMajorScale))
	assert.Equal(t, 6, degreeOfNoteInScale("b", CMajorScale))
	assert.Panics(t, func() {
		degreeOfNoteInScale("x", CMajorScale)
	})

	assert.Equal(t, 0, degreeOfNoteInScale("a", AMajorScale))
	assert.Equal(t, 1, degreeOfNoteInScale("b", AMajorScale))
	assert.Equal(t, 2, degreeOfNoteInScale("c", AMajorScale))
	assert.Equal(t, 3, degreeOfNoteInScale("d", AMajorScale))
	assert.Equal(t, 4, degreeOfNoteInScale("e", AMajorScale))
	assert.Equal(t, 5, degreeOfNoteInScale("f", AMajorScale))
	assert.Equal(t, 6, degreeOfNoteInScale("g", AMajorScale))

}
