package modules

import (
	"testing"
)

// MockModule для тестирования интерфейса
type MockModule struct {
	name        string
	description string
	requiresRoot bool
	executeFunc func(func(float64, string)) error
}

func (m *MockModule) GetName() string {
	return m.name
}

func (m *MockModule) GetDescription() string {
	return m.description
}

func (m *MockModule) RequiresRoot() bool {
	return m.requiresRoot
}

func (m *MockModule) Execute(callback func(float64, string)) error {
	if m.executeFunc != nil {
		return m.executeFunc(callback)
	}
	// Симулируем выполнение
	callback(0.0, "Starting...")
	callback(0.5, "In progress...")
	callback(1.0, "Completed")
	return nil
}

func TestSystemModuleInterface(t *testing.T) {
	// Тестируем, что все наши модули реализуют интерфейс SystemModule
	modules := []SystemModule{
		&UpdatesModule{},
		&DriversModule{},
		&CleanupModule{},
		&OptimizationModule{},
		&HealthModule{},
	}
	
	for _, module := range modules {
		// Проверяем, что все методы интерфейса работают
		name := module.GetName()
		if name == "" {
			t.Errorf("Module %T returned empty name", module)
		}
		
		desc := module.GetDescription()
		if desc == "" {
			t.Errorf("Module %T returned empty description", module)
		}
		
		// RequiresRoot должен возвращать bool (проверяется компилятором)
		_ = module.RequiresRoot()
		
		// Execute должен принимать callback функцию
		callbackCalled := false
		err := module.Execute(func(progress float64, message string) {
			callbackCalled = true
			if progress < 0 || progress > 1 {
				t.Errorf("Module %T returned invalid progress: %f", module, progress)
			}
		})
		
		// Некоторые модули могут возвращать ошибки в тестовой среде
		if err != nil {
			t.Logf("Module %T returned error (may be expected): %v", module, err)
		}
		
		if !callbackCalled {
			t.Errorf("Module %T did not call progress callback", module)
		}
	}
}

func TestMockModule(t *testing.T) {
	mock := &MockModule{
		name:         "Test Module",
		description:  "Test Description",
		requiresRoot: true,
	}
	
	// Проверяем базовые методы
	if mock.GetName() != "Test Module" {
		t.Error("MockModule GetName() failed")
	}
	
	if mock.GetDescription() != "Test Description" {
		t.Error("MockModule GetDescription() failed")
	}
	
	if !mock.RequiresRoot() {
		t.Error("MockModule RequiresRoot() failed")
	}
	
	// Проверяем Execute
	callCount := 0
	var progressValues []float64
	var messages []string
	
	err := mock.Execute(func(progress float64, message string) {
		callCount++
		progressValues = append(progressValues, progress)
		messages = append(messages, message)
	})
	
	if err != nil {
		t.Errorf("MockModule Execute() returned error: %v", err)
	}
	
	if callCount != 3 {
		t.Errorf("Expected 3 callback calls, got %d", callCount)
	}
	
	expectedProgress := []float64{0.0, 0.5, 1.0}
	for i, expected := range expectedProgress {
		if i >= len(progressValues) || progressValues[i] != expected {
			t.Errorf("Expected progress[%d] = %f, got %f", i, expected, progressValues[i])
		}
	}
	
	expectedMessages := []string{"Starting...", "In progress...", "Completed"}
	for i, expected := range expectedMessages {
		if i >= len(messages) || messages[i] != expected {
			t.Errorf("Expected message[%d] = %s, got %s", i, expected, messages[i])
		}
	}
}

func TestMockModuleWithCustomExecute(t *testing.T) {
	mock := &MockModule{
		name:        "Custom Module",
		description: "Custom Description",
		executeFunc: func(callback func(float64, string)) error {
			callback(0.25, "Custom step 1")
			callback(0.75, "Custom step 2")
			callback(1.0, "Custom completed")
			return nil
		},
	}
	
	var messages []string
	err := mock.Execute(func(progress float64, message string) {
		messages = append(messages, message)
	})
	
	if err != nil {
		t.Errorf("Custom execute returned error: %v", err)
	}
	
	expected := []string{"Custom step 1", "Custom step 2", "Custom completed"}
	if len(messages) != len(expected) {
		t.Errorf("Expected %d messages, got %d", len(expected), len(messages))
	}
	
	for i, exp := range expected {
		if i >= len(messages) || messages[i] != exp {
			t.Errorf("Expected message[%d] = %s, got %s", i, exp, messages[i])
		}
	}
}