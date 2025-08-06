package modules

import (
	"fmt"
	"os/exec"
	"strings"
)

type UpdatesModule struct{}

func (m *UpdatesModule) GetName() string {
	return "System Updates"
}

func (m *UpdatesModule) GetDescription() string {
	return "Update system packages and security patches"
}

func (m *UpdatesModule) RequiresRoot() bool {
	return true
}

func (m *UpdatesModule) Execute(progressCallback func(progress float64, message string)) error {
	progressCallback(0.1, "Updating package lists...")
	
	// Обновляем списки пакетов
	cmd := exec.Command("sudo", "apt", "update")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update package lists: %v", err)
	}
	
	progressCallback(0.3, "Checking for upgradeable packages...")
	
	// Проверяем доступные обновления
	cmd = exec.Command("apt", "list", "--upgradable")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check upgradeable packages: %v", err)
	}
	
	upgradeable := strings.Count(string(output), "\n") - 1 // -1 для заголовка
	if upgradeable <= 0 {
		progressCallback(1.0, "No packages to update")
		return nil
	}
	
	progressCallback(0.5, fmt.Sprintf("Found %d upgradeable packages", upgradeable))
	
	// Выполняем обновление
	progressCallback(0.6, "Installing package updates...")
	cmd = exec.Command("sudo", "apt", "upgrade", "-y")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to upgrade packages: %v", err)
	}
	
	progressCallback(0.8, "Checking for snap updates...")
	
	// Обновляем snap пакеты
	cmd = exec.Command("sudo", "snap", "refresh")
	cmd.Run() // Игнорируем ошибки snap
	
	progressCallback(0.9, "Cleaning up...")
	
	// Очищаем кэш
	cmd = exec.Command("sudo", "apt", "autoremove", "-y")
	cmd.Run()
	
	progressCallback(1.0, fmt.Sprintf("Successfully updated %d packages", upgradeable))
	
	return nil
}