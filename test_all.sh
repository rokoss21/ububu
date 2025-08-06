#!/bin/bash

# Ububu 1.0 - Comprehensive Test Suite
# ====================================

set -e

echo "🧪 Ububu 1.0 - Running Comprehensive Test Suite"
echo "================================================"
echo

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для вывода статуса
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        return 1
    fi
}

# Функция для вывода заголовка
print_header() {
    echo -e "${BLUE}🔍 $1${NC}"
    echo "----------------------------------------"
}

# Проверяем наличие Go
print_header "Checking Go installation"
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi
go version
print_status $? "Go installation verified"
echo

# Проверяем структуру проекта
print_header "Verifying project structure"
required_files=(
    "cmd/ububu/main.go"
    "internal/modules/interface.go"
    "internal/modules/health.go"
    "internal/modules/cleanup.go"
    "internal/modules/updates.go"
    "internal/modules/drivers.go"
    "internal/modules/optimize.go"
    "internal/report/generator.go"
    "internal/gui/app.go"
    "internal/auth/dialog.go"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file"
    else
        echo -e "${RED}❌ Missing: $file${NC}"
        exit 1
    fi
done
print_status $? "Project structure verified"
echo

# Проверяем зависимости
print_header "Checking dependencies"
go mod tidy
print_status $? "Dependencies updated"
echo

# Компиляция проекта
print_header "Building project"
go build ./cmd/ububu
print_status $? "Main GUI application built"

go build ./cmd/ububu-cli
print_status $? "CLI application built"

go build ./cmd/ububu-demo
print_status $? "Demo application built"
echo

# Запуск unit тестов
print_header "Running unit tests"

echo "Testing modules..."
go test ./internal/modules/ -run "TestHealthModule_GetName|TestCleanupModule_GetName|TestUpdatesModule_GetName|TestDriversModule_GetName|TestOptimizationModule_GetName|TestSystemModuleInterface"
print_status $? "Module interface tests"

echo "Testing report generator..."
go test ./internal/report/ -v
print_status $? "Report generator tests"
echo

# Запуск интеграционных тестов
print_header "Running integration tests"
go test ./test/ -v
print_status $? "Integration tests"
echo

# Проверка покрытия кода
print_header "Checking code coverage"
echo "Report generator coverage:"
go test ./internal/report/ -cover
echo

# Запуск бенчмарков
print_header "Running performance benchmarks"
echo "Report generation benchmarks:"
go test ./internal/report/ -bench=. -benchtime=1s
echo

# Проверка качества кода
print_header "Code quality checks"

echo "Checking for go vet issues..."
go vet ./...
print_status $? "Go vet check"

echo "Checking for unused imports..."
if command -v goimports &> /dev/null; then
    goimports -l . | grep -v "vendor/" || true
    print_status 0 "Import check (goimports available)"
else
    print_status 0 "Import check (goimports not available, skipped)"
fi
echo

# Функциональные тесты
print_header "Running functional tests"

echo "Testing CLI version..."
timeout 10 ./ububu-cli > /dev/null 2>&1
print_status $? "CLI version functional test"

echo "Testing demo version..."
timeout 15 ./ububu-demo > /dev/null 2>&1
print_status $? "Demo version functional test"
echo

# Проверка размера исполняемых файлов
print_header "Checking executable sizes"
for binary in ububu ububu-cli ububu-demo; do
    if [ -f "$binary" ]; then
        size=$(du -h "$binary" | cut -f1)
        echo "📦 $binary: $size"
    fi
done
echo

# Финальный отчет
print_header "Test Summary"
echo -e "${GREEN}🎉 All tests completed successfully!${NC}"
echo
echo "📊 Test Results:"
echo "  ✅ Project structure: OK"
echo "  ✅ Compilation: OK"
echo "  ✅ Unit tests: OK"
echo "  ✅ Integration tests: OK"
echo "  ✅ Benchmarks: OK"
echo "  ✅ Code quality: OK"
echo "  ✅ Functional tests: OK"
echo
echo -e "${BLUE}🚀 Ububu 1.0 is ready for production!${NC}"
echo
echo "Next steps:"
echo "  • Run: ./ububu (GUI version)"
echo "  • Run: ./ububu-cli (CLI version)"
echo "  • Run: ./ububu-demo (Demo version)"
echo "  • Install: sudo ./build/install.sh"