package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rokoss21/ububu/internal/modules"
)

// Компактные стили для стандартного терминала 80x24
var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D4AA")).
		Bold(true)

	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	taskStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(0, 1)

	selectedTaskStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00D4AA")).
		Padding(0, 1)

	logStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)
)

type Task struct {
	Name        string
	Description string
	Icon        string
	Module      modules.SystemModule
	Selected    bool
	Progress    float64
	Status      string
	Error       error
	StartTime   time.Time
	EndTime     time.Time
	Details     []string
}

type model struct {
	tasks         []Task
	cursor        int
	running       bool
	currentTask   int
	progress      progress.Model
	spinner       spinner.Model
	logs          []string
	width         int
	height        int
	phase         string // "select", "running", "complete", "report"
	totalTasks    int
	completedTasks int
	overallProgress float64
	reportGenerated bool
}

type taskCompleteMsg struct {
	taskIndex int
	success   bool
	message   string
	error     error
	startTime time.Time
	endTime   time.Time
	details   []string
}

type progressMsg struct {
	taskIndex int
	progress  float64
	message   string
}

type logMsg struct {
	level   string
	message string
}

func initialModel() model {
	// Компактные задачи для стандартного терминала
	tasks := []Task{
		{
			Name:        "Health Check",
			Description: "System health & performance",
			Icon:        "🏥",
			Module:      &modules.HealthModule{},
			Selected:    true,
		},
		{
			Name:        "Cleanup",
			Description: "Clean temp files & caches",
			Icon:        "🧹",
			Module:      &modules.CleanupModule{},
			Selected:    true,
		},
		{
			Name:        "Updates",
			Description: "System & security updates",
			Icon:        "🔄",
			Module:      &modules.UpdatesModule{},
			Selected:    false,
		},
		{
			Name:        "Drivers",
			Description: "Driver updates",
			Icon:        "🖥️",
			Module:      &modules.DriversModule{},
			Selected:    false,
		},
		{
			Name:        "Optimization",
			Description: "Performance optimization",
			Icon:        "⚡",
			Module:      &modules.OptimizationModule{},
			Selected:    false,
		},
	}

	// Компактный прогресс-бар
	prog := progress.New(progress.WithDefaultGradient())
	prog.Width = 50

	// Компактный спиннер
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D4AA"))

	return model{
		tasks:   tasks,
		cursor:  0,
		phase:   "select",
		progress: prog,
		spinner: s,
		logs:    []string{},
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if msg.Width > 60 {
			m.progress.Width = msg.Width - 20
		} else {
			m.progress.Width = 40
		}
		return m, nil

	case tea.KeyMsg:
		switch m.phase {
		case "select":
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.tasks)-1 {
					m.cursor++
				}
			case " ":
				m.tasks[m.cursor].Selected = !m.tasks[m.cursor].Selected
			case "enter":
				return m.startTasks()
			case "a":
				for i := range m.tasks {
					m.tasks[i].Selected = true
				}
			case "n":
				for i := range m.tasks {
					m.tasks[i].Selected = false
				}
			}
		case "running":
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "p":
				return m.generateReport()
			}
		case "complete":
			switch msg.String() {
			case "ctrl+c", "q", "enter":
				return m, tea.Quit
			case "p":
				return m.generateReport()
			}
		case "report":
			switch msg.String() {
			case "ctrl+c", "q", "enter":
				return m, tea.Quit
			}
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case progressMsg:
		if msg.taskIndex < len(m.tasks) {
			m.tasks[msg.taskIndex].Progress = msg.progress
			cmd := m.progress.SetPercent(msg.progress)
			m.addLog("PROGRESS", msg.message)
			return m, cmd
		}
		return m, nil

	case taskCompleteMsg:
		if msg.taskIndex < len(m.tasks) {
			m.tasks[msg.taskIndex].Progress = 1.0
			m.tasks[msg.taskIndex].Error = msg.error
			m.tasks[msg.taskIndex].StartTime = msg.startTime
			m.tasks[msg.taskIndex].EndTime = msg.endTime
			m.tasks[msg.taskIndex].Details = msg.details
			
			if msg.success {
				m.tasks[msg.taskIndex].Status = "✅ Complete"
			} else {
				m.tasks[msg.taskIndex].Status = "❌ Failed"
			}
			m.addLog("INFO", msg.message)
			
			// Обновляем общий прогресс
			m.completedTasks++
			m.overallProgress = float64(m.completedTasks) / float64(m.totalTasks)
			cmd := m.progress.SetPercent(m.overallProgress)
			
			// Переходим к следующей задаче
			m.currentTask++
			if m.currentTask >= len(m.getSelectedTasks()) {
				m.phase = "complete"
				m.running = false
				m.addLog("SUCCESS", "🎉 All tasks completed!")
				return m, cmd
			} else {
				// Запускаем следующую задачу
				return m, tea.Batch(cmd, m.executeNextTask())
			}
		}
		return m, nil

	case logMsg:
		m.addLog(msg.level, msg.message)
		return m, nil
	}

	return m, cmd
}

