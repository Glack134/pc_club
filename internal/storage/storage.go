package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Glack134/pc_club/internal/domain"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage interface {
	// PC методы
	CreatePC(ctx context.Context, pc *domain.PC) error
	GetPC(ctx context.Context, id string) (*domain.PC, error)
	UpdatePC(ctx context.Context, pc *domain.PC) error
	ListPCs(ctx context.Context) ([]*domain.PC, error)

	// User методы
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetUserByCredentials(ctx context.Context, name, passwordHash string) (*domain.User, error)

	// Session методы
	CreateSession(ctx context.Context, session *domain.Session) error
	GetActiveSession(ctx context.Context, pcID string) (*domain.Session, error)
	EndSession(ctx context.Context, sessionID string, endTime time.Time, cost float64) error

	// Transaction методы
	BeginTx(ctx context.Context) (Tx, error)
	Close() error
}

type Tx interface {
	Commit() error
	Rollback() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
