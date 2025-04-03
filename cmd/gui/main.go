package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Glack134/pc_club/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type AdminApp struct {
	client     rpc.AdminServiceClient
	conn       *grpc.ClientConn
	authToken  string
	app        fyne.App
	mainWindow fyne.Window
}

func main() {
	a := app.New()
	w := a.NewWindow("PC Club Admin")

	// Конфигурация подключения
	conn, err := grpc.Dial("192.168.1.14:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	adminApp := &AdminApp{
		client:     rpc.NewAdminServiceClient(conn),
		conn:       conn,
		app:        a,
		mainWindow: w,
	}

	adminApp.showLoginWindow()
	a.Run()
}

func (app *AdminApp) showLoginWindow() {
	loginWindow := app.app.NewWindow("Admin Login")
	loginWindow.Resize(fyne.NewSize(300, 200))

	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	status := widget.NewLabel("")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: username},
			{Text: "Password", Widget: password},
		},
		OnSubmit: func() {
			go func() {
				resp, err := app.client.Login(context.Background(), &rpc.LoginRequest{
					Username: username.Text,
					Password: password.Text,
				})

				// Используем fyne.CurrentApp() вместо app.app.Driver()
				fyne.CurrentApp().SendNotification(fyne.NewNotification("Login", "Processing..."))

				if err != nil {
					fyne.CurrentApp().SendNotification(fyne.NewNotification("Error", "Login failed"))
					status.SetText("Login failed: " + err.Error())
					return
				}

				if resp.Success {
					app.authToken = resp.Token
					loginWindow.Hide()
					app.showMainWindow()
				} else {
					fyne.CurrentApp().SendNotification(fyne.NewNotification("Error", "Invalid credentials"))
					status.SetText("Invalid credentials")
				}
			}()
		},
	}

	loginWindow.SetContent(container.NewVBox(
		widget.NewLabel("Admin Login"),
		form,
		status,
	))
	loginWindow.Show()
}

func (app *AdminApp) showMainWindow() {
	// Форма выдачи доступа
	userEntry := widget.NewEntry()
	pcEntry := widget.NewEntry()
	timeEntry := widget.NewEntry()
	resultLabel := widget.NewLabel("")

	// Список активных сессий
	sessionList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i int, o fyne.CanvasObject) {},
	)

	// Функция безопасного обновления списка сессий
	updateSessions := func() {
		go func() {
			md := metadata.Pairs("authorization", "Bearer "+app.authToken)
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			resp, err := app.client.GetActiveSessions(ctx, &rpc.Empty{})
			if err != nil {
				fyne.CurrentApp().SendNotification(fyne.NewNotification("Error", "Failed to get sessions"))
				resultLabel.SetText("Failed to get sessions: " + err.Error())
				return
			}

			// Обновляем UI через главное окно
			app.mainWindow.Canvas().SetContent(
				container.NewAppTabs(
					container.NewTabItem("Grant Access", container.NewVBox(
						widget.NewForm(
							widget.NewFormItem("User ID", userEntry),
							widget.NewFormItem("PC ID", pcEntry),
							widget.NewFormItem("Minutes", timeEntry),
						),
						resultLabel,
					)),
					container.NewTabItem("Sessions", container.NewVScroll(sessionList)),
				),
			)

			sessionList.Length = func() int { return len(resp.Sessions) }
			sessionList.UpdateItem = func(i int, o fyne.CanvasObject) {
				s := resp.Sessions[i]
				o.(*widget.Label).SetText(
					fmt.Sprintf("PC: %s, User: %s, Expires: %s",
						s.PcId, s.UserId, time.Unix(s.ExpiresAt, 0).Format("15:04:05")))
			}
			sessionList.Refresh()
		}()
	}

	// Форма выдачи доступа
	grantAccessForm := widget.NewForm(
		widget.NewFormItem("User ID", userEntry),
		widget.NewFormItem("PC ID", pcEntry),
		widget.NewFormItem("Minutes", timeEntry),
	)

	grantAccessForm.OnSubmit = func() {
		go func() {
			minutes, err := strconv.Atoi(timeEntry.Text)
			if err != nil {
				fyne.CurrentApp().SendNotification(fyne.NewNotification("Error", "Invalid minutes"))
				resultLabel.SetText("Invalid minutes")
				return
			}

			md := metadata.Pairs("authorization", "Bearer "+app.authToken)
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			resp, err := app.client.GrantAccess(ctx, &rpc.GrantRequest{
				UserId:  userEntry.Text,
				PcId:    pcEntry.Text,
				Minutes: int32(minutes),
			})

			if err != nil {
				fyne.CurrentApp().SendNotification(fyne.NewNotification("Error", "Grant access failed"))
				resultLabel.SetText("Error: " + err.Error())
				return
			}

			fyne.CurrentApp().SendNotification(fyne.NewNotification("Success", resp.Message))
			resultLabel.SetText(resp.Message)
			updateSessions()
		}()
	}

	// Кнопка принудительной блокировки
	forceLockBtn := widget.NewButton("Force Lock", func() {
		// Реализация аналогична GrantAccess
	})

	// Вкладки
	tabs := container.NewAppTabs(
		container.NewTabItem("Grant Access", container.NewVBox(
			grantAccessForm,
			resultLabel,
		)),
		container.NewTabItem("Sessions", container.NewBorder(
			nil,
			forceLockBtn,
			nil, nil,
			container.NewVScroll(sessionList),
		)),
	)

	// Обновление данных каждые 5 секунд
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			updateSessions()
		}
	}()

	app.mainWindow.SetContent(tabs)
	app.mainWindow.Resize(fyne.NewSize(600, 400))
	app.mainWindow.Show()
	updateSessions()
}