func (m *model) addLog(level, message string) {
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("[%s] %s", timestamp, message)
	m.logs = append(m.logs, logEntry)
	
	// Ограничиваем количество логов для компактности
	if len(m.logs) > 20 {
		m.logs = m.logs[len(m.logs)-20:]
	}
}

func (m model) getSelectedTasks() []int {
	var selected []int
	for i, task := range m.tasks {
		if task.Selected {
			selected = append(selected, i)
		}
	}
	return selected
}

func (m model) startTasks() (tea.Model, tea.Cmd) {
	selectedTasks := m.getSelectedTasks()
	if len(selectedTasks) == 0 {
		m.addLog("ERROR", "No tasks selected!")
		return m, nil
	}

	m.phase = "running"
	m.running = true
	m.currentTask = 0
	m.totalTasks = len(selectedTasks)
	m.completedTasks = 0
	m.overallProgress = 0.0
	m.addLog("INFO", "🚀 Starting Ububu optimization...")

	// Инициализируем время начала для выбранных задач
	for _, taskIndex := range selectedTasks {
		m.tasks[taskIndex].StartTime = time.Now()
	}

	// Запускаем выполнение задач
	return m, tea.Batch(
		m.progress.Init(),
		m.executeNextTask(),
	)
}

func (m model) executeNextTask() tea.Cmd {
	selectedTasks := m.getSelectedTasks()
	if m.currentTask >= len(selectedTasks) {
		return nil
	}

	taskIndex := selectedTasks[m.currentTask]
	task := m.tasks[taskIndex]

	return func() tea.Msg {
		// Отмечаем время начала
		startTime := time.Now()
		var details []string
		
		// Выполняем задачу синхронно с callback для прогресса
		err := task.Module.Execute(func(progress float64, message string) {
			details = append(details, fmt.Sprintf("%.0f%% - %s", progress*100, message))
			// Прогресс будет обновляться через общий прогресс задач
		})

		endTime := time.Now()
		success := err == nil
		var message string
		if success {
			message = fmt.Sprintf("✅ %s completed successfully", task.Name)
		} else {
			message = fmt.Sprintf("❌ %s failed: %v", task.Name, err)
		}

		return taskCompleteMsg{
			taskIndex: taskIndex,
			success:   success,
			message:   message,
			error:     err,
			startTime: startTime,
			endTime:   endTime,
			details:   details,
		}
	}
}

func (m model) View() string {
	var b strings.Builder

	// Компактный заголовок
	title := titleStyle.Render("🐧 Ububu 1.0 - Ubuntu System Optimizer")
	b.WriteString(title + "\n")
	b.WriteString(strings.Repeat("=", 50) + "\n\n")

	switch m.phase {
	case "select":
		b.WriteString(m.renderTaskSelection())
	case "running":
		b.WriteString(m.renderRunning())
	case "complete":
		b.WriteString(m.renderComplete())
	case "report":
		b.WriteString(m.renderReport())
	}

	return b.String()
}

