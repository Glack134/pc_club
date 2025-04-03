package client

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

type LockScreen struct {
	window         fyne.Window
	pcID           string
	unlockCallback func()
}

func NewLockScreen(pcID string) *LockScreen {
	a := app.New()
	w := a.NewWindow("PC Locked - " + pcID)
	w.SetFullScreen(true)

	return &LockScreen{
		window: w,
		pcID:   pcID,
	}
}

func (l *LockScreen) SetUnlockCallback(callback func()) {
	l.unlockCallback = callback
}

func (l *LockScreen) Show() {
	unlockBtn := widget.NewButton("Unlock", func() {
		if l.unlockCallback != nil {
			l.unlockCallback()
		}
		l.window.Close()
	})

	l.window.SetContent(container.NewCenter(
		container.NewVBox(
			widget.NewLabel("PC is locked"),
			widget.NewLabel("PC ID: "+l.pcID),
			unlockBtn,
		),
	))
	l.window.Show()
}
