package locker

type Locker interface {
	Lock() error
	Unlock() error
	IsLocked() (bool, error)
}

type hardwareMonitor interface {
	GetCPUUsage() (float64, error)
	GetGPUUsage() (float64, error)
	GetRAMUsage() (float64, error)
}