func (m model) renderTaskSelection() string {
	var b strings.Builder

	b.WriteString(headerStyle.Render("📋 Select Tasks:") + "\n\n")

	for i, task := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checkbox := "☐"
		if task.Selected {
			checkbox = "☑"
		}

		// Компактный формат - одна строка на задачу
		taskText := fmt.Sprintf("%s %s %s %s - %s", 
			cursor, checkbox, task.Icon, task.Name, task.Description)

		if m.cursor == i {
			b.WriteString(selectedTaskStyle.Render(taskText) + "\n")
		} else {
			b.WriteString(taskStyle.Render(taskText) + "\n")
		}
	}

	// Компактные инструкции
	b.WriteString("\n" + headerStyle.Render("Controls:") + " ↑/↓ Navigate • Space Toggle • Enter Start • q Quit\n")

	return b.String()
}

func (m model) renderRunning() string {
	var b strings.Builder

	b.WriteString(headerStyle.Render("🚀 Running Optimization...") + "\n\n")

	// Показываем текущую задачу и общий прогресс
	selectedTasks := m.getSelectedTasks()
	if m.currentTask < len(selectedTasks) {
		currentTaskIndex := selectedTasks[m.currentTask]
		currentTask := m.tasks[currentTaskIndex]
		
		b.WriteString(fmt.Sprintf("%s %s %s\n", 
			m.spinner.View(), currentTask.Icon, currentTask.Name))
	}
	
	// Общий прогресс всех задач
	b.WriteString(fmt.Sprintf("Overall Progress: %d/%d tasks (%.1f%%)\n", 
		m.completedTasks, m.totalTasks, m.overallProgress*100))
	b.WriteString(m.progress.View() + "\n\n")

	// Компактный статус задач
	b.WriteString(headerStyle.Render("📊 Status:") + "\n")
	for i, task := range m.tasks {
		if !task.Selected {
			continue
		}

		status := "⏳ Pending"
		if i < len(selectedTasks) && selectedTasks[m.currentTask] == i {
			status = "🔄 Running"
		} else if task.Status != "" {
			status = task.Status
		}

		b.WriteString(fmt.Sprintf("  %s %s - %s\n", task.Icon, task.Name, status))
	}

	// Компактные логи (только последние 3)
	if len(m.logs) > 0 {
		b.WriteString("\n" + headerStyle.Render("📝 Log:") + "\n")
		start := 0
		if len(m.logs) > 3 {
			start = len(m.logs) - 3
		}
		for _, log := range m.logs[start:] {
			b.WriteString(logStyle.Render("  " + log) + "\n")
		}
	}

	b.WriteString("\n" + headerStyle.Render("Press p for report • q to quit") + "\n")

	return b.String()
}

func (m model) renderComplete() string {
	var b strings.Builder

	b.WriteString(successStyle.Render("🎉 Optimization Complete!") + "\n\n")

	// Компактные результаты
	b.WriteString(headerStyle.Render("📊 Results:") + "\n")
	for _, task := range m.tasks {
		if !task.Selected {
			continue
		}
		b.WriteString(fmt.Sprintf("  %s %s - %s\n", task.Icon, task.Name, task.Status))
	}

	b.WriteString("\n" + headerStyle.Render("Press Enter or q to exit • p for report") + "\n")

	return b.String()
}

