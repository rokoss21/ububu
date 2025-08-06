package modules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanupModule_GetName(t *testing.T) {
	module := &CleanupModule{}
	expected := "System Cleanup"
	if got := module.GetName(); got != expected {
		t.Errorf("GetName() = %v, want %v", got, expected)
	}
}

func TestCleanupModule_GetDescription(t *testing.T) {
	module := &CleanupModule{}
	desc := module.GetDescription()
	if desc == "" {
		t.Error("GetDescription() returned empty string")
	}
	if !strings.Contains(strings.ToLower(desc), "clean") {
		t.Error("Description should contain 'clean'")
	}
}

func TestCleanupModule_RequiresRoot(t *testing.T) {
	module := &CleanupModule{}
	if module.RequiresRoot() {
		t.Error("CleanupModule should not require root privileges")
	}
}

func TestCleanupModule_Execute(t *testing.T) {
	module := &CleanupModule{}
	
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
}

func TestCleanupModule_GetDirSize(t *testing.T) {
	module := &CleanupModule{}
	
	// Создаем временную директорию для тестирования
	tempDir, err := os.MkdirTemp("", "ububu_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Создаем тестовый файл
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World! This is a test file."
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Тестируем getDirSize
	size, err := module.getDirSize(tempDir)
	if err != nil {
		t.Errorf("getDirSize() returned error: %v", err)
	}
	
	expectedSize := int64(len(testContent))
	if size != expectedSize {
		t.Errorf("getDirSize() = %d, want %d", size, expectedSize)
	}
}

func TestCleanupModule_GetDirSize_NonExistent(t *testing.T) {
	module := &CleanupModule{}
	
	// Тестируем с несуществующей директорией
	size, err := module.getDirSize("/non/existent/directory")
	
	// Функция может вернуть ошибку или 0 (в зависимости от реализации)
	if err != nil {
		t.Logf("getDirSize() returned expected error for non-existent directory: %v", err)
	}
	
	// Размер должен быть 0 для несуществующей директории
	if size != 0 {
		t.Errorf("getDirSize() should return 0 for non-existent directory, got %d", size)
	}
}

func TestCleanupModule_CleanPackageCache(t *testing.T) {
	module := &CleanupModule{}
	
	// Этот тест может потребовать root права, поэтому проверяем
	if os.Geteuid() != 0 {
		t.Skip("Skipping package cache test - requires root privileges")
	}
	
	size, err := module.cleanPackageCache()
	
	// Проверяем, что функция выполнилась без критических ошибок
	if err != nil {
		t.Logf("cleanPackageCache() returned error (may be expected): %v", err)
	}
	
	// Размер должен быть неотрицательным
	if size < 0 {
		t.Errorf("cleanPackageCache() returned negative size: %d", size)
	}
}

func TestCleanupModule_CleanBrowserCache(t *testing.T) {
	module := &CleanupModule{}
	
	// Создаем временную домашнюю директорию для тестирования
	tempHome, err := os.MkdirTemp("", "ububu_home_test")
	if err != nil {
		t.Fatalf("Failed to create temp home dir: %v", err)
	}
	defer os.RemoveAll(tempHome)
	
	// Создаем поддиректории кэша браузера
	cacheDir := filepath.Join(tempHome, ".cache", "google-chrome")
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create cache dir: %v", err)
	}
	
	// Создаем тестовый файл кэша
	cacheFile := filepath.Join(cacheDir, "cache.dat")
	err = os.WriteFile(cacheFile, []byte("test cache data"), 0644)
	if err != nil {
		t.Fatalf("Failed to create cache file: %v", err)
	}
	
	// Временно изменяем HOME для теста
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)
	
	size, err := module.cleanBrowserCache()
	
	// Проверяем результат
	if err != nil {
		t.Errorf("cleanBrowserCache() returned error: %v", err)
	}
	
	if size < 0 {
		t.Errorf("cleanBrowserCache() returned negative size: %d", size)
	}
	
	// Проверяем, что файл кэша был удален
	if _, err := os.Stat(cacheFile); !os.IsNotExist(err) {
		t.Error("Cache file should have been deleted")
	}
}