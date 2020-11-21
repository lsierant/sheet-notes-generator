package notes

import (
	"fmt"
	"strings"
)

type Chord struct {
	Scale    Scale
	Notes    []Note
	RootNote Note
	Type     ChordType
}

func (c Chord) Name() string {
	name := c.RootNote.NameWithSharpFlatModifier()

	switch c.Type {
	case ChordTypeMajorTriad:
		name += " maj"
	case ChordTypeMinorTriad:
		name += " min"
	case ChordTypeDiminishedTriad:
		name += " dim"
	case ChordTypeAugmentedTriad:
		name += " aug"
	case ChordTypeMinorSeventh:
		name += " min7"
	case ChordTypeDominantSeventh, ChordTypeMajorSeventh:
		name += " maj7"
	case ChordTypeHalfDiminishedSeventh:
		name += " halfdim7"
	case ChordTypeDiminishedSeventh:
		name += " dim7"
	}

	return name
}

func (c Chord) RomanNumeral() string {
	if c.isTriad() {
		return c.Scale.Degree(c.RootNote.BaseName).RomanNumeralTriad
	}

	return c.Scale.Degree(c.RootNote.BaseName).RomanNumeralSeventh
}

func (c Chord) isTriad() bool {
	switch c.Type {
	case ChordTypeMajorTriad, ChordTypeMinorTriad, ChordTypeDiminishedTriad, ChordTypeAugmentedTriad:
		return true
	}

	return false
}

type ChordOnClefs struct {
	TrebleClefNotes []Note
	BassClefNotes   []Note
}

type ChordType int

const (
	ChordTypeMinorTriad ChordType = iota
	ChordTypeMajorTriad
	ChordTypeDiminishedTriad
	ChordTypeAugmentedTriad
	ChordTypeMinorSeventh
	ChordTypeMajorSeventh
	ChordTypeDominantSeventh
	ChordTypeDiminishedSeventh
	ChordTypeHalfDiminishedSeventh
)

type ChordQuality int

const (
	ChordQualityMinor ChordQuality = iota
	ChordQualityMajor
	ChordQualityDiminished
)

func GenerateAllDiatonicTriadsInScale(scale Scale) []Chord {
	var resultChords []Chord

	scaleNotes := ApplyScale(AllNotes, scale)
	for i := 0; i < len(scaleNotes); i++ {
		fmt.Printf("%d\t notes[%d]=%+v, %s\n", i, scaleNotes[i].BaseNoteIndex, scaleNotes[i], scaleNotes[i].LilypondSymbol())
	}

	triadNotes := noteNamesOfAllTriads()

	for i := 0; i < len(triadNotes); i++ {
		chords := generateAllTriadsForChord(scaleNotes, triadNotes[i], scale)
		fmt.Printf("triadNotes: %s\n", triadNotes)

		for chordIdx := 0; chordIdx < len(chords); chordIdx++ {
			fmt.Printf("chord: %s (%s)\n", chords[chordIdx].RootNote, chords[chordIdx].Type)
			chordsOnClefs := generateChordsOnClefs(chords[chordIdx], scale)
			resultChords = append(resultChords, chordsOnClefs...)
		}
	}

	return resultChords
}

func GenerateAllDiatonicSeventhsInScaleWithoutFifths(scale Scale) []Chord {
	var resultChords []Chord

	scaleNotes := ApplyScale(AllNotes, scale)
	for i := 0; i < len(scaleNotes); i++ {
		fmt.Printf("%d\t notes[%d]=%+v, %s\n", i, scaleNotes[i].BaseNoteIndex, scaleNotes[i], scaleNotes[i].LilypondSymbol())
	}

	seventhNotes := noteNamesOfAllSeventhsWithoutFifth()

	for i := 0; i < len(seventhNotes); i++ {
		chords := generateAllTriadsForChord(scaleNotes, seventhNotes[i], scale)

		fmt.Printf("all chords size = %d\n", len(chords))
		filteredChords := filterChordsUsingPredicate(chords, chordNotesHasMinimumDistanceBetweenNotes(3))

		fmt.Printf("filtered size = %d\n", len(filteredChords))

		chords = filteredChords

		fmt.Printf("triadNotes: %s\n", seventhNotes)

		for chordIdx := 0; chordIdx < len(chords); chordIdx++ {
			fmt.Printf("chord: %s (%s)\n", chords[chordIdx].RootNote, chords[chordIdx].Type)
			chordsOnClefs := generateChordsOnClefs(chords[chordIdx], scale)
			resultChords = append(resultChords, chordsOnClefs...)
		}
	}

	return resultChords
}

func filterChordsUsingPredicate(chords []Chord, p func(Chord) bool) []Chord {
	var filteredChords []Chord
	for i := 0; i < len(chords); i++ {
		if p(chords[i]) {
			filteredChords = append(filteredChords, chords[i])
		}
	}

	return filteredChords
}

