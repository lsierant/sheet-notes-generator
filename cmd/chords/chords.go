package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/lsierant/notes-gen/pkg/lilypond"
	"github.com/lsierant/notes-gen/pkg/notes"
	"github.com/lsierant/notes-gen/pkg/utils"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
)

var debug = true

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	utils.HandleSignals(func(code os.Signal) {
		log.Printf("Received signal %d", code)
		cancel()
	})

	tmpDir := flag.String("tmpDir", "tmp", "temp directory for generating lilypond images")
	imageDir := flag.String("imageDir", "images", "destination directory for storing generated images")
	deckFilePath := flag.String("deckFilePath", "deck.csv", "path to generated deck file")
	htmlFilePath := flag.String("htmlFilePath", "", "path to generated html file with all images")
	parallel := flag.Int("parallel", runtime.NumCPU(), "level of parallelism, defaults to number of CPUs")
	scaleFlag := flag.String("scale", "c major", `scale to use, e.g. "c flat major", "d minor", "c sharp minor", default: "c major"`)
	accidentals := flag.Int("accidentals", 7, "filter scales up to given number of accidentals")
	triads := flag.Bool("triads", true, "generate triads")
	//sevenths := flag.Bool("sevenths", false, "generate triads")
	onePager := flag.Bool("onePager", false, "generate one pager instead of deck")

	flag.Parse()

	scales, err := utils.FilterScales(*scaleFlag, *accidentals)
	if err != nil {
		log.Fatal(err)
	}

	renderer := lilypond.Renderer{WorkingDir: *tmpDir}

	if *onePager {
		renderAllDiatonicTriadsOnOnePage(ctx, renderer, *imageDir, scales)
	} else {
		if *triads {
			triads := generateAllTriadsInScales(scales)

			if deckFilePath != nil && *deckFilePath != "" {
				deckFileContent := prepareDeck(triads)
				err = ioutil.WriteFile(*deckFilePath, []byte(deckFileContent), 0660)
				if err != nil {
					log.Fatalf("errors while rendering file:\n%v", err)
				}
			}

			if htmlFilePath != nil && *htmlFilePath != "" {
				htmlFileContent := prepareHtml(triads)

				err = ioutil.WriteFile(*htmlFilePath, []byte(htmlFileContent), 0660)
				if err != nil {
					log.Fatalf("errors while rendering html file:\n%v", err)
				}
			}
			renderAllDiatonicTriadsAsSeparateImages(ctx, renderer, *imageDir, *parallel, triads)
		}
	}
}

