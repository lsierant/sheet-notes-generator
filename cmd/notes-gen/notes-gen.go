package main

import (
	"bytes"
	"fmt"
	"github.com/lsierant/notes-gen/pkg/lilypond"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"text/template"
)

func main() {
	renderer := lilypond.Renderer{WorkingDir: "tmp"}

	trebleClefNotes := []string{"c'", "cs'", "d'", "ds'", "e'", "f'", "fs'", "g'", "gs'", "a'", "as'", "b'",
								"c''", "cs''", "d''", "ds''", "e''", "f''", "fs''", "g''", "gs''", "a''", "as''", "b''"}

	wg := sync.WaitGroup{}
	wg.Add(len(trebleClefNotes) - 1)
	for i := 1; i < len(trebleClefNotes); i++ {
		go func(index int) {
			defer wg.Done()

			firstNote := "c'"

			interval := Interval{
				FirstNote:  firstNote,
				SecondNote: trebleClefNotes[index],
				Name:       fmt.Sprintf("%s_%s", firstNote, trebleClefNotes[index]),
			}
			source, err := parseAndRenderTextTemplate("interval", intervalTemplate, interval)
			if err != nil {
				log.Fatalf("failed to render interval template: %v", err)
			}

			fmt.Printf("Rendering interval: %+v\n", interval)

			png, err := renderer.RenderPNG(source)
			if err != nil {
				log.Fatalf("failed to render PNG from source: %s: %v", source, err)
			}

			err = ioutil.WriteFile(fmt.Sprintf("tmp/intervals/%s.png", interval.Name), png, os.FileMode(0660))
			if err != nil {
				log.Fatalf("failed to write png file: %v", err)
			}
		}(i)
	}

	fmt.Println("Waiting...")
	wg.Wait()
	fmt.Println("Done...")
}

type Interval struct {
	FirstNote  string
	SecondNote string
	Name       string
}

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

var intervalTemplate = `
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
      <{{.FirstNote}} {{.SecondNote}}>4
  }
}`
