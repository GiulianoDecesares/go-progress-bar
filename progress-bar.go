package progressbar

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/buger/goterm"
	"github.com/fatih/color"
)

const (
	maxBarLength = 50
	minBarLength = 1
)

var barColor = color.New(color.FgGreen).SprintFunc()

type Bar struct {
	lastPercent float64 // last recorded progress percentage
	total       int64   // total value for progress
	graph       string  // the fill value for progress bar
	background  string  // the background value for the progress bar
	processId   string  // the process the progress bar is showing

	currentBarLength int
}

func NewProgressBar(start int64, total int64, processName string) *Bar {
	bar := &Bar{
		lastPercent:      -1, // To draw the first time with 0% progress
		total:            total,
		processId:        processName,
		graph:            "█",
		background:       "▓",
		currentBarLength: maxBarLength,
	}

	return bar
}

func (bar *Bar) Update(currentRawProgress int64) {
	currentPercent := bar.getPercent(currentRawProgress)

	if currentPercent != bar.lastPercent { // Avoid redrawing UI with same percent
		bar.lastPercent = currentPercent
		fmt.Print(bar.getDefaultLayout(bar.lastPercent))
	}
}

func (bar *Bar) Finish() {
	fmt.Println()
}

func (bar *Bar) getPercent(currentRawProgress int64) float64 {
	return (float64(currentRawProgress) / float64(bar.total)) * float64(100)
}

func (bar *Bar) getDefaultLayout(percentProgress float64) string {
	layout := bar.GetFullLayout(percentProgress)

	// Calculate correction if necessary
	if correction := utf8.RuneCountInString(layout) - goterm.Width(); correction > 0 {
		bar.currentBarLength -= correction

		if bar.currentBarLength < 0 {
			bar.currentBarLength = minBarLength
		}

		layout = bar.GetFullLayout(percentProgress)
	}

	return layout
}

func (bar *Bar) GetBarLayout(percentProgress float64) string {
	normalizedPercent := (int(percentProgress) * bar.currentBarLength) / 100

	fullBarPortion := strings.Repeat(bar.graph, int(normalizedPercent))
	emptyBarPortion := strings.Repeat(bar.background, bar.currentBarLength-int(normalizedPercent))

	return fullBarPortion + emptyBarPortion
}

func (bar *Bar) GetFullLayout(percentProgress float64) string {
	barString := bar.GetBarLayout(percentProgress)

	const layoutTemplate string = "\r%s [%s] %.2f%%"

	return fmt.Sprintf(layoutTemplate, barColor(barString), bar.processId, percentProgress)
}
