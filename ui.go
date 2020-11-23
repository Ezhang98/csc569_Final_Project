package main

import (
    "fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"fmt"
)

import "fyne.io/fyne"

type diagonal struct {
}

func (d *diagonal) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := 0, 0
	for _, o := range objects {
		childSize := o.MinSize()

		w += childSize.Width
		h += childSize.Height
	}
	return fyne.NewSize(w, h)
}



func (d *diagonal) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, containerSize.Height - d.MinSize(objects).Height)
	for _, o := range objects {
		size := o.MinSize()
		o.Resize(size)
		o.Move(pos)

		pos = pos.Add(fyne.NewPos(size.Width, size.Height))
	}
}

type enterEntry struct {
    widget.Entry
}

func (e *enterEntry) onEnter() {
    fmt.Println(e.Entry.Text)
    e.Entry.SetText("")
}

func newEnterEntry() *enterEntry {
    entry := &enterEntry{}
    entry.ExtendBaseWidget(entry)
    return entry
}

func (e *enterEntry) KeyDown(key *fyne.KeyEvent) {
    switch key.Name {
    case fyne.KeyReturn:
        e.onEnter()
    default:
        e.Entry.KeyDown(key)
    }

}

func main() {
    // a := app.New()
    // win := a.NewWindow("Hello World")
    // win.SetContent(widget.NewVBox(
    //     widget.NewLabel("Hello World!"),
    //     widget.NewButton("Quit", func() {
    //         a.Quit()
    //     }),
    // ))
	// win.ShowAndRun()
	

	// a := app.New()
	// w := a.NewWindow("Diagonal")

	// text1 := widget.NewLabel("topleft")
	// text2 := widget.NewLabel("Middle Label")
	// text3 := widget.NewLabel("bottomright")

	// w.SetContent(fyne.NewContainerWithLayout(&diagonal{}, text1, text2, text3))
	// w.ShowAndRun()

	a := app.New()
    w := a.NewWindow("Messenger")

    entry := newEnterEntry()

    w.SetContent(entry)
    w.ShowAndRun()
}