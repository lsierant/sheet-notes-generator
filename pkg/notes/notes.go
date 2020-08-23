package notes

import "fmt"

type NoteModifier int

const (
	NoteModifierNone  NoteModifier = 0
	NoteModifierSharp NoteModifier = 1
	NoteModifierFlat  NoteModifier = 2
)

type Note struct {
	Name       string
	Modifier   NoteModifier
	Index      int
	TrebleClef bool
	BassClef   bool
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

func (n Note) LilypondSymbol() string {
	return fmt.Sprintf("%s%s", n.Name, octaveModifier(n))
}

func (n Note) Symbol() string {
	return fmt.Sprintf("%s%s", n.Name, octaveModifierForFileName(n))
}

func octaveModifier(note Note) string {
	modifiers := []string{",", "", "'", "''", "'''"}
	return modifiers[note.Index/12]
}

func octaveModifierForFileName(note Note) string {
	modifiers := []string{"l", "", "u", "uu", "uuu"}
	return modifiers[note.Index/12]
}

type Interval struct {
	FirstNote  Note
	SecondNote Note
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
	case diff >= 0 && diff <= 12:
		return simpleIntervals[diff]
	default:
		panic(fmt.Errorf("interval not supported: %d", diff))
	}

}

func (i Interval) Distance() int {
	return i.SecondNote.Index - i.FirstNote.Index
}

func note(index int, name string, trebleClef bool, bassClef bool) Note {
	return Note{
		Name:       name,
		Modifier:   NoteModifierNone,
		Index:      index,
		TrebleClef: trebleClef,
		BassClef:   bassClef,
	}
}

var AllNotes = []Note{
	note(0, "c", false, true),
	note(2, "d", false, true),
	note(4, "e", false, true),
	note(5, "f", false, true),
	note(7, "g", false, true),
	note(9, "a", false, true),
	note(11, "b", false, true),
	note(12, "c", false, true),
	note(14, "d", false, true),
	note(16, "e", false, true),
	note(17, "f", false, true),
	note(19, "g", false, true),
	note(21, "a", true, true),
	note(23, "b", true, true),
	note(24, "c", true, true),
	note(26, "d", true, true),
	note(28, "e", true, true),
	note(29, "f", true, false),
	note(31, "g", true, false),
	note(33, "a", true, false),
	note(35, "b", true, false),
	note(36, "c", true, false),
	note(38, "d", true, false),
	note(40, "e", true, false),
	note(41, "f", true, false),
	note(43, "g", true, false),
	note(45, "a", true, false),
	note(47, "b", true, false),
	note(48, "c", true, false),
}
