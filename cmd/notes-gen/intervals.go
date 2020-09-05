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

	flag.Parse()

	scales, err := filterScales(*scaleFlag, *accidentals)
	if err != nil {
		log.Fatal(err)
	}

	intervals := generateIntervals(scales)

	renderer := lilypond.Renderer{WorkingDir: *tmpDir}

	deckFileContent := prepareDeck(intervals)

	err = ioutil.WriteFile(*deckFilePath, []byte(deckFileContent), 0660)
	if err != nil {
		log.Fatalf("errors while rendering file:\n%v", err)
	}

	if htmlFilePath != nil && *htmlFilePath != "" {
		htmlFileContent := prepareHtml(intervals)

		err = ioutil.WriteFile(*htmlFilePath, []byte(htmlFileContent), 0660)
		if err != nil {
			log.Fatalf("errors while rendering html file:\n%v", err)
		}
	}

	err = utils.RunInParallel(ctx, len(intervals), *parallel, func(idx int) error {
		intervalFileName := fmt.Sprintf("%s/%s", *imageDir, fmt.Sprintf("%s.png", intervalFileName(intervals[idx])))
		if _, err := os.Stat(intervalFileName); err == nil {
			fmt.Printf("Skipping rendering: %s\n", intervalFileName)
			return nil
		}

		return renderIntervalAndWriteFile(ctx, renderer, intervals[idx], intervalFileName)
	})

	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}

	fmt.Println("Done...")
}

func generateIntervals(scales []notes.Scale) []notes.Interval {
	var intervals []notes.Interval

	for _, scale := range scales {
		notesInScale := notes.ApplyScale(notes.AllNotes, scale)
		intervalsInScale := notes.GenerateIntervals(notesInScale, 0, len(notesInScale), 12)
		fmt.Printf("Scale: %s\n", scale.Name)
		for i := 0; i < len(intervalsInScale); i++ {
			intervalsInScale[i].Scale = scale
			fmt.Printf("%s %s -> %s %s\n", intervalsInScale[i].FirstNote.NameWithModifier(), intervalsInScale[i].FirstNote.LilypondSymbol(), intervalsInScale[i].SecondNote.NameWithModifier(), intervalsInScale[i].SecondNote.LilypondSymbol())
		}
		intervals = append(intervals, intervalsInScale...)
	}

	return intervals
}

func filterScales(scaleFlag string, accidentals int) ([]notes.Scale, error) {
	var scales []notes.Scale
	if scaleFlag == "major" {
		for _, scale := range notes.ScaleMap {
			if scale.Mode == notes.ScaleModeMajor && len(scale.NotesModified) <= accidentals {
				scales = append(scales, scale)
			}
		}
	} else {
		scale, ok := notes.ScaleMap[scaleFlag]
		if !ok {
			return nil, fmt.Errorf("invalid scale: %s", scaleFlag)
		}
		scales = append(scales, scale)
	}

	if len(scales) == 0 {
		return nil, fmt.Errorf("no scales found for: %s", scaleFlag)
	}

	return scales, nil
}

func prepareDeck(intervals []notes.Interval) string {
	deckLines := make([]string, 0)

	for i := 0; i < len(intervals); i++ {
		deckLines = append(deckLines, deckLine(intervals[i]))
	}

	sort.Strings(deckLines)
	return strings.Join(deckLines, "\n")
}

func prepareHtml(intervals []notes.Interval) string {
	sort.Slice(intervals, func(i int, j int) bool {
		if len(intervals[i].Scale.NotesModified) != len(intervals[j].Scale.NotesModified) {
			return len(intervals[i].Scale.NotesModified) < len(intervals[j].Scale.NotesModified)
		}

		if intervals[i].Scale.Name != intervals[j].Scale.Name {
			return intervals[i].Scale.Name < intervals[j].Scale.Name
		}

		if intervals[i].Distance() != intervals[j].Distance() {
			return intervals[i].Distance() < intervals[j].Distance()
		}

		if intervals[i].FirstNote.ToneIndex() != intervals[j].FirstNote.ToneIndex() {
			return intervals[i].FirstNote.ToneIndex() < intervals[j].FirstNote.ToneIndex()
		}

		return intervals[i].SecondNote.ToneIndex() != intervals[j].SecondNote.ToneIndex()
	})

	lines := make([]string, 0)
	for i := 0; i < len(intervals); i++ {
		imageSrc := fmt.Sprintf("images/%s.png", intervalFileName(intervals[i]))
		lines = append(lines, fmt.Sprintf(`
<div><img style="vertical-align:middle" src="%s" width="150"><span style="margin-left: 30pt;">%s</span></div><br><hr>`, imageSrc, backText(intervals[i])))
	}

	return fmt.Sprintf(`
<html>
<body>
%s
</body>
</html>
`, strings.Join(lines, "\n"))
}

func deckLine(interval notes.Interval) string {
	return fmt.Sprintf(`"%s";"%s"`, frontText(interval), backText(interval))
}

func backText(interval notes.Interval) string {
	return fmt.Sprintf("%s (%d), %s -> %s", interval.Name(), interval.Distance(), interval.FirstNote.NameWithModifier(), interval.SecondNote.NameWithModifier())
}

func frontText(interval notes.Interval) string {
	return fmt.Sprintf("<img src=\"\"%s.png\"\">", intervalFileName(interval))
}

func renderIntervalAndWriteFile(ctx context.Context, renderer lilypond.Renderer, interval notes.Interval, intervalFilePath string) error {
	png, err := lilypond.RenderIntervalImage(ctx, &renderer, interval)

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
	scaleName := strings.ReplaceAll(interval.Scale.Name, " ", "_")
	intervalFileName := fmt.Sprintf("%s_%s_%s", scaleName, interval.FirstNote, interval.SecondNote)
	md5Hash := fmt.Sprintf("%x", md5.Sum([]byte(intervalFileName)))
	return fmt.Sprintf("ng-%s-%s", md5Hash, intervalFileName)
}
