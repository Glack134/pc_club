package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Glack134/pc_club/internal/auth"
	"github.com/Glack134/pc_club/internal/storage"
	"github.com/Glack134/pc_club/pkg/config"
	"github.com/Glack134/pc_club/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var clientConnections = make(map[string]rpc.AdminServiceClient)

type server struct {
	rpc.UnimplementedAdminServiceServer
	config *config.Config
}

func getClientForPC(pcID string) (rpc.AdminServiceClient, error) {
	// В реальной реализации здесь должно быть:
	// 1. Поиск адреса клиента по pcID (из конфига или БД)
	// 2. Установка gRPC соединения
	// 3. Сохранение соединения в clientConnections

	if client, ok := clientConnections[pcID]; ok {
		return client, nil
	}

	// Пример для тестирования - используйте реальные адреса клиентов
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:50051", pcID), // Замените на реальный адрес клиента
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := rpc.NewAdminServiceClient(conn)
	clientConnections[pcID] = client
	return client, nil
}

func (s *server) TerminateSession(ctx context.Context, req *rpc.SessionRequest) (*rpc.Response, error) {
	claims, _ := ctx.Value("claims").(*auth.Claims)
	if !claims.IsAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin access required")
	}

	session, err := storage.GetSession(req.SessionId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "session not found")
	}

	if err := storage.TerminateSession(req.SessionId); err != nil {
		return nil, status.Error(codes.Internal, "failed to terminate session")
	}

	go func() {
		client, err := getClientForPC(session.PcID)
		if err == nil {
			_, _ = client.LockPC(context.Background(), &rpc.PCRequest{PcId: session.PcID})
		}
	}()

	_ = storage.LogAction(claims.UserID, "session_terminated",
		fmt.Sprintf("Session %s terminated", req.SessionId))

	return &rpc.Response{Success: true, Message: "Session terminated"}, nil
}

func (s *server) LockPC(ctx context.Context, req *rpc.PCRequest) (*rpc.Response, error) {
	// Проверяем права администратора
	claims, ok := ctx.Value("claims").(*auth.Claims)
	if !ok || !claims.IsAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin access required")
	}

	// В реальной реализации здесь должна быть логика блокировки PC
	log.Printf("Received LockPC request for PC: %s", req.PcId)

	return &rpc.Response{
		Success: true,
		Message: fmt.Sprintf("PC %s locked successfully", req.PcId),
	}, nil
}

func (s *server) UnlockPC(ctx context.Context, req *rpc.PCRequest) (*rpc.Response, error) {
	// Проверяем права администратора
	claims, ok := ctx.Value("claims").(*auth.Claims)
	if !ok || !claims.IsAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin access required")
	}

	// В реальной реализации здесь должна быть логика разблокировки PC
	log.Printf("Received UnlockPC request for PC: %s", req.PcId)

	return &rpc.Response{
		Success: true,
		Message: fmt.Sprintf("PC %s unlocked successfully", req.PcId),
	}, nil
}

func (s *server) ForceLockPC(ctx context.Context, req *rpc.PCRequest) (*rpc.Response, error) {
	claims, _ := ctx.Value("claims").(*auth.Claims)
	if !claims.IsAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin access required")
	}

	client, err := getClientForPC(req.PcId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "PC client not available")
	}

	if _, err := client.LockPC(context.Background(), req); err != nil {
		return nil, status.Error(codes.Internal, "failed to lock PC")
	}

	_ = storage.LogAction(claims.UserID, "force_locked",
		fmt.Sprintf("PC %s force locked", req.PcId))

	return &rpc.Response{Success: true, Message: "PC locked"}, nil
}

func (s *server) GetActiveSessions(ctx context.Context, _ *rpc.Empty) (*rpc.SessionsResponse, error) {
	sessions, err := storage.GetActiveSessions()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get sessions")
	}

	var pbSessions []*rpc.Session
	for _, s := range sessions {
		pbSessions = append(pbSessions, &rpc.Session{
			Id:        s.ID,
			UserId:    s.UserID,
			PcId:      s.PcID,
			ExpiresAt: s.ExpiresAt.Unix(),
		})
	}

	return &rpc.SessionsResponse{Sessions: pbSessions}, nil
}

func (s *server) GrantAccess(ctx context.Context, req *rpc.GrantRequest) (*rpc.Response, error) {
	// Claims уже проверены в interceptor, просто получаем их
	claims, ok := ctx.Value("claims").(*auth.Claims)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to get auth claims")
	}

	if !claims.IsAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin access required")
	}

	// 2. Создание сессии
	sessionID, err := storage.CreateSession(req.UserId, req.PcId, time.Duration(req.Minutes)*time.Minute)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		return nil, status.Error(codes.Internal, "failed to create session")
	}

	log.Printf("Access granted. Session: %s, User: %s, PC: %s, Duration: %d mins",
		sessionID, req.UserId, req.PcId, req.Minutes)

	return &rpc.Response{
		Success: true,
		Message: "Access granted successfully",
	}, nil
}

func (s *server) Login(ctx context.Context, req *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	user, err := storage.GetUserByUsername(req.Username)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &rpc.LoginResponse{
		Token:   token,
		Success: true,
	}, nil
}

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация секретного ключа для JWT
	auth.SetSecretKey(cfg.JWTSecret)

	// Инициализация БД
	if err := storage.Init(cfg.DBPath); err != nil {
		log.Fatalf("Failed to init storage: %v", err)
	}
	defer storage.Close()

	// Настройка gRPC сервера
	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)
	rpc.RegisterAdminServiceServer(srv, &server{config: cfg})

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		srv.GracefulStop()
	}()

	log.Printf("Server started on port %s", cfg.ServerPort)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Пропускаем аутентификацию для метода Login
	if info.FullMethod == "/rpc.AdminService/Login" {
		return handler(ctx, req)
	}

	// Получаем токен из метаданных
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata not provided")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token not provided")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	// Добавляем claims в контекст для использования в обработчиках
	ctx = context.WithValue(ctx, "claims", claims)

	return handler(ctx, req)
}
