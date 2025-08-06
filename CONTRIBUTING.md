# Contributing to Ububu

Thank you for your interest in contributing to Ububu! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites
- Go 1.21.5 or higher
- Ubuntu 20.04+ (or compatible Linux distribution)
- Git for version control

### Development Setup
```bash
# Clone the repository
git clone https://github.com/rokoss21/ububu.git
cd ububu

# Build the project
go build ./cmd/ububu

# Run tests
./test_all.sh
```

## 🏗️ Project Structure

```
ububu/
├── cmd/ububu/              # Main application
├── internal/
│   ├── modules/           # System optimization modules
│   ├── report/           # Report generation
│   └── auth/             # Authentication dialogs
├── test/                 # Integration tests
└── docs/                 # Documentation
```

## 🔧 Development Guidelines

### Code Style
- Follow standard Go conventions
- Use `gofmt` for formatting
- Write clear, descriptive variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and small

### Testing
- Write unit tests for new functionality
- Ensure all tests pass before submitting PR
- Include integration tests for system modules
- Test on multiple Ubuntu versions when possible

### Commit Messages
- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, Remove)
- Keep first line under 50 characters
- Include detailed description if needed

Example:
```
Add progress tracking to cleanup module

- Implement real-time progress callbacks
- Add error handling for disk space checks
- Update tests to cover new functionality
```

## 🐛 Bug Reports

When reporting bugs, please include:
- Ubuntu version and system specifications
- Steps to reproduce the issue
- Expected vs actual behavior
- Error messages or logs
- Screenshots if applicable

## ✨ Feature Requests

For new features:
- Describe the problem you're trying to solve
- Explain your proposed solution
- Consider backward compatibility
- Discuss performance implications

## 🔄 Pull Request Process

1. **Fork the repository** and create a feature branch
2. **Make your changes** following the coding guidelines
3. **Add tests** for new functionality
4. **Run the test suite** and ensure all tests pass
5. **Update documentation** if needed
6. **Submit a pull request** with a clear description

### PR Checklist
- [ ] Code follows project style guidelines
- [ ] All tests pass (`./test_all.sh`)
- [ ] New functionality includes tests
- [ ] Documentation updated if needed
- [ ] Commit messages are clear and descriptive
- [ ] No breaking changes (or clearly documented)

## 🧪 Testing

### Running Tests
```bash
# All tests
./test_all.sh

# Unit tests only
go test ./internal/modules/ ./internal/report/

# Integration tests
go test ./test/ -v

# Benchmarks
go test ./internal/report/ -bench=.

# Coverage
go test -cover ./internal/modules/
```

### Writing Tests
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Test both success and error cases
- Include edge cases and boundary conditions

## 📝 Documentation

- Update README.md for user-facing changes
- Add inline comments for complex code
- Update AGENTS.md for development changes
- Include examples in documentation

## 🏷️ Release Process

1. Update version numbers
2. Update CHANGELOG.md
3. Create release notes
4. Tag the release
5. Build and test release binaries

## 🤝 Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Maintain a professional tone

## 📞 Getting Help

- Open an issue for bugs or feature requests
- Check existing issues before creating new ones
- Join discussions in pull requests
- Ask questions in issue comments

## 🎯 Areas for Contribution

We welcome contributions in these areas:
- **New system modules** (monitoring, optimization)
- **Performance improvements** (speed, memory usage)
- **UI/UX enhancements** (better progress display, colors)
- **Testing** (more test coverage, edge cases)
- **Documentation** (examples, tutorials, guides)
- **Platform support** (other Linux distributions)

Thank you for contributing to Ububu! 🚀