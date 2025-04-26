# GoSyn - Go Syntax Quick Reference Tool

A command-line tool providing quick access to Go programming syntax with colorized output. Perfect for developers who need quick reminders of Go concepts without leaving the terminal.

## Features

- **Colorized Output**: Easy-to-read syntax examples with ANSI colors
- **Comprehensive Sections**: Covers Variables, Conditionals, Loops, Functions, Concurrency, and more
- **Quick Navigation**: Jump directly to specific syntax patterns
- **Alias Support**: Short commands for frequent actions (lsec, lsub)

## Installation

```bash
# Install with Go
go install github.com/bbarrington0099/gosyn@latest

# Verify installation
gosyn help
```

## Usage

### Basic Commands

```bash
# Show help
gosyn help

# List all sections
gosyn listSections
gosyn lsec

# List subsections in a section
gosyn listSubsections Variables
gosyn lsub Functions

# Get syntax for a subsection
gosyn Functions Declaration
gosyn Slices BasicOperations
```

### Command Aliases
| Full Command         | Alias | Example                   |
|----------------------|-------|---------------------------|
| `help`               | `h`   | `gosyn h`                 |
| `listSections`       | `lsec`| `gosyn lsec`              |
| `listSubsections`    | `lsub`| `gosyn lsub Concurrency`  |


## Roadmap & Contributions

We welcome contributions! Here are priority areas:

### High Priority
1. **Tab Completion**
   - Implement bash/zsh/fish completion
   - Auto-complete section/subsection names
2. **Section Improvements**
   - Add improved examples
   - Expand Error Handling patterns
   - Add additional sections

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes with clear messages
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---