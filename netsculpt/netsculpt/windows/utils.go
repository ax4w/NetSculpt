package windows

import (
	"github.com/rivo/tview"
	"strconv"
)

func newText(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
}

func prefixInputValidator(textToCheck string, lastChar rune) bool {
	if textToCheck == "-" {
		return true
	}
	v, err := strconv.Atoi(textToCheck)
	if err != nil {
		return false
	}
	return v >= 1 && v <= 30
}

func ipInputValidator(textToCheck string, lastChar rune) bool {
	if textToCheck == "-" {
		return true
	}
	v, err := strconv.Atoi(textToCheck)
	if err != nil {
		return false
	}
	return v >= 0 && v <= 255
}

const inputWidth = 15
