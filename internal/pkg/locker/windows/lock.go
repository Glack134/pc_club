//go:build windows
// +build windows

package windows

import (
	"golang.org/x/sys/windows"
)

var (
	user32           = windows.NewLazySystemDLL("user32.dll")
	lockWorkStation  = user32.NewProc("LockWorkStation")
	sendMessage      = user32.NewProc("SendMessageW")
	getSystemMetrics = user32.NewProc("GetSystemMetrics")
)

type WindowsLocker struct{}

func New() *WindowsLocker {
	return &WindowsLocker{}
}

func (wl *WindowsLocker) Lock() error {
	_, _, err := lockWorkStation.Call()
	if err != windows.Errno(0) {
		return err
	}
	return nil
}

func (wl *WindowsLocker) Unlock() error {
	// В Windows разблокировка обычно требует ввода пароля
	return nil
}

func (wl *WindowsLocker) IsLocked() (bool, error) {
	const SM_REMOTESESSION = 0x1000
	ret, _, _ := getSystemMetrics.Call(SM_REMOTESESSION)
	return ret != 0, nil
}
