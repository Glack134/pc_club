package storage

import (
	"database/sql"

	"github.com/Glack134/pc_club/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func CreateUser(user *models.User) error {
	_, err := DB.Exec(
		"INSERT INTO users (id, username, password, is_admin) VALUES (?, ?, ?, ?)",
		user.ID, user.Username, user.Password, user.IsAdmin,
	)
	return err
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow(
		"SELECT id, username, password, is_admin FROM users WHERE username = ?",

		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin)
	return &user, err
}
