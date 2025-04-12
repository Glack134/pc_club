package monitoring

type HardwareMonitor interface {
	GetCPUUsage() (float64, error)
	GetRAMUsage() (float64, error)
}

type ProcessTracker interface {
	// Методы для работы с процессами
}

// Реализации интерфейсов
type hardwareMonitor struct{}

func NewHardwareMonitor() HardwareMonitor {
	return &hardwareMonitor{}
}

func (h *hardwareMonitor) GetCPUUsage() (float64, error) {
	// Реализация получения загрузки CPU
	return 0.0, nil
}

func (h *hardwareMonitor) GetRAMUsage() (float64, error) {
	// Реализация получения использования RAM
	return 0.0, nil
}

type processTracker struct{}

func NewProcessTracker() ProcessTracker {
	return &processTracker{}
}
