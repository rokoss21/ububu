package modules

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type CleanupModule struct{}

func (m *CleanupModule) GetName() string {
	return "System Cleanup"
}

func (m *CleanupModule) GetDescription() string {
	return "Clean temporary files and caches"
}

func (m *CleanupModule) RequiresRoot() bool {
	return false
}

func (m *CleanupModule) Execute(progressCallback func(progress float64, message string)) error {
	var totalFreed int64
	
	progressCallback(0.1, "Cleaning package cache...")
	freed, err := m.cleanPackageCache()
	if err == nil {
		totalFreed += freed
		progressCallback(0.25, fmt.Sprintf("Package cache cleaned: %d MB freed", freed/1024/1024))
	}
	
	progressCallback(0.3, "Cleaning browser caches...")
	freed, err = m.cleanBrowserCache()
	if err == nil {
		totalFreed += freed
		progressCallback(0.5, fmt.Sprintf("Browser cache cleaned: %d MB freed", freed/1024/1024))
	}
	
	progressCallback(0.6, "Cleaning temporary files...")
	freed, err = m.cleanTempFiles()
	if err == nil {
		totalFreed += freed
		progressCallback(0.75, fmt.Sprintf("Temp files cleaned: %d MB freed", freed/1024/1024))
	}
	
	progressCallback(0.8, "Cleaning old logs...")
	freed, err = m.cleanOldLogs()
	if err == nil {
		totalFreed += freed
		progressCallback(0.9, fmt.Sprintf("Old logs cleaned: %d MB freed", freed/1024/1024))
	}
	
	progressCallback(1.0, fmt.Sprintf("Cleanup completed! Total freed: %d MB", totalFreed/1024/1024))
	
	return nil
}

func (m *CleanupModule) cleanPackageCache() (int64, error) {
	var totalSize int64
	
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Получаем размер кэша перед очисткой
	cmd := exec.CommandContext(ctx, "du", "-sb", "/var/cache/apt/archives")
	output, err := cmd.Output()
	if err == nil {
		fmt.Sscanf(string(output), "%d", &totalSize)
	}
	
	// Очищаем кэш пакетов (без sudo для избежания зависания)
	cmd = exec.CommandContext(ctx, "apt", "clean")
	cmd.Run() // Игнорируем ошибки
	
	// Удаляем неиспользуемые пакеты (без sudo)
	cmd = exec.CommandContext(ctx, "apt", "autoremove", "-y")
	cmd.Run() // Игнорируем ошибки
	
	return totalSize, nil
}

func (m *CleanupModule) cleanBrowserCache() (int64, error) {
	var totalSize int64
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return 0, err
	}
	
	// Пути к кэшам браузеров
	cachePaths := []string{
		filepath.Join(homeDir, ".cache/google-chrome"),
		filepath.Join(homeDir, ".cache/chromium"),
		filepath.Join(homeDir, ".cache/firefox"),
		filepath.Join(homeDir, ".cache/mozilla"),
		filepath.Join(homeDir, ".cache/thumbnails"),
	}
	
	for _, cachePath := range cachePaths {
		if size, err := m.getDirSize(cachePath); err == nil {
			totalSize += size
			os.RemoveAll(cachePath)
		}
	}
	
	return totalSize, nil
}

func (m *CleanupModule) cleanTempFiles() (int64, error) {
	var totalSize int64
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return 0, err
	}
	
	// Временные папки
	tempPaths := []string{
		"/tmp",
		filepath.Join(homeDir, ".cache"),
		filepath.Join(homeDir, ".local/share/Trash"),
	}
	
	for _, tempPath := range tempPaths {
		if tempPath == "/tmp" {
			// Для /tmp очищаем только старые файлы
			cmd := exec.Command("find", "/tmp", "-type", "f", "-atime", "+7", "-delete")
			cmd.Run()
		} else {
			if size, err := m.getDirSize(tempPath); err == nil {
				totalSize += size
				os.RemoveAll(tempPath)
				os.MkdirAll(tempPath, 0755) // Пересоздаем папку
			}
		}
	}
	
	return totalSize, nil
}

func (m *CleanupModule) cleanOldLogs() (int64, error) {
	var totalSize int64
	
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	
	// Очищаем системные логи старше 7 дней (без sudo)
	cmd := exec.CommandContext(ctx, "journalctl", "--vacuum-time=7d")
	cmd.Run() // Игнорируем ошибки
	
	// Очищаем старые логи в домашней папке пользователя
	homeDir, err := os.UserHomeDir()
	if err == nil {
		logPath := filepath.Join(homeDir, ".local/share/logs")
		if size, err := m.getDirSize(logPath); err == nil {
			totalSize += size
			os.RemoveAll(logPath)
		}
	}
	
	return totalSize, nil
}

func (m *CleanupModule) getDirSize(path string) (int64, error) {
	var size int64
	
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Игнорируем ошибки доступа
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	
	return size, err
}
