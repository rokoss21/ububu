package modules

import (
	"fmt"
	"os/exec"
	"strings"
)

type DriversModule struct{}

func (m *DriversModule) GetName() string {
	return "Driver Updates"
}

func (m *DriversModule) GetDescription() string {
	return "Check and update system drivers"
}

func (m *DriversModule) RequiresRoot() bool {
	return true
}

func (m *DriversModule) Execute(progressCallback func(progress float64, message string)) error {
	progressCallback(0.1, "Detecting available drivers...")
	
	// Проверяем доступные драйверы
	cmd := exec.Command("ubuntu-drivers", "devices")
	output, err := cmd.Output()
	if err != nil {
		progressCallback(0.5, "ubuntu-drivers not available, checking manually...")
		return m.checkManualDrivers(progressCallback)
	}
	
	progressCallback(0.3, "Analyzing driver recommendations...")
	
	outputStr := string(output)
	if strings.Contains(outputStr, "No devices") || len(strings.TrimSpace(outputStr)) == 0 {
		progressCallback(1.0, "No additional drivers needed")
		return nil
	}
	
	progressCallback(0.5, "Installing recommended drivers...")
	
	// Устанавливаем рекомендуемые драйверы
	cmd = exec.Command("sudo", "ubuntu-drivers", "autoinstall")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install drivers: %v", err)
	}
	
	progressCallback(0.8, "Checking NVIDIA drivers...")
	
	// Проверяем NVIDIA драйверы отдельно
	cmd = exec.Command("nvidia-smi")
	if err := cmd.Run(); err == nil {
		progressCallback(0.9, "NVIDIA drivers are working correctly")
	}
	
	progressCallback(1.0, "Driver updates completed")
	
	return nil
}

func (m *DriversModule) checkManualDrivers(progressCallback func(progress float64, message string)) error {
	progressCallback(0.4, "Checking for NVIDIA hardware...")
	
	// Проверяем наличие NVIDIA карты
	cmd := exec.Command("lspci")
	output, err := cmd.Output()
	if err != nil {
		progressCallback(1.0, "Could not detect hardware")
		return nil
	}
	
	outputStr := strings.ToLower(string(output))
	
	if strings.Contains(outputStr, "nvidia") {
		progressCallback(0.6, "NVIDIA hardware detected")
		
		// Проверяем установлен ли драйвер
		cmd = exec.Command("nvidia-smi")
		if err := cmd.Run(); err != nil {
			progressCallback(0.8, "NVIDIA driver not installed or not working")
			// Здесь можно добавить установку драйвера
		} else {
			progressCallback(0.8, "NVIDIA driver is working")
		}
	}
	
	if strings.Contains(outputStr, "amd") || strings.Contains(outputStr, "radeon") {
		progressCallback(0.7, "AMD hardware detected")
		// AMD драйверы обычно включены в ядро
	}
	
	progressCallback(1.0, "Hardware check completed")
	return nil
}