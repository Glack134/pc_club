package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Glack134/pc_club/internal/auth"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Session struct {
	ID        string
	UserID    string
	PCID      string
	ExpiresAt time.Time
}

func Init(dbPath string) error {
	// Создаем директорию если не существует
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return fmt.Errorf("failed to create db directory: %v", err)
	}

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Создаем таблицы
	if _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id TEXT PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            is_admin BOOLEAN NOT NULL DEFAULT FALSE
        );
        
        CREATE TABLE IF NOT EXISTS sessions (
            id TEXT PRIMARY KEY,
            user_id TEXT NOT NULL,
            pc_id TEXT NOT NULL,
            expires_at DATETIME NOT NULL,
            FOREIGN KEY(user_id) REFERENCES users(id)
        );
    `); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	// Создаем тестового админа если нет пользователей
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return fmt.Errorf("failed to count users: %v", err)
	}

	if count == 0 {
		hashedPassword, _ := auth.HashPassword("admin123")
		if _, err := db.Exec(
			"INSERT INTO users (id, username, password, is_admin) VALUES (?, ?, ?, ?)",
			"admin-1", "admin", hashedPassword, true,
		); err != nil {
			return fmt.Errorf("failed to create admin user: %v", err)
		}
		log.Println("Created default admin user: admin / admin123")
	}

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func CreateSession(userID, pcID string, duration time.Duration) (string, error) {
	id := generateSessionID()
	_, err := db.Exec(
		"INSERT INTO sessions VALUES (?, ?, ?, ?)",
		id,
		userID,
		pcID,
		time.Now().Add(duration),
	)
	return id, err
}

type User struct {
	ID       string
	Username string
	Password string
	IsAdmin  bool
}

func generateSessionID() string {
	return fmt.Sprintf("session-%d", time.Now().UnixNano())
}
