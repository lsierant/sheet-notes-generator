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
  \key des \major

  c'1 cis'1 ces'1 d'1 dis'1 des'1 e'1 eis'1 es'1 f'1 fis'1 fes'1 g'1 gis'1 ges'1 a'1 ais'1 as'1 b'1 bis'1 bes'1
}

lower = {
    \once \override Staff.TimeSignature #'transparent = ##t
    \key c \major
    \clef bass

	s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1 s1
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
