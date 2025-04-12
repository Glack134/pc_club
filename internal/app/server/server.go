package server

import (
	"context"
	"sync"

	"github.com/Glack134/pc_club/internal/domain" // Добавлен импорт domain
	"github.com/Glack134/pc_club/internal/storage"
	"github.com/Glack134/pc_club/pkg/logger"
	"github.com/Glack134/pc_club/pkg/rpc/admin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PcClubServer struct {
	admin.UnimplementedAdminServiceServer
	logger  logger.Logger
	storage storage.Storage
	mu      sync.Mutex
	pcs     map[string]*domain.PC // Изменено с PcStatus на PC
}

func NewPcClubServer(logger logger.Logger, storage storage.Storage) *PcClubServer {
	return &PcClubServer{
		logger:  logger,
		storage: storage,
		pcs:     make(map[string]*domain.PC), // Изменено с PcStatus на PC
	}
}

func (s *PcClubServer) RegisterServices(grpcServer *grpc.Server) {
	admin.RegisterAdminServiceServer(grpcServer, s)
}

func (s *PcClubServer) GetPcStatus(ctx context.Context, req *admin.PcStatusRequest) (*admin.PcStatusResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pc, ok := s.pcs[req.PcId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "PC not found")
	}

	return &admin.PcStatusResponse{
		PcId:            req.PcId,
		IsLocked:        pc.IsLocked,
		CpuUsage:        pc.CpuUsage,
		RamUsage:        pc.RamUsage,
		RunningPrograms: pc.RunningPrograms,
	}, nil
}

func (s *PcClubServer) LockPc(ctx context.Context, req *admin.LockPcRequest) (*admin.LockPcResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pc, ok := s.pcs[req.PcId]
	if !ok {
		pc = &domain.PC{ID: req.PcId}
		s.pcs[req.PcId] = pc
	}

	pc.IsLocked = true
	s.logger.Info("PC locked")
	return &admin.LockPcResponse{Success: true}, nil
}
