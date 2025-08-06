package modules

import (
	"os"
	"strings"
	"testing"
)

func TestUpdatesModule_GetName(t *testing.T) {
	module := &UpdatesModule{}
	expected := "System Updates"
	if got := module.GetName(); got != expected {
		t.Errorf("GetName() = %v, want %v", got, expected)
	}
}

func TestUpdatesModule_GetDescription(t *testing.T) {
	module := &UpdatesModule{}
	desc := module.GetDescription()
	if desc == "" {
		t.Error("GetDescription() returned empty string")
	}
	if !strings.Contains(strings.ToLower(desc), "update") {
		t.Error("Description should contain 'update'")
	}
}

func TestUpdatesModule_RequiresRoot(t *testing.T) {
	module := &UpdatesModule{}
	if !module.RequiresRoot() {
		t.Error("UpdatesModule should require root privileges")
	}
}

func TestUpdatesModule_Execute_NonRoot(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping non-root test - running as root")
	}
	
	module := &UpdatesModule{}
	
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

func TestUpdatesModule_Execute_Root(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("Skipping root test - not running as root")
	}
	
	module := &UpdatesModule{}
	
	callCount := 0
	var messages []string
	
	progressCallback := func(progress float64, message string) {
		callCount++
		if message != "" {
			messages = append(messages, message)
		}
	}
	
	err := module.Execute(progressCallback)
	
	// С root правами выполнение может пройти успешно или с ошибкой (зависит от системы)
	if err != nil {
		t.Logf("Execute() returned error (may be expected in test environment): %v", err)
	}
	
	// Callback должен вызываться
	if callCount == 0 {
		t.Error("Progress callback was never called")
	}
	
	// Должны быть сообщения о прогрессе
	if len(messages) == 0 {
		t.Error("No progress messages received")
	}
	
	// Первое сообщение должно быть о обновлении списков пакетов
	if len(messages) > 0 && !strings.Contains(strings.ToLower(messages[0]), "updating") {
		t.Errorf("First message should be about updating, got: %s", messages[0])
	}
}