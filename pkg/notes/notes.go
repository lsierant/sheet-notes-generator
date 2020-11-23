//go:generate stringer -type=ChordType,ChordQuality,NoteModifier -output=enums_string.go
package notes

import (
	"fmt"
	"strings"
)

type NoteModifier int

const (
	NoteModifierNone  NoteModifier = 0
	NoteModifierSharp NoteModifier = +1
	NoteModifierFlat  NoteModifier = -1
)

type Note struct {
	BaseName      string
	Modifier      NoteModifier
	BaseNoteIndex int
	TrebleClef    bool
	BassClef      bool
}

func (n Note) NotesOnClefs(trebleClef bool) []Note {
	var notes []Note
	if n.TrebleClef && trebleClef {
		trebleNote := n
		trebleNote.BassClef = false
		notes = append(notes, trebleNote)
	}

	if n.BassClef {
		bassNote := n
		bassNote.TrebleClef = false
		notes = append(notes, bassNote)
	}

	return notes
}

func (n Note) ToneIndex() int {
	return n.BaseNoteIndex + int(n.Modifier)
}

func (n Note) String() string {
	result := n.Symbol()
	if n.Modifier == NoteModifierSharp {
		result += "_s"
	} else if n.Modifier == NoteModifierFlat {
		result += "_f"
	} else {
		result += "_"
	}

	if n.TrebleClef && n.BassClef {
		result += "_TB"
	} else if n.TrebleClef {
		result += "_T"
	} else if n.BassClef {
		result += "_B"
	}

	return result
}

func noteNameWithModifier(noteName string, modifier NoteModifier) string {
	switch modifier {
	case NoteModifierSharp:
		return noteName + "is"
	case NoteModifierFlat:
		if noteName == "e" {
			return "es"
		}
		return noteName + "es"
	default:
		return noteName
	}
}

func noteNameWithSharpFlatModifier(noteName string, modifier NoteModifier) string {
	switch modifier {
	case NoteModifierSharp:
		return strings.ToUpper(noteName) + "♯"
	case NoteModifierFlat:
		return strings.ToUpper(noteName) + "♭"
	default:
		return strings.ToUpper(noteName)
	}
}

func (n Note) NameWithModifier() string {
	return noteNameWithModifier(n.BaseName, n.Modifier)
}

func (n Note) NameWithSharpFlatModifier() string {
	return noteNameWithSharpFlatModifier(n.BaseName, n.Modifier)
}

func (n Note) LilypondSymbol() string {
	return fmt.Sprintf("%s%s", n.NameWithModifier(), octaveModifier(n))
}

func (n Note) Symbol() string {
	return fmt.Sprintf("%s%s", n.NameWithModifier(), octaveModifierForFileName(n))
}

func octaveModifier(note Note) string {
	modifiers := []string{",", "", "'", "''", "'''"}
	return modifiers[note.BaseNoteIndex/12]
}

func octaveModifierForFileName(note Note) string {
	modifiers := []string{"l", "", "u", "uu", "uuu"}
	return modifiers[note.BaseNoteIndex/12]
}

type Interval struct {
	FirstNote  Note
	SecondNote Note
	Scale      Scale
}

var simpleIntervals = map[int]string{
	0:  "Perfect unison",
	1:  "Minor second",
	2:  "Major second",
	3:  "Minor third",
	4:  "Major third",
	5:  "Perfect fourth",
	6:  "Tritone",
	7:  "Perfect fifth",
	8:  "Minor sixth",
	9:  "Major sixth",
	10: "Minor seventh",
	11: "Major seventh",
	12: "Perfect Octave",
}

func (i Interval) Name() string {
	switch diff := i.Distance(); {

	case diff == 6:
		cScale := "cdefgabcdefgab"
		firstIndex := strings.Index(cScale, i.FirstNote.BaseName)
		secondIndex := strings.Index(cScale[firstIndex+1:], i.SecondNote.BaseName) + firstIndex + 1

		if secondIndex-firstIndex == 3 {
			return "Augmented fourth"
		} else if secondIndex-firstIndex == 4 {
			return "Diminished fifth"
		} else {
			panic(fmt.Errorf("invalid scale diff: %d, %d, %+v", firstIndex, secondIndex, i))
		}
	case diff >= 0 && diff <= 12:
		return simpleIntervals[diff]
	default:
		panic(fmt.Errorf("interval not supported: %d", diff))
	}
}

func (i Interval) Distance() int {
	if dist := i.SecondNote.ToneIndex() - i.FirstNote.ToneIndex(); dist < 0 {
		return dist * -1
	} else {
		return dist
	}
}

func note(toneIndex int, name string, trebleClef bool, bassClef bool) Note {
	return Note{
		BaseName:      name,
		Modifier:      NoteModifierNone,
		BaseNoteIndex: toneIndex,
		TrebleClef:    trebleClef,
		BassClef:      bassClef,
	}
}

var AllNotes = []Note{
	note(-5, "g", false, true),
	note(-3, "a", false, true),
	note(-1, "b", false, true),
	note(0, "c", false, true),
	note(2, "d", false, true),
	note(4, "e", false, true),
	note(5, "f", false, true),
	note(7, "g", false, true),
	note(9, "a", false, true),
	note(11, "b", false, true),
	note(12, "c", false, true),
	note(14, "d", false, true),
	note(16, "e", true, true),
	note(17, "f", true, true),
	note(19, "g", true, true),
	note(21, "a", true, true),
	note(23, "b", true, true),
	note(24, "c", true, true),
	note(26, "d", true, true),
	note(28, "e", true, true),
	note(29, "f", true, true),
	note(31, "g", true, true),
	note(33, "a", true, true),
	note(35, "b", true, false),
	note(36, "c", true, false),
	note(38, "d", true, false),
	note(40, "e", true, false),
	note(41, "f", true, false),
	note(43, "g", true, false),
	note(45, "a", true, false),
	note(47, "b", true, false),
	note(48, "c", true, false),
	note(50, "d", true, false),
	note(52, "e", true, false),
	note(53, "f", true, false),
}
