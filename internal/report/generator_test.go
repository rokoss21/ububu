package report

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestGenerator_New(t *testing.T) {
	gen := NewGenerator()
	if gen == nil {
		t.Error("NewGenerator() returned nil")
	}
}

func TestGenerator_GenerateFileName(t *testing.T) {
	gen := NewGenerator()
	filename := gen.generateFileName()
	
	// Проверяем, что имя файла не пустое
	if filename == "" {
		t.Error("generateFileName() returned empty string")
	}
	
	// Проверяем, что имя содержит префикс
	if !strings.HasPrefix(filename, "ububu-report-") {
		t.Errorf("Filename should start with 'ububu-report-', got: %s", filename)
	}
	
	// Проверяем, что имя содержит расширение .log
	if !strings.HasSuffix(filename, ".log") {
		t.Errorf("Filename should end with '.log', got: %s", filename)
	}
	
	// Проверяем, что имя содержит дату
	now := time.Now()
	expectedDate := now.Format("2006-01-02")
	if !strings.Contains(filename, expectedDate) {
		t.Errorf("Filename should contain today's date %s, got: %s", expectedDate, filename)
	}
}

func TestGenerator_GenerateFullReport(t *testing.T) {
	gen := NewGenerator()
	testLogContent := "Test log content\nLine 2\nLine 3"
	
	report := gen.GenerateFullReport(testLogContent)
	
	// Проверяем, что отчет не пустой
	if report == "" {
		t.Error("generateFullReport() returned empty string")
	}
	
	// Проверяем, что отчет содержит заголовок
	if !strings.Contains(report, "Ububu 1.0") {
		t.Error("Report should contain 'Ububu 1.0'")
	}
	
	// Проверяем, что отчет содержит дату
	now := time.Now()
	expectedDate := now.Format("2006-01-02")
	if !strings.Contains(report, expectedDate) {
		t.Errorf("Report should contain today's date %s", expectedDate)
	}
	
	// Проверяем, что отчет содержит переданный лог
	if !strings.Contains(report, testLogContent) {
		t.Error("Report should contain the provided log content")
	}
	
	// Проверяем, что отчет содержит системную информацию
	if !strings.Contains(report, "SYSTEM INFORMATION") {
		t.Error("Report should contain system information section")
	}
	
	// Проверяем, что отчет содержит сводку
	if !strings.Contains(report, "SUMMARY") {
		t.Error("Report should contain summary section")
	}
	
	// Проверяем, что отчет содержит лог секцию
	if !strings.Contains(report, "OPTIMIZATION LOG") {
		t.Error("Report should contain optimization log section")
	}
}

func TestGenerator_GetSystemInfo(t *testing.T) {
	gen := NewGenerator()
	sysInfo := gen.getSystemInfo()
	
	// Проверяем, что системная информация не пустая
	if sysInfo == "" {
		t.Error("getSystemInfo() returned empty string")
	}
	
	// Проверяем, что содержит информацию о пользователе
	if user := os.Getenv("USER"); user != "" {
		if !strings.Contains(sysInfo, user) {
			t.Errorf("System info should contain username %s", user)
		}
	}
	
	// Проверяем, что содержит hostname
	if hostname, err := os.Hostname(); err == nil {
		if !strings.Contains(sysInfo, hostname) {
			t.Errorf("System info should contain hostname %s", hostname)
		}
	}
}

func TestGenerator_GenerateFullReport_WithEmptyLog(t *testing.T) {
	gen := NewGenerator()
	report := gen.GenerateFullReport("")
	
	// Даже с пустым логом отчет должен содержать структуру
	if !strings.Contains(report, "Ububu 1.0") {
		t.Error("Report should contain header even with empty log")
	}
	
	if !strings.Contains(report, "SYSTEM INFORMATION") {
		t.Error("Report should contain system information section even with empty log")
	}
}

func TestGenerator_GenerateFullReport_WithSpecialCharacters(t *testing.T) {
	gen := NewGenerator()
	testLogContent := "Test with special chars: áéíóú ñ 中文 🐧 <>&\""
	
	report := gen.GenerateFullReport(testLogContent)
	
	// Проверяем, что специальные символы корректно обрабатываются
	if !strings.Contains(report, testLogContent) {
		t.Error("Report should handle special characters correctly")
	}
}

func TestGenerator_GenerateFullReport_Structure(t *testing.T) {
	gen := NewGenerator()
	testLogContent := "Sample log content"
	
	report := gen.GenerateFullReport(testLogContent)
	
	// Проверяем правильную структуру отчета
	sections := []string{
		"Ububu 1.0 - System Optimization Report",
		"SYSTEM INFORMATION",
		"OPTIMIZATION LOG",
		"SUMMARY",
		"End of Report",
	}
	
	for _, section := range sections {
		if !strings.Contains(report, section) {
			t.Errorf("Report should contain section: %s", section)
		}
	}
	
	// Проверяем, что секции идут в правильном порядке
	lastIndex := -1
	for _, section := range sections {
		index := strings.Index(report, section)
		if index <= lastIndex {
			t.Errorf("Section '%s' should appear after previous sections", section)
		}
		lastIndex = index
	}
}

// Бенчмарк для генерации отчета
func BenchmarkGenerator_GenerateFullReport(b *testing.B) {
	gen := NewGenerator()
	testLogContent := strings.Repeat("Test log line\n", 1000) // 1000 строк лога
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gen.GenerateFullReport(testLogContent)
	}
}

// Бенчмарк для генерации имени файла
func BenchmarkGenerator_GenerateFileName(b *testing.B) {
	gen := NewGenerator()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gen.generateFileName()
	}
}