package client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type LockScreen struct {
	Window       fyne.Window
	OnUnlock     func()
	pcID         string
	adminLogin   *widget.Entry
	adminPass    *widget.Entry
	unlockButton *widget.Button
}

func NewLockScreen(pcID string) *LockScreen {
	a := app.New()
	w := a.NewWindow("PC Locked - " + pcID)
	w.SetFullScreen(true)

	return &LockScreen{
		Window: w,
		pcID:   pcID,
	}
}

func (ls *LockScreen) Show() {
	ls.adminLogin = widget.NewEntry()
	ls.adminPass = widget.NewPasswordEntry()
	ls.unlockButton = widget.NewButton("Unlock", func() {
		if ls.OnUnlock != nil {
			ls.OnUnlock()
		}
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Admin Login", Widget: ls.adminLogin},
			{Text: "Password", Widget: ls.adminPass},
		},
		OnSubmit: func() {
			ls.unlockButton.OnTapped()
		},
	}

	ls.Window.SetContent(container.NewCenter(
		container.NewVBox(
			widget.NewLabel("PC is locked"),
			widget.NewLabel("ID: "+ls.pcID),
			form,
			ls.unlockButton,
		),
	))
}

func (ls *LockScreen) SetUnlockCallback(callback func()) {
	ls.OnUnlock = callback
}

func (ls *LockScreen) Run() {
	ls.Show()
	ls.Window.ShowAndRun()
}
