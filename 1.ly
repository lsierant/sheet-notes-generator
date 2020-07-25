
\version "2.14.1"
\include "english.ly"
\include "lilypond-book-preamble.ly" 

\paper{
  indent=0\mm
  line-width=120\mm
  oddFooterMarkup=##f
  oddHeaderMarkup=##f
  bookTitleMarkup = ##f
  scoreTitleMarkup = ##f
}

\score {
  \new Staff {
    \once \override Staff.TimeSignature #'transparent = ##t
    \key d \major
    \clef treble
      <cs' g'>4
  }
}