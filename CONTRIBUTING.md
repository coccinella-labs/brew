# Contributing to Coffee Maker Firmware

We welcome contributions to the coffee maker firmware project! This document outlines how to contribute effectively.

## Getting Started

### Prerequisites

```bash
# Install ARM toolchain
brew install arm-none-eabi-gcc  # macOS
sudo apt-get install gcc-arm-none-eabi  # Ubuntu

# Install development tools
pip install pre-commit
```

### Setup

```bash
git clone https://github.com/harpertoken/brew.git
cd brew
./setup_dev.sh
pre-commit install
```

## Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Follow MISRA C:2012 guidelines
- Write safety-critical code only
- No dynamic memory allocation
- Document all functions

### 3. Test Changes

```bash
make clean && make          # Build firmware
make quality               # Run quality checks
./scripts/quality_simple.sh  # Additional checks
```

### 4. Commit Changes

```bash
git add .
git commit -m "feat: add temperature safety check"
```

Pre-commit hooks will automatically run quality checks.

## Code Standards

### Safety Requirements

- **No dynamic allocation** - Use static memory only
- **Bounds checking** - Validate all array access
- **Fail-safe defaults** - Safe state on errors
- **Interrupt safety** - Atomic operations where needed

### Code Style

```c
// Function naming: snake_case
void heater_init(void);

// Constants: UPPER_CASE
#define MAX_TEMP_CELSIUS    95

// Variables: snake_case
static uint32_t current_temp;

// Always check return values
if (heater_init() != SUCCESS) {
    // Handle error
}
```

### Documentation

```c
/**
 * @brief Initialize heater control system
 * @return SUCCESS on success, ERROR_CODE on failure
 * @safety Critical - must be called before heater_on()
 */
heater_result_t heater_init(void);
```

## Testing

### Hardware Testing

```bash
make flash    # Flash to development board
make debug    # Start GDB debugging session
```

### Static Analysis

All code must pass:
- cppcheck (no warnings)
- MISRA C compliance
- Stack usage analysis
- Duplicate code detection

## Pull Request Process

### 1. Quality Gates

- [ ] All tests pass
- [ ] Static analysis clean
- [ ] Memory usage within limits
- [ ] Documentation updated

### 2. Review Checklist

- [ ] Safety-critical code reviewed
- [ ] No security vulnerabilities
- [ ] Performance impact assessed
- [ ] Hardware compatibility verified

### 3. Merge Requirements

- Minimum 2 approvals required
- All CI checks must pass
- No merge conflicts
- Squash commits before merge

## Issue Guidelines

### Bug Reports

```markdown
**Hardware**: STM32F4, Rev 2.1
**Firmware Version**: v1.2.3
**Steps to Reproduce**:
1. Power on device
2. Set temperature to 95°C
3. Observe behavior

**Expected**: Heater activates
**Actual**: No response
**Logs**: [attach debug output]
```

### Feature Requests

```markdown
**Feature**: Automatic descaling cycle
**Justification**: Maintenance requirement
**Safety Impact**: Low
**Implementation**: [brief description]
```

## Architecture Guidelines

### Memory Layout

```
Flash (256KB):
├── Bootloader (32KB)
├── Application (200KB)
└── Configuration (24KB)

RAM (64KB):
├── Stack (4KB)
├── Heap (0KB - prohibited)
└── Static data (60KB)
```

### Module Structure

```
src/
├── control.c     # State machine
├── heater.c      # Temperature control
├── pump.c        # Water flow control
└── safety.c      # Safety interlocks
```

## Security Guidelines

### Code Review Focus

- Buffer overflow prevention
- Integer overflow checks
- Input validation
- Privilege escalation prevention

### Secure Coding

```c
// Always validate inputs
if (temp > MAX_TEMP_CELSIUS || temp < MIN_TEMP_CELSIUS) {
    return ERROR_INVALID_TEMP;
}

// Use safe string functions
strncpy(buffer, input, sizeof(buffer) - 1);
buffer[sizeof(buffer) - 1] = '\0';
```

## Release Process

### Version Numbering

- `MAJOR.MINOR.PATCH` (semantic versioning)
- Major: Breaking changes
- Minor: New features
- Patch: Bug fixes

### Release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] Security review complete
- [ ] Hardware validation done
- [ ] Release notes prepared

## Community

### Communication

- **Issues**: Bug reports and feature requests
- **Discussions**: Design questions and help
- **Security**: security@harpertoken.com

### Code of Conduct

- Be respectful and professional
- Focus on technical merit
- Welcome newcomers
- Maintain safety-first mindset

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Questions?** Open a discussion or contact the maintainers.
