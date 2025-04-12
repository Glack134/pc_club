package domain

import "time"

type PC struct {
	ID              string
	Name            string
	IsLocked        bool
	CpuUsage        float64
	RamUsage        float64
	RunningPrograms []string
	LastHeartbeat   time.Time
}

type User struct {
	ID           string
	Name         string
	PasswordHash string
	Balance      float64
	IsAdmin      bool
	CreatedAt    time.Time
}

type Session struct {
	ID        string
	UserID    string
	PCID      string
	StartTime time.Time
	EndTime   time.Time
	Cost      float64
}