func chordNotesHasMinimumDistanceBetweenNotes(minDistance int) func(Chord) bool {
	return func(chord Chord) bool {
		for j := 0; j < len(chord.Notes); j++ {
			for k := j + 1; k < len(chord.Notes); k++ {
				interval := Interval{FirstNote: chord.Notes[j], SecondNote: chord.Notes[k]}
				if interval.Distance() < minDistance {
					return false
				}
			}
		}

		return true
	}
}

func ChordToChordOnClefs(chord Chord) ChordOnClefs {
	chordOnClefs := ChordOnClefs{}
	for i := 0; i < len(chord.Notes); i++ {
		if chord.Notes[i].TrebleClef && chord.Notes[i].BassClef {
			panic(fmt.Errorf("both clefs set in note: %v", chord.Notes[i]))
		}
		if chord.Notes[i].TrebleClef {
			chordOnClefs.TrebleClefNotes = append(chordOnClefs.TrebleClefNotes, chord.Notes[i])
		}
		if chord.Notes[i].BassClef {
			chordOnClefs.BassClefNotes = append(chordOnClefs.BassClefNotes, chord.Notes[i])
		}
	}

	return chordOnClefs
}

func generateChordsOnClefs(chord Chord, scale Scale) []Chord {
	var resultingChords []Chord
	var recursive func(noteIdx int, trebleClef bool, currentNotes []Note)
	recursive = func(noteIdx int, trebleClef bool, currentNotes []Note) {
		if noteIdx > len(chord.Notes) {
			return
		}

		if len(currentNotes) == 3 {
			newChord := chord
			newChord.Notes = currentNotes
			newChord.Scale = scale
			resultingChords = append(resultingChords, newChord)
			return
		}

		chordNotes := chord.Notes[noteIdx].NotesOnClefs(trebleClef)
		for i := 0; i < len(chordNotes); i++ {
			recursive(noteIdx+1, chordNotes[i].TrebleClef, append(currentNotes, chordNotes[i]))
		}
	}

	recursive(0, true, make([]Note, 0))

	return resultingChords
}

func generateAllTriadsForChord(noteList []Note, triadNotes string, scale Scale) []Chord {
	chordNotes := filterNotesByChordNotes(noteList, triadNotes)
	currentToneIdx := len(chordNotes) - 1
	resultChords := make([]Chord, 0)
	rootNoteLetter := string(triadNotes[0])
	var rootNote Note
	for ; currentToneIdx >= 2; currentToneIdx-- {
		sopranoNote := chordNotes[currentToneIdx]
		if sopranoNote.BaseName == rootNoteLetter {
			rootNote = sopranoNote
		}

		altoNote := chordNotes[currentToneIdx-1]
		if altoNote.BaseName == rootNoteLetter {
			rootNote = altoNote
		}

		tenorNote := chordNotes[currentToneIdx-2]
		if tenorNote.BaseName == rootNoteLetter {
			rootNote = tenorNote
		}

		degree := scale.Degree(rootNote.BaseName)
		chord := Chord{
			Scale:    scale,
			Notes:    []Note{sopranoNote, altoNote, tenorNote},
			RootNote: rootNote,
			Type:     degree.TriadType,
		}

		//fmt.Printf("chord: %s: %s,%s,%s: %s, %s, %s: %s, %s\n", triadNotes, tenorNote.LilypondSymbol(), altoNote.LilypondSymbol(), sopranoNote.LilypondSymbol(),
		//	Interval{tenorNote, altoNote, scale}.Name(), Interval{altoNote, sopranoNote, scale}.Name(), Interval{tenorNote, sopranoNote, scale}.Name(), chord.Type, fname)

		resultChords = append(resultChords, chord)
	}

	return resultChords
}

func filterNotesByChordNotes(noteList []Note, chord string) []Note {
	var filtered []Note
	for i := 0; i < len(noteList); i++ {
		if strings.Contains(chord, noteList[i].BaseName) {
			filtered = append(filtered, noteList[i])
		}
	}

	return filtered
}

func noteNamesOfAllTriads() []string {
	var triads []string
	cScale := "cdefgabcdefgab"
	for i := 0; i < 7; i++ {
		triads = append(triads, fmt.Sprintf("%c%c%c", cScale[i], cScale[i+2], cScale[i+4]))
	}

	return triads
}

func NoteNamesOfAllTriads() []string {
	var triads []string
	cScale := "cdefgabcdefgab"
	for i := 0; i < 7; i++ {
		triads = append(triads, fmt.Sprintf("%c%c%c", cScale[i], cScale[i+2], cScale[i+4]))
	}

	return triads
}

func noteNamesOfAllSeventhsWithoutFifth() []string {
	return []string{
		"ceb",
		"dfc",
		"egd",
		"fae",
		"gbf",
		"acg",
		"bda",
	}
}
