package storage

import (
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	dbPath := "file:test.db?mode=memory&cache=shared"
	err := Init(dbPath)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Close()

	sessionID, err := CreateSession("user1", "pc1", 1*time.Minute)
	if err != nil {
		t.Errorf("CreateSession failed: %v", err)
	}
	if sessionID == "" {
		t.Error("Empty session ID")
	}
}
