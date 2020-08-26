package notes

func GenerateIntervals(noteList []Note, startIndex int, endIndex int, maxDist int) []Interval {
	var intervals []Interval
	for i := startIndex; i < endIndex; i++ {
		for j := i + 1; j < endIndex && noteList[j].ToneIndex()-noteList[i].ToneIndex() <= maxDist; j++ {
			first := noteList[i]
			second := noteList[j]

			if first.BassClef && second.BassClef {
				interval := Interval{FirstNote: first, SecondNote: second}
				interval.FirstNote.TrebleClef = false
				interval.SecondNote.TrebleClef = false
				intervals = append(intervals, interval)
			}

			if first.BassClef && second.TrebleClef {
				interval := Interval{FirstNote: first, SecondNote: second}
				interval.FirstNote.TrebleClef = false
				interval.SecondNote.BassClef = false
				intervals = append(intervals, interval)
			}

			if first.TrebleClef && second.TrebleClef {
				interval := Interval{FirstNote: first, SecondNote: second}
				interval.FirstNote.BassClef = false
				interval.SecondNote.BassClef = false
				intervals = append(intervals, interval)
			}
		}
	}

	return intervals
}
