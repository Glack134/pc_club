package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("PC Club Admin")

	// Форма выдачи доступа
	userEntry := widget.NewEntry()
	pcEntry := widget.NewEntry()
	timeEntry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "User ID", Widget: userEntry},
			{Text: "PC ID", Widget: pcEntry},
			{Text: "Minutes", Widget: timeEntry},
		},
		OnSubmit: func() {
			// Здесь будет обработка выдачи доступа
		},
	}

	// Список активных сессий
	sessionsList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i int, o fyne.CanvasObject) {},
	)

	// Создаем вкладки
	tabs := container.NewAppTabs(
		container.NewTabItem("Grant Access", form),
		container.NewTabItem("Sessions", sessionsList),
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
