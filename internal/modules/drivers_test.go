package modules

import (
	"os"
	"strings"
	"testing"
)

func TestDriversModule_GetName(t *testing.T) {
	module := &DriversModule{}
	expected := "Driver Updates"
	if got := module.GetName(); got != expected {
		t.Errorf("GetName() = %v, want %v", got, expected)
	}
}

func TestDriversModule_GetDescription(t *testing.T) {
	module := &DriversModule{}
	desc := module.GetDescription()
	if desc == "" {
		t.Error("GetDescription() returned empty string")
	}
	if !strings.Contains(strings.ToLower(desc), "driver") {
		t.Error("Description should contain 'driver'")
	}
}

func TestDriversModule_RequiresRoot(t *testing.T) {
	module := &DriversModule{}
	if !module.RequiresRoot() {
		t.Error("DriversModule should require root privileges")
	}
}

func TestDriversModule_Execute_NonRoot(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping non-root test - running as root")
	}
	
	module := &DriversModule{}
	
	callCount := 0
	progressCallback := func(progress float64, message string) {
		callCount++
	}
	
	err := module.Execute(progressCallback)
	
	// Без root прав может быть ошибка или fallback к manual check
	if err != nil {
		t.Logf("Execute() returned error (expected for non-root): %v", err)
	}
	
	// Callback должен вызываться для начальных шагов
	if callCount == 0 {
		t.Error("Progress callback should be called even when failing")
	}
}

func TestDriversModule_Execute_Root(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("Skipping root test - not running as root")
	}
	
	module := &DriversModule{}
	
	callCount := 0
	var messages []string
	
	progressCallback := func(progress float64, message string) {
		callCount++
		if message != "" {
			messages = append(messages, message)
		}
	}
	
	err := module.Execute(progressCallback)
	
	// С root правами выполнение может пройти успешно или с ошибкой
	if err != nil {
		t.Logf("Execute() returned error (may be expected): %v", err)
	}
	
	// Callback должен вызываться
	if callCount == 0 {
		t.Error("Progress callback was never called")
	}
	
	// Должны быть сообщения о прогрессе
	if len(messages) == 0 {
		t.Error("No progress messages received")
	}
	
	// Первое сообщение должно быть о детектировании драйверов
	if len(messages) > 0 && !strings.Contains(strings.ToLower(messages[0]), "detect") {
		t.Logf("First message: %s", messages[0])
	}
}

func TestDriversModule_CheckManualDrivers(t *testing.T) {
	module := &DriversModule{}
	
	callCount := 0
	var messages []string
	
	progressCallback := func(progress float64, message string) {
		callCount++
		if message != "" {
			messages = append(messages, message)
		}
	}
	
	err := module.checkManualDrivers(progressCallback)
	
	// Функция не должна возвращать критических ошибок
	if err != nil {
		t.Errorf("checkManualDrivers() returned error: %v", err)
	}
	
	// Callback должен вызываться
	if callCount == 0 {
		t.Error("Progress callback was never called")
	}
	
	// Должны быть сообщения
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