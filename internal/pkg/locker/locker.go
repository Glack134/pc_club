package locker

type Locker interface {
	Lock() error
	Unlock() error
	IsLocked() (bool, error)
}

// WindowsLocker - реализация для Windows
type WindowsLocker struct{}

func NewWindowsLocker() *WindowsLocker {
	return &WindowsLocker{}
}

func (w *WindowsLocker) Lock() error {
	// Реализация блокировки
	return nil
}

func (w *WindowsLocker) Unlock() error {
	// Реализация разблокировки
	return nil
}

func (w *WindowsLocker) IsLocked() (bool, error) {
	// Проверка состояния блокировки
	return false, nil
}

type hardwareMonitor interface {
	GetCPUUsage() (float64, error)
	GetGPUUsage() (float64, error)
	GetRAMUsage() (float64, error)
}