func renderAllDiatonicTriadsOnOnePage(ctx context.Context, renderer lilypond.Renderer, destDir string, scales []notes.Scale) {
	for s := 0; s < len(scales); s++ {
		chords := notes.GenerateAllDiatonicTriadsInScale(scales[s])

		chordFilePath := fmt.Sprintf("%s/ng-chord-all-%s.png", destDir, scales[s].Name)
		fmt.Println(chordFilePath)
		if _, err := os.Stat(chordFilePath); err == nil {
			fmt.Printf("Skipping rendering: %s\n", chordFilePath)
			continue
		}

		multipleChords := lilypond.MultipleChords{
			Scale:  scales[s].LilypondSymbol,
			Chords: convertToLilypondChords(chords),
		}

		err := renderChordAndWriteFile(ctx, renderer, multipleChords, chordFilePath)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Done...")
}

func generateAllTriadsInScales(scales []notes.Scale) []notes.Chord {
	var chords []notes.Chord
	for i := 0; i < len(scales); i++ {
		chords = append(chords, notes.GenerateAllDiatonicTriadsInScale(scales[i])...)
	}

	return chords
}

func renderAllDiatonicTriadsAsSeparateImages(ctx context.Context, renderer lilypond.Renderer, destDir string, parallel int, chords []notes.Chord) {
	err := utils.RunInParallel(ctx, len(chords), parallel, func(idx int) error {
		chord := convertToLilypondChord(chords[idx])
		multipleChords := lilypond.MultipleChords{
			Scale:  chords[idx].Scale.LilypondSymbol,
			Chords: []lilypond.SingleChord{chord},
		}

		chordFilePath := chordFilePath(destDir, chords[idx])
		fmt.Println(chordFilePath)

		if _, err := os.Stat(chordFilePath); err == nil {
			fmt.Printf("Skipping rendering: %s\n", chordFilePath)
			return nil
		}

		err := renderChordAndWriteFile(ctx, renderer, multipleChords, chordFilePath)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[%d]%v, ", idx, chords[idx])

		return nil
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("\n")

	fmt.Println("Done...")
}

func chordFilePath(imageDir string, chord notes.Chord) string {
	return fmt.Sprintf("%s/%s", imageDir, fmt.Sprintf("%s.png", chordFileName(chord)))
}

func convertToLilypondChords(chord []notes.Chord) []lilypond.SingleChord {
	var result []lilypond.SingleChord
	for i := 0; i < len(chord); i++ {
		result = append(result, convertToLilypondChord(chord[i]))
	}

	return result
}

func convertToLilypondChord(chord notes.Chord) lilypond.SingleChord {
	chordOnClefs := notes.ChordToChordOnClefs(chord)
	lilypondChord := lilypond.SingleChord{}

	for i := 0; i < len(chordOnClefs.TrebleClefNotes); i++ {
		lilypondChord.TrebleNotes = append(lilypondChord.TrebleNotes, chordOnClefs.TrebleClefNotes[i].LilypondSymbol())
	}

	if len(chordOnClefs.TrebleClefNotes) == 0 {
		lilypondChord.TrebleRaw = "s4"
	}

	for i := 0; i < len(chordOnClefs.BassClefNotes); i++ {
		lilypondChord.BassNotes = append(lilypondChord.BassNotes, chordOnClefs.BassClefNotes[i].LilypondSymbol())
	}

	if len(chordOnClefs.BassClefNotes) == 0 {
		lilypondChord.BassRaw = "s4"
	}

	return lilypondChord
}

func renderChordAndWriteFile(ctx context.Context, renderer lilypond.Renderer, chord lilypond.MultipleChords, chordFilePath string) error {
	png, err := lilypond.RenderChordImage(ctx, &renderer, chord)

	if err != nil {
		return fmt.Errorf("failed to render lilypond image: %v", err)
	}

	if debug {
		source, err := lilypond.RenderChordSource(ctx, &renderer, chord)
		if err != nil {
			return fmt.Errorf("failed to render lilypond source: %v", err)
		}

		if err = ioutil.WriteFile(chordFilePath+".ll", []byte(source), os.FileMode(0660)); err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(chordFilePath, png, os.FileMode(0660))
	if err != nil {
		return fmt.Errorf("failed to write png file: %v", err)
	}

	log.Printf("Rendered file: %s\n", chordFilePath)

	return nil
}

func chordFileName(chord notes.Chord) string {
	scaleName := strings.ReplaceAll(chord.Scale.Name, " ", "_")
	var chordNotesArr []string
	for i := 0; i < len(chord.Notes); i++ {
		chordNotesArr = append(chordNotesArr, fmt.Sprintf("%s", chord.Notes[i]))
	}
	chordFileName := fmt.Sprintf("%s_%s", scaleName, strings.Join(chordNotesArr, "_"))
	md5Hash := fmt.Sprintf("%x", md5.Sum([]byte(chordFileName)))
	return fmt.Sprintf("ng-chord-%s-%s", md5Hash, chordFileName)
}

func prepareDeck(chords []notes.Chord) string {
	deckLines := make([]string, 0)

	for i := 0; i < len(chords); i++ {
		deckLines = append(deckLines, deckLine(chords[i]))
	}

	sort.Strings(deckLines)
	return strings.Join(deckLines, "\n")
}

func prepareHtml(chords []notes.Chord) string {
	sort.Slice(chords, func(i int, j int) bool {
		if len(chords[i].Scale.NotesModified) != len(chords[j].Scale.NotesModified) {
			return len(chords[i].Scale.NotesModified) < len(chords[j].Scale.NotesModified)
		}

		if chords[i].Scale.Name != chords[j].Scale.Name {
			return chords[i].Scale.Name < chords[j].Scale.Name
		}

		if chords[i].RootNote.BaseName != chords[j].RootNote.BaseName {
			return chords[i].RootNote.BaseName < chords[j].RootNote.BaseName
		}

		for n := 0; n < len(chords[i].Notes); n++ {
			if chords[i].Notes[n].ToneIndex() != chords[j].Notes[n].ToneIndex() {
				return chords[i].Notes[n].ToneIndex() < chords[j].Notes[n].ToneIndex()
			}
		}

		return false
	})

	lines := make([]string, 0)
	for i := 0; i < len(chords); i++ {
		imageSrc := fmt.Sprintf("images/%s.png", chordFileName(chords[i]))
		lines = append(lines, fmt.Sprintf(`
<div><img style="vertical-align:middle" src="%s" width="150"><span style="margin-left: 30pt;">%s</span></div><br><hr>`, imageSrc, backText(chords[i])))
	}

	return fmt.Sprintf(`
<?xml version="1.0" encoding="utf-8" ?>
<html>
<body>
%s
</body>
</html>
`, strings.Join(lines, "\n"))
}

func deckLine(chord notes.Chord) string {
	return fmt.Sprintf(`"%s";"%s"`, frontText(chord), backText(chord))
}

func backText(chord notes.Chord) string {
	return fmt.Sprintf("%s", chord.Name())
}

func frontText(chord notes.Chord) string {
	return fmt.Sprintf("<img src=\"\"%s.png\"\">", chordFileName(chord))
}
