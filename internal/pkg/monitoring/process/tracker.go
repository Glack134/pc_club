package process

import (
	"errors"
	"os/exec"
	"strings"
)

type ProcessTracker struct{}

func NewProcessTracker() *ProcessTracker {
	return &ProcessTracker{}
}

func (pt *ProcessTracker) IsProcessRunning(name string) (bool, error) {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq "+name)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 3 {
		return false, nil
	}

	return strings.Contains(lines[2], name), nil
}

func (pt *ProcessTracker) KillProcess(name string) error {
	cmd := exec.Command("taskkill", "/IM", name, "/F")
	return cmd.Run()
}

func (pt *ProcessTracker) StartProcess(path string) error {
	if path == "" {
		return errors.New("empty path")
	}
	return exec.Command("cmd", "/C", "start", "", path).Start()
}
