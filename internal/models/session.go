package models

import "time"

type Session struct {
	ID        string
	UserID    string
	PCID      string
	ExpiresAt time.Time
	CreatedAt time.Time
}
