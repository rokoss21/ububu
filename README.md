# 🐧 Ububu 1.0 - Ubuntu System Optimizer

[![Go Version](https://img.shields.io/badge/Go-1.21.5+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux-orange.svg)](https://www.linux.org)

**Interactive terminal-based system optimization tool for Ubuntu and Linux distributions**

Ububu is a modern, user-friendly system optimizer that provides real-time progress tracking, detailed reporting, and comprehensive system maintenance through an intuitive terminal interface.

## ✨ Features

- **🎯 Interactive Task Selection** - Choose which optimizations to run with keyboard navigation
- **📊 Real-time Progress Tracking** - Live progress bars with overall completion percentage
- **📋 Detailed Reporting** - Generate comprehensive reports with timing and error details
- **🖥️ Compact Design** - Optimized for standard 80x24 terminal windows
- **🏥 System Health Monitoring** - Check disk, memory, CPU, and temperature
- **🧹 Automated Cleanup** - Clean temporary files, caches, and old logs
- **🔄 System Updates** - Update packages and security patches
- **🖥️ Driver Management** - Check and update system drivers
- **⚡ Performance Optimization** - Optimize system settings and performance

## 🚀 Quick Start

### Prerequisites
- Go 1.21.5 or higher
- Ubuntu 20.04+ (or compatible Linux distribution)
- Terminal with at least 80x24 character support

### Installation

```bash
# Clone the repository
git clone https://github.com/rokoss21/ububu.git
cd ububu

# Build the application
go build ./cmd/ububu

# Run the optimizer
./ububu
```

### Alternative Build Commands
```bash
# Build all components
go build ./cmd/ububu && mv ububu ububu-cli

# Run tests
./test_all.sh

# Run with verbose output
./ububu-cli --verbose
```

## 🎮 Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate tasks |
| `Space` | Toggle task selection |
| `a` | Select all tasks |
| `n` | Select no tasks |
| `Enter` | Start optimization |
| `p` | Generate detailed report (during/after execution) |
| `q` | Quit application |

## 📋 Available Tasks

| Task | Description | Default | Duration |
|------|-------------|---------|----------|
| 🏥 **Health Check** | System health & performance monitoring | ✅ Selected | ~200ms |
| 🧹 **Cleanup** | Clean temp files, caches, and logs | ✅ Selected | ~3s |
| 🔄 **Updates** | System & security updates | ⬜ Optional | ~30s |
| 🖥️ **Drivers** | Driver updates and management | ⬜ Optional | ~15s |
| ⚡ **Optimization** | Performance optimization tweaks | ⬜ Optional | ~10s |

## 📊 Progress Tracking & Reporting

### Real-time Progress
- **Overall Progress**: Shows `X/Y tasks (Z%)` completion
- **Task Progress**: Individual task completion percentage
- **Live Updates**: Progress bars update in real-time during execution

### Detailed Reports
Press `p` during or after execution to generate comprehensive reports:
- **Session Summary**: Total duration, success/failure counts
- **Task Details**: Individual timing, errors, step-by-step progress
- **Error Diagnostics**: Detailed error information and recommendations
- **Auto-saved**: Reports saved as `ububu_report_YYYY-MM-DD_HH-MM-SS.txt`

## 🏗️ Project Structure

```
ububu/
├── cmd/ububu/              # Main application
│   └── main.go            # CLI interface with progress tracking
├── internal/
│   ├── modules/           # System optimization modules
│   │   ├── health.go     # System health checks
│   │   ├── cleanup.go    # File cleanup operations
│   │   ├── updates.go    # System updates
│   │   ├── drivers.go    # Driver management
│   │   └── optimize.go   # Performance optimization
│   ├── report/           # Report generation system
│   │   └── generator.go  # Report formatting and export
│   └── auth/             # Authentication dialogs
├── test/                 # Integration tests
├── ububu                 # Compiled executable
└── test_all.sh          # Comprehensive test suite
```

## 🔧 Development

### Build Commands
```bash
# Build main application
go build ./cmd/ububu

# Build and rename
go build ./cmd/ububu && mv ububu ububu-cli

# Run comprehensive tests
./test_all.sh

# Run specific test suites
go test ./internal/modules/ ./internal/report/
go test ./test/ -v

# Run benchmarks
go test ./internal/report/ -bench=.
```

### Dependencies
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Terminal UI framework
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling and layout
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - UI components (progress bars, spinners)

## 📊 Example Output

```
🐧 Ububu 1.0 - Ubuntu System Optimizer
==================================================

📋 Select Tasks:

> ☑ 🏥 Health Check - System health & performance
  ☑ 🧹 Cleanup - Clean temp files & caches  
  ☐ 🔄 Updates - System & security updates

Controls: ↑/↓ Navigate • Space Toggle • Enter Start • q Quit
```

## 📝 Usage Examples

### Basic Usage
```bash
# Run the optimizer
./ububu

# Select tasks with keyboard navigation
# Press Space to toggle, Enter to start
# Press 'p' during execution to generate reports
```

### Advanced Usage
```bash
# Build from source
go build ./cmd/ububu

# Run comprehensive tests
./test_all.sh

# Generate reports during execution
# Press 'p' key at any time to create detailed reports
```

### Sample Output
```
🐧 Ububu 1.0 - Ubuntu System Optimizer
==================================================

📋 Select Tasks:

> ☑ 🏥 Health Check - System health & performance
  ☑ 🧹 Cleanup - Clean temp files & caches  
  ☐ 🔄 Updates - System & security updates

Overall Progress: 1/2 tasks (50.0%)
Task: Health Check ████████████████████ 100%

Controls: ↑/↓ Navigate • Space Toggle • p Report • Enter Start • q Quit
```

## 🛠️ Technical Details

- **Language:** Go 1.21.5+
- **UI Framework:** Bubble Tea (terminal-based)
- **Architecture:** Modular system with pluggable optimization modules
- **Compatibility:** Ubuntu 20.04+ (works on most Linux distributions)
- **Dependencies:** Minimal - only terminal UI libraries
- **Performance:** Fast startup, efficient execution, timeout protection

## 🎯 Key Features

### ✅ **Core Functionality:**
- **Interactive Task Selection** - Keyboard-driven interface
- **Real-time Progress Tracking** - Live progress bars with percentages
- **Detailed Reporting** - Comprehensive reports with timing and errors
- **Timeout Protection** - No hanging on system commands
- **Compact Design** - Fits standard terminal windows (80x24)
- **Professional Interface** - Clean, debug-free output

### 🚀 **Performance Optimizations:**
- **Fast Startup** - No GUI dependencies, minimal resource usage
- **Efficient Execution** - Parallel processing where possible
- **Memory Optimized** - Minimal memory footprint
- **Reliable Operation** - Comprehensive error handling and recovery

## 🧪 Testing

The project includes comprehensive testing:

```bash
# Run all tests
./test_all.sh

# Unit tests
go test ./internal/modules/ ./internal/report/

# Integration tests  
go test ./test/ -v

# Benchmarks
go test ./internal/report/ -bench=.
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) for the excellent Bubble Tea framework
- Ubuntu community for inspiration and feedback
- Contributors and testers who helped improve the project

---

**Ububu 1.0** - Making Ubuntu optimization simple, interactive, and reliable! 🚀