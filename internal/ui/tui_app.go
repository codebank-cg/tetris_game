package ui

import (
	"github.com/rivo/tview"
)

// NewApplication creates a new tview application.
func NewApplication() *tview.Application {
	return tview.NewApplication()
}

// NewFlex creates a new flex container.
func NewFlex() *tview.Flex {
	return tview.NewFlex()
}

// NewTextView creates a new text view.
func NewTextView() *tview.TextView {
	return tview.NewTextView()
}

// NewBox creates a new box.
func NewBox() *tview.Box {
	return tview.NewBox()
}
