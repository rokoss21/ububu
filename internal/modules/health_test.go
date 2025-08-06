package modules

import (
	"strings"
	"testing"
)

func TestHealthModule_GetName(t *testing.T) {
	module := &HealthModule{}
	expected := "System Health Check"
	if got := module.GetName(); got != expected {
		t.Errorf("GetName() = %v, want %v", got, expected)
	}
}

func TestHealthModule_GetDescription(t *testing.T) {
	module := &HealthModule{}
	desc := module.GetDescription()
	if desc == "" {
		t.Error("GetDescription() returned empty string")
	}
	if !strings.Contains(desc, "health") {
		t.Error("Description should contain 'health'")
	}
}

func TestHealthModule_RequiresRoot(t *testing.T) {
	module := &HealthModule{}
	if module.RequiresRoot() {
		t.Error("HealthModule should not require root privileges")
	}
}

func TestHealthModule_Execute(t *testing.T) {
	module := &HealthModule{}
	
	// Счетчик вызовов callback
	callCount := 0
	var lastProgress float64
	var messages []string
	
	progressCallback := func(progress float64, message string) {
		callCount++
		lastProgress = progress
		if message != "" {
			messages = append(messages, message)
		}
	}
	
	err := module.Execute(progressCallback)
	
	// Проверяем, что выполнение прошло без ошибок
	if err != nil {
		t.Errorf("Execute() returned error: %v", err)
	}
	
	// Проверяем, что callback вызывался
	if callCount == 0 {
		t.Error("Progress callback was never called")
	}
	
	// Проверяем, что финальный прогресс = 1.0
	if lastProgress != 1.0 {
		t.Errorf("Final progress = %v, want 1.0", lastProgress)
	}
	
	// Проверяем, что были сообщения
	if len(messages) == 0 {
		t.Error("No progress messages received")
	}
	
	// Проверяем, что последнее сообщение содержит "completed"
	lastMessage := messages[len(messages)-1]
	if !strings.Contains(strings.ToLower(lastMessage), "completed") {
		t.Errorf("Last message should contain 'completed', got: %s", lastMessage)
	}
}

func TestHealthModule_CheckDiskHealth(t *testing.T) {
	module := &HealthModule{}
	
	result, err := module.checkDiskHealth()
	
	// Проверяем, что функция не возвращает ошибку
	if err != nil {
		t.Errorf("checkDiskHealth() returned error: %v", err)
	}
	
	// Проверяем, что результат не пустой
	if result == "" {
		t.Error("checkDiskHealth() returned empty result")
	}
	
	// Проверяем, что результат содержит статус
	if !strings.Contains(result, "GOOD") && !strings.Contains(result, "WARNING") && !strings.Contains(result, "CAUTION") {
		t.Errorf("Result should contain status, got: %s", result)
	}
}

func TestHealthModule_CheckMemoryUsage(t *testing.T) {
	module := &HealthModule{}
	
	result, err := module.checkMemoryUsage()
	
	// Проверяем, что функция не возвращает ошибку
	if err != nil {
		t.Errorf("checkMemoryUsage() returned error: %v", err)
	}
	
	// Проверяем, что результат не пустой
	if result == "" {
		t.Error("checkMemoryUsage() returned empty result")
	}
	
	// Проверяем, что результат содержит статус памяти
	if !strings.Contains(result, "NORMAL") && !strings.Contains(result, "HIGH") && !strings.Contains(result, "CRITICAL") {
		t.Errorf("Result should contain memory status, got: %s", result)
	}
	
	// Проверяем, что результат содержит проценты
	if !strings.Contains(result, "%") {
		t.Errorf("Result should contain percentage, got: %s", result)
	}
}

func TestHealthModule_AnalyzeProcesses(t *testing.T) {
	module := &HealthModule{}
	
	result, err := module.analyzeProcesses()
	
	// Проверяем, что функция не возвращает ошибку
	if err != nil {
		t.Errorf("analyzeProcesses() returned error: %v", err)
	}
	
	// Проверяем, что результат не пустой
	if result == "" {
		t.Error("analyzeProcesses() returned empty result")
	}
	
	// Проверяем, что результат содержит информацию о нагрузке
	if !strings.Contains(result, "load") {
		t.Errorf("Result should contain load information, got: %s", result)
	}
	
	// Проверяем, что результат содержит количество процессов
	if !strings.Contains(result, "processes") {
		t.Errorf("Result should contain process count, got: %s", result)
	}
}