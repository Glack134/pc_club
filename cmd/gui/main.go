package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Game Access Control")

	userEntry := widget.NewEntry()
	pcEntry := widget.NewEntry()
	timeEntry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "User", Widget: userEntry},
			{Text: "PC", Widget: pcEntry},
			{Text: "Minutes", Widget: timeEntry},
		},
		OnSubmit: func() {
			// Обработка отправки формы
		},
	}

	w.SetContent(container.NewVBox(
		widget.NewLabel("Grant Access"),
		form,
	))
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}
