package modules

// SystemModule определяет интерфейс для всех модулей системного обслуживания
type SystemModule interface {
	// Execute выполняет задачу модуля
	// progressCallback вызывается для обновления прогресса (0.0 - 1.0)
	Execute(progressCallback func(progress float64, message string)) error
	
	// GetName возвращает название модуля
	GetName() string
	
	// GetDescription возвращает описание модуля
	GetDescription() string
	
	// RequiresRoot возвращает true если модуль требует root права
	RequiresRoot() bool
}

// ProgressCallback тип для функции обратного вызова прогресса
type ProgressCallback func(progress float64, message string)