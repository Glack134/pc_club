package steam

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	steamRegKey = "SOFTWARE\\Valve\\Steam"
)

type SteamManager struct {
	installPath string
}

func NewSteamManager() (*SteamManager, error) {
	// Поиск пути установки Steam
	path, err := findSteamPath()
	if err != nil {
		return nil, err
	}
	return &SteamManager{installPath: path}, nil
}

func findSteamPath() (string, error) {
	// Проверка стандартных путей
	paths := []string{
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Steam"),
		filepath.Join(os.Getenv("ProgramFiles"), "Steam"),
		"C:\\Steam",
	}

	for _, path := range paths {
		if _, err := os.Stat(filepath.Join(path, "steam.exe")); err == nil {
			return path, nil
		}
	}

	return "", errors.New("Steam installation not found")
}

func (sm *SteamManager) LaunchGame(appID string) error {
	steamCmd := filepath.Join(sm.installPath, "steam.exe")
	cmd := exec.Command(steamCmd, "-applaunch", appID)
	return cmd.Start()
}

func (sm *SteamManager) CloseGame(appID string) error {
	// Реализация закрытия игры
	return nil
}
