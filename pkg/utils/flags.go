package utils

import (
	"fmt"
	"github.com/lsierant/notes-gen/pkg/notes"
)

func FilterScales(scaleFlag string, accidentals int) ([]notes.Scale, error) {
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