func (m model) generateReport() (tea.Model, tea.Cmd) {
	m.phase = "report"
	m.reportGenerated = true
	m.addLog("INFO", "📄 Generating detailed report...")
	
	// Генерируем отчет
	report := m.createDetailedReport()
	
	// Сохраняем отчет в файл
	filename := fmt.Sprintf("ububu_report_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	err := os.WriteFile(filename, []byte(report), 0644)
	
	if err != nil {
		m.addLog("ERROR", fmt.Sprintf("Failed to save report: %v", err))
	} else {
		m.addLog("SUCCESS", fmt.Sprintf("📄 Report saved to: %s", filename))
	}
	
	return m, nil
}

func (m model) createDetailedReport() string {
	var b strings.Builder
	
	// Заголовок отчета
	b.WriteString("🐧 UBUBU 1.0 - SYSTEM OPTIMIZATION REPORT\n")
	b.WriteString(strings.Repeat("=", 60) + "\n")
	b.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	b.WriteString(fmt.Sprintf("Total Tasks: %d\n", m.totalTasks))
	b.WriteString(fmt.Sprintf("Completed: %d\n", m.completedTasks))
	b.WriteString(fmt.Sprintf("Overall Progress: %.1f%%\n\n", m.overallProgress*100))
	
	// Детальная информация по задачам
	b.WriteString("📋 TASK DETAILS\n")
	b.WriteString(strings.Repeat("-", 40) + "\n\n")
	
	for _, task := range m.tasks {
		if !task.Selected {
			continue
		}
		
		b.WriteString(fmt.Sprintf("%s %s\n", task.Icon, task.Name))
		b.WriteString(fmt.Sprintf("Description: %s\n", task.Description))
		b.WriteString(fmt.Sprintf("Status: %s\n", task.Status))
		
		if !task.StartTime.IsZero() {
			b.WriteString(fmt.Sprintf("Start Time: %s\n", task.StartTime.Format("15:04:05")))
		}
		if !task.EndTime.IsZero() {
			b.WriteString(fmt.Sprintf("End Time: %s\n", task.EndTime.Format("15:04:05")))
			duration := task.EndTime.Sub(task.StartTime)
			b.WriteString(fmt.Sprintf("Duration: %v\n", duration))
		}
		
		if task.Error != nil {
			b.WriteString(fmt.Sprintf("❌ ERROR: %v\n", task.Error))
		}
		
		if len(task.Details) > 0 {
			b.WriteString("Progress Details:\n")
			for _, detail := range task.Details {
				b.WriteString(fmt.Sprintf("  • %s\n", detail))
			}
		}
		
		b.WriteString("\n" + strings.Repeat("-", 40) + "\n\n")
	}
	
	// Логи выполнения
	if len(m.logs) > 0 {
		b.WriteString("📝 EXECUTION LOG\n")
		b.WriteString(strings.Repeat("-", 40) + "\n\n")
		for _, log := range m.logs {
			b.WriteString(log + "\n")
		}
		b.WriteString("\n")
	}
	
	// Рекомендации
	b.WriteString("💡 RECOMMENDATIONS\n")
	b.WriteString(strings.Repeat("-", 40) + "\n\n")
	
	failedTasks := 0
	for _, task := range m.tasks {
		if task.Selected && task.Error != nil {
			failedTasks++
		}
	}
	
	if failedTasks > 0 {
		b.WriteString(fmt.Sprintf("⚠️  %d task(s) failed. Review error details above.\n", failedTasks))
		b.WriteString("• Check system permissions for failed operations\n")
		b.WriteString("• Ensure sufficient disk space for cleanup operations\n")
		b.WriteString("• Verify network connectivity for update operations\n\n")
	} else {
		b.WriteString("✅ All selected tasks completed successfully!\n")
		b.WriteString("• System optimization completed without errors\n")
		b.WriteString("• Consider running optimization regularly for best performance\n\n")
	}
	
	b.WriteString("Generated by Ububu 1.0 - Ubuntu System Optimizer\n")
	b.WriteString("https://github.com/rokoss21/ububu\n")
	
	return b.String()
}

func (m model) renderReport() string {
	var b strings.Builder
	
	b.WriteString(successStyle.Render("📄 Report Generated!") + "\n\n")
	
	b.WriteString(headerStyle.Render("📊 Summary:") + "\n")
	b.WriteString(fmt.Sprintf("  Total Tasks: %d\n", m.totalTasks))
	b.WriteString(fmt.Sprintf("  Completed: %d\n", m.completedTasks))
	b.WriteString(fmt.Sprintf("  Overall Progress: %.1f%%\n\n", m.overallProgress*100))
	
	// Показываем результаты
	failedTasks := 0
	for _, task := range m.tasks {
		if task.Selected {
			if task.Error != nil {
				failedTasks++
				b.WriteString(fmt.Sprintf("  ❌ %s %s - %s\n", task.Icon, task.Name, task.Status))
				b.WriteString(fmt.Sprintf("     Error: %v\n", task.Error))
			} else {
				b.WriteString(fmt.Sprintf("  ✅ %s %s - %s\n", task.Icon, task.Name, task.Status))
			}
		}
	}
	
	if failedTasks > 0 {
		b.WriteString(fmt.Sprintf("\n⚠️  %d task(s) had errors. Check the detailed report file.\n", failedTasks))
	}
	
	b.WriteString("\n" + headerStyle.Render("Press Enter or q to exit") + "\n")
	
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
