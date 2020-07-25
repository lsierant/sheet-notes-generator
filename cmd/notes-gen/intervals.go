package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/lsierant/notes-gen/pkg/lilypond"
	"github.com/lsierant/notes-gen/pkg/notes"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	tmpDir := flag.String("tmpDir", "tmp", "temp directory for generating lilypond images")
	imageDir := flag.String("imageDir", "images", "destination directory for storing generated images")
	deckFilePath := flag.String("deckFilePath", "deck.csv", "path to generated deck file")

	flag.Parse()

	intervals := notes.GenerateIntervals(notes.AllNotes, 0, len(notes.AllNotes), 12)
	for i := 0; i < len(intervals); i++ {
		fmt.Printf("%s %s -> %s %s\n", intervals[i].FirstNote, intervals[i].FirstNote.LilypondSymbol(), intervals[i].SecondNote, intervals[i].SecondNote.LilypondSymbol())
	}

	renderer := lilypond.Renderer{WorkingDir: *tmpDir}

	for i := 0; i < len(intervals); i++ {
		err := renderIntervalAndWriteFile(renderer, intervals[i], fmt.Sprintf("%s/%s", *imageDir, fmt.Sprintf("%s.png", intervalFileName(intervals[i]))))
		if err != nil {
			log.Fatal(err)
		}
	}

	deckFileContent := prepareDeck(intervals)

	err := ioutil.WriteFile(*deckFilePath, []byte(deckFileContent), 0660)
	if err != nil {
		fmt.Printf("error writing file: %v", err)
	}

	fmt.Println("Done...")
}

func prepareDeck(intervals []notes.Interval) string {
	deckLines := make([]string, 0)

	for i := 0; i < len(intervals); i++ {
		deckLines = append(deckLines, deckLine(intervals[i]))
	}

	sort.Strings(deckLines)
	return strings.Join(deckLines, "\n")
}

func deckLine(interval notes.Interval) string {
	frontText := fmt.Sprintf("<img src=\"\"%s.png\"\">", intervalFileName(interval))
	backText := fmt.Sprintf("%s (%d), %s -> %s", interval.Name(), interval.Distance(), interval.FirstNote.Name, interval.SecondNote.Name)

	return fmt.Sprintf(`"%s";"%s"`, frontText, backText)
}

func renderIntervalAndWriteFile(renderer lilypond.Renderer, interval notes.Interval, intervalFilePath string) error {
	png, err := lilypond.RenderIntervalImage(&renderer, interval)

	if err != nil {
		return fmt.Errorf("failed to render lilypond image: %v", err)
	}

	err = ioutil.WriteFile(intervalFilePath, png, os.FileMode(0660))
	if err != nil {
		return fmt.Errorf("failed to write png file: %v", err)
	}

	log.Printf("Rendered file: %s\n", intervalFilePath)

	return nil
}

func intervalFileName(interval notes.Interval) string {
	intervalFileName := fmt.Sprintf("%s_%s", interval.FirstNote, interval.SecondNote)
	md5Hash := fmt.Sprintf("%x", md5.Sum([]byte(intervalFileName)))
	return fmt.Sprintf("ng-%s-%s", md5Hash, intervalFileName)
}
