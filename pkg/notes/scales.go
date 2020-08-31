package notes

type ScaleMode int

const (
	ScaleModeMajor         = 0
	ScaleModeMinorHarmonic = 1
)

type Scale struct {
	Note           string
	Mode           ScaleMode
	Modifier       NoteModifier
	NotesModified  []string
	Name           string
	LilypondSymbol string
}

var (
	CMajorScale = Scale{
		Note:           "c",
		LilypondSymbol: `c \major`,
		Name:           "c major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierNone,
		NotesModified:  []string{},
	}

	AMinorScale = Scale{
		Note:           "a",
		LilypondSymbol: `a \minor`,
		Name:           "a minor",
		Mode:           ScaleModeMinorHarmonic,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"g"},
	}

	GMajorScale = Scale{
		Note:           "g",
		LilypondSymbol: `g \major`,
		Name:           "g major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"f"},
	}

	EMinorScale = Scale{
		Note:           "e",
		LilypondSymbol: `e \minor`,
		Name:           "e minor",
		Mode:           ScaleModeMinorHarmonic,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"f", "d"},
	}

	DMajorScale = Scale{
		Note:           "d",
		LilypondSymbol: `d \major`,
		Name:           "d major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"c", "f"},
	}

	BMinorScale = Scale{
		Note:           "b",
		LilypondSymbol: `b \minor`,
		Name:           "b minor",
		Mode:           ScaleModeMinorHarmonic,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"c", "f", "a"},
	}

	AMajorScale = Scale{
		Note:           "a",
		LilypondSymbol: `a \major`,
		Name:           "a major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"c", "f", "g"},
	}

	FSharpMinorScale = Scale{
		Note:           "f",
		LilypondSymbol: `fis \minor`,
		Name:           "f sharp minor",
		Mode:           ScaleModeMinorHarmonic,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"c", "f", "g", "e"},
	}

	EMajorScale = Scale{
		Note:           "e",
		LilypondSymbol: `e \major`,
		Name:           "e major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"c", "f", "g", "d"},
	}

	CSharpMinorScale = Scale{
		Note:           "c",
		LilypondSymbol: `cis \minor`,
		Name:           "c sharp minor",
		Mode:           ScaleModeMinorHarmonic,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"c", "f", "g", "d", "b"},
	}

	BMajorScale = Scale{
		Note:           "b",
		LilypondSymbol: `b \major`,
		Name:           "b major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"c", "f", "g", "d", "a"},
	}

	FSharpMajor = Scale{
		Note:           "f",
		LilypondSymbol: `fis \major`,
		Name:           "f sharp major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierSharp,
		NotesModified:  []string{"f", "c", "g", "d", "a", "e"},
	}

	FMajorScale = Scale{
		Note:           "f",
		LilypondSymbol: `f \major`,
		Name:           "f major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierFlat,
		NotesModified:  []string{"b"},
	}

	BFlatMajorScale = Scale{
		Note:           "b",
		LilypondSymbol: `bes \major`,
		Name:           "b flat major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierFlat,
		NotesModified:  []string{"b", "e"},
	}

	EFlatMajorScale = Scale{
		Note:           "e",
		LilypondSymbol: `es \major`,
		Name:           "e flat major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierFlat,
		NotesModified:  []string{"b", "e", "a"},
	}

	AFlatMajorScale = Scale{
		Note:           "a",
		LilypondSymbol: `as \major`,
		Name:           "a flat major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierFlat,
		NotesModified:  []string{"b", "e", "a", "d"},
	}

	DFlatMajorScale = Scale{
		Note:           "a",
		LilypondSymbol: `des \major`,
		Name:           "d flat major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierFlat,
		NotesModified:  []string{"b", "e", "a", "d", "g"},
	}

	GFlatMajorScale = Scale{
		Note:           "g",
		LilypondSymbol: `ges \major`,
		Name:           "g flat major",
		Mode:           ScaleModeMajor,
		Modifier:       NoteModifierFlat,
		NotesModified:  []string{"b", "e", "a", "d", "g", "c"},
	}
	//DMinorScale = Scale{
	//	Note:          "d",
	//Name: "d minor",
	//	Mode:          ScaleModeMinorHarmonic,
	//	Modifier:      NoteModifierFlat,
	//	NotesModified: []string{"b"},
	//}

	//GMinorScale = Scale{
	//	Note:          "b",
	//Name: "g minor",
	//	Mode:          ScaleModeMinorHarmonic,
	//	Modifier:      NoteModifierFlat,
	//	NotesModified: []string{"b", "e"},
	//}
	//CMinorScale = Scale{
	//	Note:          "c",
	//Name: "c minor",
	//	Mode:          ScaleModeMinorHarmonic,
	//	Modifier:      NoteModifierFlat,
	//	NotesModified: []string{"b", "e", "a"},
	//}

	//FMinorScale = Scale{
	//	Note:          "f",
	//Name: "f minor",
	//	Mode:          ScaleModeMinorHarmonic,
	//	Modifier:      NoteModifierFlat,
	//	NotesModified: []string{"b", "e", "a", "d"},
	//}
	//BFlatMinorScale = Scale{
	//	Note:          "b",
	//Name: "b flat minor",
	//	Mode:          ScaleModeMinorHarmonic,
	//	Modifier:      NoteModifierFlat,
	//	NotesModified: []string{"b", "e", "a", "d", "g"},
	//}

	//
	//EFlatMinorScale = Scale{
	//	Note:          "e",
	//Name: "e flat major",
	//	Mode:          ScaleModeMinorHarmonic,
	//	Modifier:      NoteModifierFlat,
	//	NotesModified: []string{"b", "e", "a", "d", "g", "c"},
	//}
)

var ScaleMap = map[string]Scale{
	CMajorScale.Name:      CMajorScale,
	AMinorScale.Name:      AMinorScale,
	GMajorScale.Name:      GMajorScale,
	EMinorScale.Name:      EMinorScale,
	DMajorScale.Name:      DMajorScale,
	BMinorScale.Name:      BMinorScale,
	AMajorScale.Name:      AMajorScale,
	FSharpMinorScale.Name: FSharpMinorScale,
	EMajorScale.Name:      EMajorScale,
	CSharpMinorScale.Name: CSharpMinorScale,
	BMajorScale.Name:      BMajorScale,
	FSharpMajor.Name:      FSharpMajor,
	FMajorScale.Name:      FMajorScale,
	//DMinorScale.Name: DMinorScale,
	BFlatMajorScale.Name: BFlatMajorScale,
	//GMinorScale.Name: GMinorScale,
	EFlatMajorScale.Name: EFlatMajorScale,
	//CMinorScale.Name: CMinorScale,
	AFlatMajorScale.Name: AFlatMajorScale,
	//FMinorScale.Name: FMinorScale,
	DFlatMajorScale.Name: DFlatMajorScale,
	//BFlatMinorScale.Name: BFlatMinorScale,
	GFlatMajorScale.Name: GFlatMajorScale,
}

func ApplyScale(notes []Note, scale Scale) []Note {
	var notesInScale []Note
	for i := 0; i < len(notes); i++ {
		note := notes[i]
		note.Scale = scale
		for j := 0; j < len(scale.NotesModified); j++ {
			if note.BaseName == scale.NotesModified[j] {
				note.Modifier = scale.Modifier
			}
		}
		notesInScale = append(notesInScale, note)
	}

	return notesInScale
}
