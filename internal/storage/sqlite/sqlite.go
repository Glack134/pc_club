package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Glack134/pc_club/internal/domain"
	"github.com/Glack134/pc_club/internal/storage"
	"go.uber.org/zap"
)

type SQLiteStorage struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewStorage создает новое подключение к SQLite
func NewStorage(dsn string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &SQLiteStorage{
		db:     db,
		logger: zap.NewNop(), // По умолчанию без логгирования
	}, nil
}

// SetLogger устанавливает логгер для хранилища
func (s *SQLiteStorage) SetLogger(logger *zap.Logger) {
	s.logger = logger
}

// Init выполняет миграции базы данных
func (s *SQLiteStorage) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			balance REAL DEFAULT 0,
			is_admin INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS pcs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			is_locked INTEGER DEFAULT 0,
			last_heartbeat TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			pc_id TEXT NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP,
			cost REAL DEFAULT 0,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(pc_id) REFERENCES pcs(id)
		);

		CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
		CREATE INDEX IF NOT EXISTS idx_sessions_pc_id ON sessions(pc_id);
	`)

	return err
}

// Close закрывает соединение с базой данных
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// PC методы

func (s *SQLiteStorage) CreatePC(ctx context.Context, pc *domain.PC) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO pcs (id, name, is_locked) VALUES (?, ?, ?)",
		pc.ID, pc.Name, pc.IsLocked)
	return err
}

func (s *SQLiteStorage) GetPC(ctx context.Context, id string) (*domain.PC, error) {
	var pc domain.PC
	var lastHeartbeat sql.NullTime

	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, is_locked, last_heartbeat FROM pcs WHERE id = ?", id).
		Scan(&pc.ID, &pc.Name, &pc.IsLocked, &lastHeartbeat)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	if lastHeartbeat.Valid {
		pc.LastHeartbeat = lastHeartbeat.Time
	}

	return &pc, nil
}

func (s *SQLiteStorage) UpdatePC(ctx context.Context, pc *domain.PC) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE pcs SET name = ?, is_locked = ?, last_heartbeat = ? WHERE id = ?",
		pc.Name, pc.IsLocked, pc.LastHeartbeat, pc.ID)
	return err
}

func (s *SQLiteStorage) ListPCs(ctx context.Context) ([]*domain.PC, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, is_locked, last_heartbeat FROM pcs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pcs []*domain.PC
	for rows.Next() {
		var pc domain.PC
		var lastHeartbeat sql.NullTime

		if err := rows.Scan(&pc.ID, &pc.Name, &pc.IsLocked, &lastHeartbeat); err != nil {
			return nil, err
		}

		if lastHeartbeat.Valid {
			pc.LastHeartbeat = lastHeartbeat.Time
		}

		pcs = append(pcs, &pc)
	}

	return pcs, nil
}

// User методы

func (s *SQLiteStorage) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO users (id, name, password_hash, balance, is_admin) VALUES (?, ?, ?, ?, ?)",
		user.ID, user.Name, user.PasswordHash, user.Balance, user.IsAdmin)
	return err
}

func (s *SQLiteStorage) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, password_hash, balance, is_admin, created_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.PasswordHash, &user.Balance, &user.IsAdmin, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *SQLiteStorage) GetUserByCredentials(ctx context.Context, name, passwordHash string) (*domain.User, error) {
	var user domain.User
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, password_hash, balance, is_admin, created_at FROM users WHERE name = ? AND password_hash = ?",
		name, passwordHash).
		Scan(&user.ID, &user.Name, &user.PasswordHash, &user.Balance, &user.IsAdmin, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Session методы

func (s *SQLiteStorage) CreateSession(ctx context.Context, session *domain.Session) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO sessions (id, user_id, pc_id, start_time, end_time, cost) VALUES (?, ?, ?, ?, ?, ?)",
		session.ID, session.UserID, session.PCID, session.StartTime, session.EndTime, session.Cost)
	return err
}

func (s *SQLiteStorage) GetActiveSession(ctx context.Context, pcID string) (*domain.Session, error) {
	var session domain.Session
	err := s.db.QueryRowContext(ctx,
		"SELECT id, user_id, pc_id, start_time, end_time, cost FROM sessions WHERE pc_id = ? AND end_time IS NULL",
		pcID).
		Scan(&session.ID, &session.UserID, &session.PCID, &session.StartTime, &session.EndTime, &session.Cost)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (s *SQLiteStorage) EndSession(ctx context.Context, sessionID string, endTime time.Time, cost float64) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE sessions SET end_time = ?, cost = ? WHERE id = ?",
		endTime, cost, sessionID)
	return err
}

// Transaction методы

func (s *SQLiteStorage) BeginTx(ctx context.Context) (storage.Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &SQLiteTx{tx: tx}, nil
}

type SQLiteTx struct {
	tx *sql.Tx
}

func (t *SQLiteTx) Commit() error {
	return t.tx.Commit()
}

func (t *SQLiteTx) Rollback() error {
	return t.tx.Rollback()
}

func (t *SQLiteTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *SQLiteTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *SQLiteTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}
