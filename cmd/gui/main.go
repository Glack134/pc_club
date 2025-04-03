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
	client    rpc.AdminServiceClient
	conn      *grpc.ClientConn
	authToken string
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
		client: rpc.NewAdminServiceClient(conn),
		conn:   conn,
	}

	// Форма авторизации
	adminApp.showLoginWindow(a, w)
}

func (app *AdminApp) showLoginWindow(a fyne.App, mainWindow fyne.Window) {
	loginWindow := a.NewWindow("Admin Login")
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
			resp, err := app.client.Login(context.Background(), &rpc.LoginRequest{
				Username: username.Text,
				Password: password.Text,
			})
			if err != nil {
				status.SetText("Login failed")
				return
			}

			if resp.Success {
				app.authToken = resp.Token
				loginWindow.Hide()
				app.showMainWindow(mainWindow)
			} else {
				status.SetText("Invalid credentials")
			}
		},
	}

	loginWindow.SetContent(container.NewVBox(
		widget.NewLabel("Admin Login"),
		form,
		status,
	))
	loginWindow.Show()
}

func (app *AdminApp) showMainWindow(w fyne.Window) {
	// Форма выдачи доступа
	userEntry := widget.NewEntry()
	pcEntry := widget.NewEntry()
	timeEntry := widget.NewEntry()
	resultLabel := widget.NewLabel("")

	// Список активных сессий
	sessionList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i int, o fyne.CanvasObject) {},
	)

	// Обновление списка сессий
	updateSessions := func() {
		md := metadata.Pairs("authorization", "Bearer "+app.authToken)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		resp, err := app.client.GetActiveSessions(ctx, &rpc.Empty{})
		if err != nil {
			resultLabel.SetText("Failed to get sessions")
			return
		}

		sessionList.Length = func() int { return len(resp.Sessions) }
		sessionList.UpdateItem = func(i int, o fyne.CanvasObject) {
			s := resp.Sessions[i]
			o.(*widget.Label).SetText(
				fmt.Sprintf("PC: %s, User: %s, Expires: %s",
					s.PcId, s.UserId, time.Unix(s.ExpiresAt, 0).Format("15:04:05")))
		}
		sessionList.Refresh()
	}

	// Форма выдачи доступа
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "User ID", Widget: userEntry},
			{Text: "PC ID", Widget: pcEntry},
			{Text: "Minutes", Widget: timeEntry},
		},
		OnSubmit: func() {
			minutes, err := strconv.Atoi(timeEntry.Text)
			if err != nil {
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
				resultLabel.SetText("Error: " + err.Error())
				return
			}

			resultLabel.SetText(resp.Message)
			updateSessions()
		},
	}

	// Кнопка принудительной блокировки
	forceLockBtn := widget.NewButton("Force Lock", func() {
		// Реализация блокировки
	})

	// Вкладки
	tabs := container.NewAppTabs(
		container.NewTabItem("Grant Access", container.NewVBox(
			form,
			resultLabel,
		)),
		container.NewTabItem("Sessions", container.NewBorder(
			nil,
			forceLockBtn,
			nil, nil,
			sessionList,
		)),
	)

	// Обновление данных каждые 5 секунд
	go func() {
		for range time.Tick(5 * time.Second) {
			updateSessions()
		}
	}()

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(600, 400))
	w.Show()
}
