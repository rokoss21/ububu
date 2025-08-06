package test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rokoss21/ububu/internal/modules"
	"github.com/rokoss21/ububu/internal/report"
)

// TestModuleIntegration тестирует интеграцию между модулями
func TestModuleIntegration(t *testing.T) {
	// Создаем все модули
	allModules := []modules.SystemModule{
		&modules.HealthModule{},
		&modules.CleanupModule{},
	}
	
	// Симулируем выполнение всех модулей
	var allLogs []string
	
	for _, module := range allModules {
		t.Logf("Testing module: %s", module.GetName())
		
		if module.RequiresRoot() && os.Geteuid() != 0 {
			t.Logf("Skipping %s - requires root", module.GetName())
			continue
		}
		
		var moduleLogs []string
		progressCallback := func(progress float64, message string) {
			if message != "" {
				logEntry := time.Now().Format("15:04:05") + " " + message
				moduleLogs = append(moduleLogs, logEntry)
				allLogs = append(allLogs, logEntry)
			}
		}
		
		err := module.Execute(progressCallback)
		if err != nil {
			t.Errorf("Module %s failed: %v", module.GetName(), err)
		}
		
		if len(moduleLogs) == 0 {
			t.Errorf("Module %s produced no log messages", module.GetName())
		}
	}
	
	// Проверяем, что у нас есть логи для генерации отчета
	if len(allLogs) == 0 {
		t.Error("No logs generated from any modules")
	}
	
	// Тестируем генерацию отчета с реальными логами
	gen := report.NewGenerator()
	logContent := strings.Join(allLogs, "\n")
	reportContent := gen.GenerateFullReport(logContent)
	
	// Проверяем, что отчет содержит логи модулей
	for _, log := range allLogs {
		if !strings.Contains(reportContent, log) {
			t.Errorf("Report should contain log entry: %s", log)
		}
	}
}

// TestSystemHealthWorkflow тестирует полный workflow проверки здоровья системы
func TestSystemHealthWorkflow(t *testing.T) {
	healthModule := &modules.HealthModule{}
	
	// Собираем детальную информацию о выполнении
	var steps []string
	var progressValues []float64
	
	progressCallback := func(progress float64, message string) {
		progressValues = append(progressValues, progress)
		if message != "" {
			steps = append(steps, message)
		}
	}
	
	err := healthModule.Execute(progressCallback)
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
	
	// Проверяем, что прогресс монотонно возрастает
	for i := 1; i < len(progressValues); i++ {
		if progressValues[i] < progressValues[i-1] {
			t.Errorf("Progress should be monotonic: %f -> %f", progressValues[i-1], progressValues[i])
		}
	}
	
	// Проверяем, что финальный прогресс = 1.0
	if len(progressValues) > 0 && progressValues[len(progressValues)-1] != 1.0 {
		t.Errorf("Final progress should be 1.0, got %f", progressValues[len(progressValues)-1])
	}
	
	// Проверяем, что все основные проверки выполнены
	expectedChecks := []string{"disk", "memory", "process"}
	for _, check := range expectedChecks {
		found := false
		for _, step := range steps {
			if strings.Contains(strings.ToLower(step), check) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Health check should include %s check", check)
		}
	}
}

// TestCleanupWorkflow тестирует полный workflow очистки системы
func TestCleanupWorkflow(t *testing.T) {
	cleanupModule := &modules.CleanupModule{}
	
	var cleanupSteps []string
	var totalFreed int64
	
	progressCallback := func(progress float64, message string) {
		if message != "" {
			cleanupSteps = append(cleanupSteps, message)
			
			// Пытаемся извлечь информацию об освобожденном месте
			if strings.Contains(message, "MB freed") {
				// Простой парсинг для тестирования
				if strings.Contains(message, "Total freed:") {
					// Это финальное сообщение
				}
			}
		}
	}
	
	err := cleanupModule.Execute(progressCallback)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}
	
	// Проверяем, что все типы очистки были выполнены
	expectedCleanupTypes := []string{"package", "browser", "temp", "log"}
	for _, cleanupType := range expectedCleanupTypes {
		found := false
		for _, step := range cleanupSteps {
			if strings.Contains(strings.ToLower(step), cleanupType) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Cleanup should include %s cleanup", cleanupType)
		}
	}
	
	// Проверяем, что есть финальное сообщение о завершении
	if len(cleanupSteps) > 0 {
		lastStep := cleanupSteps[len(cleanupSteps)-1]
		if !strings.Contains(strings.ToLower(lastStep), "completed") {
			t.Error("Cleanup should end with completion message")
		}
	}
	
	t.Logf("Cleanup completed with %d steps, freed %d bytes", len(cleanupSteps), totalFreed)
}

// TestReportGeneration тестирует генерацию отчета с реальными данными
func TestReportGeneration(t *testing.T) {
	// Создаем реалистичные логи
	testLogs := []string{
		"[INFO] 14:30:00 Starting system optimization",
		"[PROGRESS] 14:30:05 Checking disk health...",
		"[SUCCESS] 14:30:10 Disk health: GOOD (15% used)",
		"[PROGRESS] 14:30:15 Analyzing memory usage...",
		"[SUCCESS] 14:30:20 Memory usage: NORMAL (45% used)",
		"[PROGRESS] 14:30:25 Cleaning package cache...",
		"[SUCCESS] 14:30:30 Package cache cleaned: 256 MB freed",
		"[INFO] 14:30:35 Optimization completed successfully",
	}
	
	logContent := strings.Join(testLogs, "\n")
	
	gen := report.NewGenerator()
	reportContent := gen.GenerateFullReport(logContent)
	
	// Проверяем структуру отчета
	requiredSections := []string{
		"System Optimization Report",
		"SYSTEM INFORMATION",
		"OPTIMIZATION LOG",
		"SUMMARY",
	}
	
	for _, section := range requiredSections {
		if !strings.Contains(reportContent, section) {
			t.Errorf("Report missing required section: %s", section)
		}
	}
	
	// Проверяем, что все логи включены в отчет
	for _, log := range testLogs {
		if !strings.Contains(reportContent, log) {
			t.Errorf("Report missing log entry: %s", log)
		}
	}
	
	// Проверяем размер отчета (должен быть разумным)
	if len(reportContent) < 1000 {
		t.Error("Report seems too short")
	}
	
	if len(reportContent) > 100000 {
		t.Error("Report seems too long")
	}
}

// TestModuleErrorHandling тестирует обработку ошибок в модулях
func TestModuleErrorHandling(t *testing.T) {
	// Тестируем модули, которые могут вернуть ошибки
	modules := []modules.SystemModule{
		&modules.UpdatesModule{}, // Требует root
		&modules.DriversModule{}, // Требует root
	}
	
	for _, module := range modules {
		if !module.RequiresRoot() {
			continue // Пропускаем модули, которые не требуют root
		}
		
		if os.Geteuid() == 0 {
			continue // Пропускаем если мы уже root
		}
		
		t.Logf("Testing error handling for: %s", module.GetName())
		
		callbackCalled := false
		progressCallback := func(progress float64, message string) {
			callbackCalled = true
		}
		
		err := module.Execute(progressCallback)
		
		// Ожидаем ошибку для модулей, требующих root
		if err == nil {
			t.Errorf("Module %s should return error when not running as root", module.GetName())
		}
		
		// Callback все равно должен вызываться для начальных шагов
		if !callbackCalled {
			t.Errorf("Module %s should call progress callback even when failing", module.GetName())
		}
	}
}
