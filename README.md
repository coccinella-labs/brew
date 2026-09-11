<p align="center">
  <img src="https://raw.githubusercontent.com/Coccinella-Labs/brew/main/.github/assets/thumbnail.png" alt="brew" width="100%">
</p>

# Coffee Maker Firmware

The official bare-metal firmware for IoT coffee makers with ARM Cortex-M4 microcontrollers.

[![CI](https://github.com/Coccinella-Labs/brew/actions/workflows/ci.yml/badge.svg)](https://github.com/Coccinella-Labs/brew/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Installation

### Prerequisites

```bash
# macOS
brew install arm-none-eabi-gcc

# Ubuntu/Debian
sudo apt-get install gcc-arm-none-eabi

# Arch Linux
sudo pacman -S arm-none-eabi-gcc
```

### Build

```bash
git clone https://github.com/Coccinella-Labs/brew.git
cd brew
make clean && make
```

### Docker

```bash
# Pull from GitHub Container Registry
docker pull ghcr.io/harpertoken/brew:latest

# Pull from Docker Hub
docker pull coccinella-labs/brew:latest

# Run the coffee maker dashboard
docker run -p 8080:8080 ghcr.io/harpertoken/brew:latest
```

## Usage

### Basic Example

```c
#include "heater.h"
#include "pump.h"
#include "control.h"

int main(void) {
    heater_init();
    pump_init();

    while(1) {
        coffee_state_machine();
    }

    return 0;
}
```

### Advanced Configuration

```c
// Temperature control
#define MAX_TEMP_CELSIUS    95
#define MIN_TEMP_CELSIUS    20
#define TEMP_HYSTERESIS     5

// Water level protection
#define MIN_WATER_LEVEL     50   // ml
#define CRITICAL_WATER      10   // ml

// Safety timing
#define MAX_BREW_TIME_MS    300000  // 5 minutes
```

### Hardware Setup

```c
// Pin configuration in board/board.h
#define HEATER_PIN        12
#define PUMP_PIN          13
#define TEMP_SENSOR_PIN   0
#define WATER_SENSOR_PIN  1
```

## API Reference

### Core Functions

#### `heater_init()`
Initialize heater control system.

#### `heater_on()` / `heater_off()`
Control heater state with safety checks.

#### `pump_init()`
Initialize pump control system.

#### `pump_on()` / `pump_off()`
Control pump state with flow monitoring.

#### `coffee_state_machine()`
Main brewing logic with safety interlocks.

### State Management

```c
typedef enum {
    IDLE,
    HEATING,
    BREWING,
    DONE
} coffee_state_t;
```

### Memory Layout

- **Flash**: 256KB (STM32F4)
- **RAM**: 64KB total
- **Stack**: 4KB (configurable)
- **Bootloader**: 32KB reserved

## Examples

### IoT Integration

```bash
# Start web dashboard
go run examples/web_dashboard.go

# Connect to AWS IoT
go run examples/iot_controller.go

# Monitor metrics
go run examples/metrics_collector.go
```

### Cloud Deployment

```bash
# Deploy infrastructure
aws cloudformation deploy --template-file cloud/infrastructure.yaml --stack-name coffee-maker

# Run monitoring stack
docker-compose up -d
```

## Development

### Quality Assurance

```bash
make quality         # Run all quality checks
make inspect         # Analyze memory layout
make size           # Memory usage report
```

### Testing

```bash
make flash          # Flash to hardware
make debug          # Start GDB session
```

### Pre-commit Hooks

```bash
pre-commit install  # Install quality gates
```

## Configuration

### Memory Limits

Configure in `linker.ld`:

```ld
MEMORY
{
    FLASH (rx) : ORIGIN = 0x08000000, LENGTH = 256K
    RAM (rwx)  : ORIGIN = 0x20000000, LENGTH = 64K
}
```

### Build Options

```bash
# Debug build
make CFLAGS="-mcpu=cortex-m4 -mthumb -g3 -DDEBUG"

# Release build
make CFLAGS="-mcpu=cortex-m4 -mthumb -O2 -DNDEBUG"
```

## Safety & Security

This firmware implements production-grade safety practices:

- **Memory safety**: No dynamic allocation, bounds checking
- **Static analysis**: cppcheck, MISRA C compliance
- **Stack protection**: Overflow detection and limits
- **Fail-safe design**: Safe defaults on all error conditions

### Critical Constraints

⚠️ **No dynamic memory allocation** - Uses static allocation only
⚠️ **No standard library** - Bare metal implementation
⚠️ **Real-time constraints** - Interrupt-driven architecture

## Contributing

We welcome contributions! Please see our [contributing guidelines](CONTRIBUTING.md).

### Development Setup

```bash
./setup_dev.sh      # Install dependencies
pre-commit install  # Setup quality gates
```

### Code Standards

- Follow MISRA C:2012 guidelines
- All code must pass static analysis
- Hardware changes require testing
- Documentation updates required

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/Coccinella-Labs/brew/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Coccinella-Labs/brew/discussions)

---

**Copyright (c) 2026 coccinella-labs**
