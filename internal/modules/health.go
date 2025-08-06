package modules

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

type HealthModule struct{}

func (m *HealthModule) GetName() string {
	return "System Health Check"
}

func (m *HealthModule) GetDescription() string {
	return "Check system health and performance metrics"
}

func (m *HealthModule) RequiresRoot() bool {
	return false
}

func (m *HealthModule) Execute(progressCallback func(progress float64, message string)) error {
	progressCallback(0.1, "Checking disk health...")
	
	diskHealth, err := m.checkDiskHealth()
	if err != nil {
		progressCallback(0.25, fmt.Sprintf("Disk health check failed: %v", err))
	} else {
		progressCallback(0.25, diskHealth)
	}
	
	progressCallback(0.3, "Checking system temperature...")
	
	tempInfo, err := m.checkTemperature()
	if err != nil {
		progressCallback(0.5, fmt.Sprintf("Temperature check failed: %v", err))
	} else {
		progressCallback(0.5, tempInfo)
	}
	
	progressCallback(0.6, "Analyzing memory usage...")
	
	memInfo, err := m.checkMemoryUsage()
	if err != nil {
		progressCallback(0.75, fmt.Sprintf("Memory check failed: %v", err))
	} else {
		progressCallback(0.75, memInfo)
	}
	
	progressCallback(0.8, "Analyzing running processes...")
	
	procInfo, err := m.analyzeProcesses()
	if err != nil {
		progressCallback(0.95, fmt.Sprintf("Process analysis failed: %v", err))
	} else {
		progressCallback(0.95, procInfo)
	}
	
	progressCallback(1.0, "System health check completed")
	
	return nil
}

func (m *HealthModule) checkDiskHealth() (string, error) {
	// Проверяем доступное место на диске
	cmd := exec.Command("df", "-h", "/")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected df output")
	}
	
	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return "", fmt.Errorf("unexpected df fields")
	}
	
	used := fields[4] // Процент использования
	usedPercent, err := strconv.Atoi(strings.TrimSuffix(used, "%"))
	if err != nil {
		return "", err
	}
	
	var status string
	if usedPercent > 90 {
		status = "⚠️ WARNING: Disk usage is high"
	} else if usedPercent > 80 {
		status = "⚠️ CAUTION: Disk usage is moderate"
	} else {
		status = "✅ GOOD: Disk usage is normal"
	}
	
	// Пробуем проверить SMART статус
	cmd = exec.Command("sudo", "smartctl", "-H", "/dev/sda")
	smartOutput, err := cmd.Output()
	if err == nil && strings.Contains(string(smartOutput), "PASSED") {
		status += " • SMART: PASSED"
	}
	
	return fmt.Sprintf("%s (%s used)", status, used), nil
}

func (m *HealthModule) checkTemperature() (string, error) {
	// Проверяем температуру CPU
	tempFiles := []string{
		"/sys/class/thermal/thermal_zone0/temp",
		"/sys/class/thermal/thermal_zone1/temp",
	}
	
	var temps []int
	for _, file := range tempFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		
		tempStr := strings.TrimSpace(string(data))
		temp, err := strconv.Atoi(tempStr)
		if err != nil {
			continue
		}
		
		// Температура в милли-градусах, конвертируем в градусы
		temps = append(temps, temp/1000)
	}
	
	if len(temps) == 0 {
		return "Temperature sensors not available", nil
	}
	
	// Берем максимальную температуру
	maxTemp := temps[0]
	for _, temp := range temps {
		if temp > maxTemp {
			maxTemp = temp
		}
	}
	
	var status string
	if maxTemp > 80 {
		status = "🔥 HOT"
	} else if maxTemp > 70 {
		status = "⚠️ WARM"
	} else {
		status = "❄️ COOL"
	}
	
	return fmt.Sprintf("%s: CPU temperature %d°C", status, maxTemp), nil
}

func (m *HealthModule) checkMemoryUsage() (string, error) {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(string(data), "\n")
	var memTotal, memAvailable int
	
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		
		value, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		
		switch fields[0] {
		case "MemTotal:":
			memTotal = value
		case "MemFree:":
			// memFree не используется, но оставляем для полноты
		case "MemAvailable:":
			memAvailable = value
		}
	}
	
	if memTotal == 0 {
		return "", fmt.Errorf("could not parse memory info")
	}
	
	usedPercent := (memTotal - memAvailable) * 100 / memTotal
	
	var status string
	if usedPercent > 90 {
		status = "🔴 CRITICAL"
	} else if usedPercent > 80 {
		status = "⚠️ HIGH"
	} else {
		status = "✅ NORMAL"
	}
	
	return fmt.Sprintf("%s: Memory usage %d%% (%d MB / %d MB)", 
		status, usedPercent, (memTotal-memAvailable)/1024, memTotal/1024), nil
}

func (m *HealthModule) analyzeProcesses() (string, error) {
	// Подсчитываем количество процессов
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(string(output), "\n")
	processCount := len(lines) - 2 // Убираем заголовок и пустую строку
	
	// Проверяем загрузку системы
	data, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return "", err
	}
	
	loadFields := strings.Fields(string(data))
	if len(loadFields) < 1 {
		return "", fmt.Errorf("could not parse load average")
	}
	
	load1min, err := strconv.ParseFloat(loadFields[0], 64)
	if err != nil {
		return "", err
	}
	
	var loadStatus string
	if load1min > 2.0 {
		loadStatus = "🔴 HIGH"
	} else if load1min > 1.0 {
		loadStatus = "⚠️ MODERATE"
	} else {
		loadStatus = "✅ LOW"
	}
	
	return fmt.Sprintf("%s load (%.2f) • %d active processes", 
		loadStatus, load1min, processCount), nil
}