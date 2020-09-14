package lilypond

import (
	"fmt"
	"strings"
)

type MultipleChords struct {
	Scale  string
	Chords []SingleChord
}
type SingleChord struct {
	TrebleRaw   string
	TrebleNotes []string
	BassRaw     string
	BassNotes   []string
}

func (c SingleChord) Bass() string {
	if c.BassRaw != "" {
		return c.BassRaw
	}

	return fmt.Sprintf("<%s>4", strings.Join(c.BassNotes, " "))
}

func (c SingleChord) Treble() string {
	if c.TrebleRaw != "" {
		return c.TrebleRaw
	}

	return fmt.Sprintf("<%s>4", strings.Join(c.TrebleNotes, " "))
}

var chordTemplate = `
\version "2.14.1"
\include "lilypond-book-preamble.ly" 

\paper{
  indent=0\mm
  line-width=120\mm
  oddFooterMarkup=##f
  oddHeaderMarkup=##f
  bookTitleMarkup = ##f
  scoreTitleMarkup = ##f
}

upper = {
  \clef treble
  \once \override Staff.TimeSignature #'transparent = ##t
  \key {{.Scale}}

{{range .Chords}}
	{{ .Treble }}
{{end}}
}

lower = {
    \once \override Staff.TimeSignature #'transparent = ##t
	\key {{.Scale}}

    \clef bass
    
{{range .Chords}}
	{{ .Bass }}
{{end}}
}

\score {
  \new PianoStaff
  <<
    \new Staff = "upper" \upper
    \new Staff = "lower" \lower
  >>
  \layout { }
  \midi { }
}
`
