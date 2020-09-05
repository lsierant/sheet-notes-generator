package lilypond

import (
	"bytes"
	"fmt"
	"text/template"
)

func parseAndRenderTextTemplate(templateName string, tpl string, data interface{}) (string, error) {
	subjectTemplate, err := template.New(templateName).Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("error parsing %s template: %v", templateName, err)
	}

	bufferString := bytes.NewBufferString("")
	err = subjectTemplate.Execute(bufferString, data)
	if err != nil {
		return "", fmt.Errorf("error rendering subject template: %v", err)
	}
	return bufferString.String(), nil
}

var trebleOnlyIntervalTemplate = `
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
  \key {{.Scale.LilypondSymbol}}

  
  <{{.FirstNote.LilypondSymbol}} {{.SecondNote.LilypondSymbol}}>4
}

lower = {
    \once \override Staff.TimeSignature #'transparent = ##t
	\key {{.Scale.LilypondSymbol}}

    \clef bass
    
	s4
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

var bassOnlyIntervalTemplate = `
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
  \key {{.Scale.LilypondSymbol}}

  
  s4
}

lower = {
    \once \override Staff.TimeSignature #'transparent = ##t
	\key {{.Scale.LilypondSymbol}}

    \clef bass
    <{{.FirstNote.LilypondSymbol}} {{.SecondNote.LilypondSymbol}}>4
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

var bassAndTrebleIntervalTemplate = `
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
  \key {{.Scale.LilypondSymbol}}
  
  {{.SecondNote.LilypondSymbol}}4
}

lower = {
    \once \override Staff.TimeSignature #'transparent = ##t
	\key {{.Scale.LilypondSymbol}}

    \clef bass

  {{.FirstNote.LilypondSymbol}}4
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
