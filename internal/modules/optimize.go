package modules

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

type OptimizationModule struct{}

func (m *OptimizationModule) GetName() string {
	return "System Optimization"
}

func (m *OptimizationModule) GetDescription() string {
	return "Optimize system performance settings"
}

func (m *OptimizationModule) RequiresRoot() bool {
	return true
}

func (m *OptimizationModule) Execute(progressCallback func(progress float64, message string)) error {
	progressCallback(0.1, "Checking SSD optimization...")
	
	if err := m.optimizeSSD(progressCallback); err != nil {
		progressCallback(0.3, fmt.Sprintf("SSD optimization failed: %v", err))
	} else {
		progressCallback(0.3, "SSD optimization completed")
	}
	
	progressCallback(0.4, "Optimizing memory settings...")
	
	if err := m.optimizeMemory(progressCallback); err != nil {
		progressCallback(0.6, fmt.Sprintf("Memory optimization failed: %v", err))
	} else {
		progressCallback(0.6, "Memory settings optimized")
	}
	
	progressCallback(0.7, "Clearing network cache...")
	
	if err := m.clearNetworkCache(progressCallback); err != nil {
		progressCallback(0.9, fmt.Sprintf("Network cache clear failed: %v", err))
	} else {
		progressCallback(0.9, "Network cache cleared")
	}
	
	progressCallback(1.0, "System optimization completed")
	
	return nil
}

func (m *OptimizationModule) optimizeSSD(progressCallback func(progress float64, message string)) error {
	// Проверяем есть ли SSD диски
	cmd := exec.Command("lsblk", "-d", "-o", "name,rota")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	
	lines := strings.Split(string(output), "\n")
	hasSSD := false
	
	for _, line := range lines[1:] { // Пропускаем заголовок
		if strings.Contains(line, "0") { // 0 означает SSD
			hasSSD = true
			break
		}
	}
	
	if !hasSSD {
		progressCallback(0.2, "No SSD detected, skipping SSD optimization")
		return nil
	}
	
	progressCallback(0.15, "SSD detected, running TRIM...")
	
	// Выполняем TRIM для всех SSD
	cmd = exec.Command("sudo", "fstrim", "-av")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("TRIM failed: %v", err)
	}
	
	progressCallback(0.25, "TRIM completed successfully")
	
	return nil
}

func (m *OptimizationModule) optimizeMemory(progressCallback func(progress float64, message string)) error {
	progressCallback(0.45, "Checking current swappiness...")
	
	// Читаем текущее значение swappiness
	data, err := ioutil.ReadFile("/proc/sys/vm/swappiness")
	if err != nil {
		return err
	}
	
	currentSwappiness, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return err
	}
	
	progressCallback(0.5, fmt.Sprintf("Current swappiness: %d", currentSwappiness))
	
	// Оптимальное значение для десктопа - 10
	optimalSwappiness := 10
	
	if currentSwappiness != optimalSwappiness {
		progressCallback(0.55, fmt.Sprintf("Setting swappiness to %d...", optimalSwappiness))
		
		// Устанавливаем новое значение
		cmd := exec.Command("sudo", "sysctl", fmt.Sprintf("vm.swappiness=%d", optimalSwappiness))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set swappiness: %v", err)
		}
		
		// Делаем изменение постоянным
		cmd = exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo 'vm.swappiness=%d' >> /etc/sysctl.conf", optimalSwappiness))
		cmd.Run() // Игнорируем ошибки, возможно уже есть
	}
	
	return nil
}

func (m *OptimizationModule) clearNetworkCache(progressCallback func(progress float64, message string)) error {
	progressCallback(0.75, "Flushing DNS cache...")
	
	// Очищаем DNS кэш
	cmd := exec.Command("sudo", "systemctl", "flush-dns")
	if err := cmd.Run(); err != nil {
		// Пробуем альтернативный способ
		cmd = exec.Command("sudo", "systemd-resolve", "--flush-caches")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to flush DNS cache: %v", err)
		}
	}
	
	progressCallback(0.85, "Clearing network manager cache...")
	
	// Перезапускаем NetworkManager для очистки кэша
	cmd = exec.Command("sudo", "systemctl", "restart", "NetworkManager")
	if err := cmd.Run(); err != nil {
		// Не критично, продолжаем
		progressCallback(0.87, "NetworkManager restart failed, continuing...")
	}
	
	return nil
}