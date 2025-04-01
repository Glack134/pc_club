package pc

import (
	"errors"
	"os/exec"
	"runtime"
)

var ErrUnsupportedOS = errors.New("unsupported operating system")

func LockPC() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("loginctl", "lock-session")
	case "windows":
		cmd = exec.Command("rundll32.exe", "user32.dll,LockWorkStation")
	default:
		return ErrUnsupportedOS
	}
	return cmd.Run()
}

func UnlockPC() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", "steam://")
	case "windows":
		cmd = exec.Command("start", "steam://")
	default:
		return ErrUnsupportedOS
	}
	return cmd.Run()
}

func DisableInput() {
	// Реализация для конкретной ОС
}

func EnableInput() {
	// Реализация для конкретной ОС
}
