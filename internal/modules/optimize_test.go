package modules

import (
	"os"
	"strings"
	"testing"
)

func TestOptimizationModule_GetName(t *testing.T) {
	module := &OptimizationModule{}
	expected := "System Optimization"
	if got := module.GetName(); got != expected {
		t.Errorf("GetName() = %v, want %v", got, expected)
	}
}

func TestOptimizationModule_GetDescription(t *testing.T) {
	module := &OptimizationModule{}
	desc := module.GetDescription()
	if desc == "" {
		t.Error("GetDescription() returned empty string")
	}
	if !strings.Contains(strings.ToLower(desc), "optim") {
		t.Error("Description should contain 'optim'")
	}
}

func TestOptimizationModule_RequiresRoot(t *testing.T) {
	module := &OptimizationModule{}
	if !module.RequiresRoot() {
		t.Error("OptimizationModule should require root privileges")
	}
}

func TestOptimizationModule_Execute_NonRoot(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping non-root test - running as root")
	}
	
	module := &OptimizationModule{}
	
	callCount := 0
	progressCallback := func(progress float64, message string) {
		callCount++
	}
	
	err := module.Execute(progressCallback)
	
	// Без root прав должна быть ошибка
	if err == nil {
		t.Error("Execute() should return error when not running as root")
	}
	
	// Callback все равно должен вызываться для начальных шагов
	if callCount == 0 {
		t.Error("Progress callback should be called even when failing")
	}
}

func TestOptimizationModule_Execute_Root(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("Skipping root test - not running as root")
	}
	
	module := &OptimizationModule{}
	
	callCount := 0
	var messages []string
	
	progressCallback := func(progress float64, message string) {
		callCount++
		if message != "" {
			messages = append(messages, message)
		}
	}
	
	err := module.Execute(progressCallback)
	
	// С root правами выполнение должно пройти успешно
	if err != nil {
		t.Errorf("Execute() returned error: %v", err)
	}
	
	// Callback должен вызываться
	if callCount == 0 {
		t.Error("Progress callback was never called")
	}
	
	// Должны быть сообщения о прогрессе
	if len(messages) == 0 {
		t.Error("No progress messages received")
	}
	
	// Проверяем, что есть сообщение о завершении
	if len(messages) > 0 {
		lastMessage := messages[len(messages)-1]
		if !strings.Contains(strings.ToLower(lastMessage), "completed") {
			t.Errorf("Last message should indicate completion, got: %s", lastMessage)
		}
	}
}

func TestOptimizationModule_OptimizeSSD(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("Skipping SSD optimization test - requires root")
	}
	
	module := &OptimizationModule{}
	
	callCount := 0
	var messages []string
	
	progressCallback := func(progress float64, message string) {
		callCount++
		if message != "" {
			messages = append(messages, message)
		}
	}
	
	err := module.optimizeSSD(progressCallback)
	
	// Функция не должна возвращать критических ошибок
	if err != nil {
		t.Errorf("optimizeSSD() returned error: %v", err)
	}
	
	// Callback должен вызываться
	if callCount == 0 {
		t.Error("Progress callback was never called")
	}
	
	// Должны быть сообщения
	if len(messages) == 0 {
		t.Error("No progress messages received")
	}
}

func TestOptimizationModule_OptimizeMemory(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("Skipping memory optimization test - requires root")
	}
	
	module := &OptimizationModule{}
	
	callCount := 0
	var messages []string
	
	progressCallback := func(progress float64, message string) {
		callCount++
		if message != "" {
			messages = append(messages, message)
		}
	}
	
	err := module.optimizeMemory(progressCallback)
	
	// Функция не должна возвращать критических ошибок
	if err != nil {
		t.Errorf("optimizeMemory() returned error: %v", err)
	}
	
	// Callback должен вызываться
	if callCount == 0 {
		t.Error("Progress callback was never called")
	}
	
	// Должны быть сообщения о swappiness
	found := false
	for _, msg := range messages {
		if strings.Contains(strings.ToLower(msg), "swappiness") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should have message about swappiness")
	}
}

func TestOptimizationModule_ClearNetworkCache(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("Skipping network cache test - requires root")
	}
	
	module := &OptimizationModule{}
	
	callCount := 0
	var messages []string
	
	progressCallback := func(progress float64, message string) {
		callCount++
		if message != "" {
			messages = append(messages, message)
		}
	}
	
	err := module.clearNetworkCache(progressCallback)
	
	// Функция может вернуть ошибку в зависимости от системы
	if err != nil {
		t.Logf("clearNetworkCache() returned error (may be expected): %v", err)
	}
	
	// Callback должен вызываться
	if callCount == 0 {
		t.Error("Progress callback was never called")
	}
	
	// Должны быть сообщения о DNS
	found := false
	for _, msg := range messages {
		if strings.Contains(strings.ToLower(msg), "dns") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should have message about DNS")
	}
}