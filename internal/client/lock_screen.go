package client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type LockScreen struct {
	Window    fyne.Window // Сделали поле экспортируемым (с большой буквы)
	onUnlock  func()
	adminMode bool
}

func NewLockScreen() *LockScreen {
	a := app.New()
	w := a.NewWindow("PC Club - Locked")
	w.SetFullScreen(true)

	return &LockScreen{
		Window: w,
	}
}

func (ls *LockScreen) Show() {
	loginBtn := widget.NewButton("Admin Login", func() {
		ls.showAdminAuth()
	})

	message := widget.NewLabel("PC is locked. Please wait for admin authorization")
	message.Alignment = fyne.TextAlignCenter

	content := container.NewCenter(
		container.NewVBox(
			message,
			loginBtn,
		),
	)

	ls.Window.SetContent(content)
	ls.Window.Show()
}

func (ls *LockScreen) showAdminAuth() {
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	form := widget.NewForm(
		widget.NewFormItem("Username", username),
		widget.NewFormItem("Password", password),
	)

	form.OnSubmit = func() {
		if ls.authenticate(username.Text, password.Text) {
			ls.onUnlock()
		}
	}

	ls.Window.SetContent(container.NewCenter(form))
}

func (ls *LockScreen) authenticate(user, pass string) bool {
	// Реализация gRPC вызова к серверу
	return true
}

func (ls *LockScreen) SetUnlockCallback(callback func()) {
	ls.onUnlock = callback
}
