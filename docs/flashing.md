# Flashing Guide

## Prerequisites

- ARM GCC toolchain installed
- OpenOCD installed
- ST-Link drivers installed
- Hardware connected properly

## Build Firmware

```bash
# Clean previous builds
make clean

# Build with safety checks
make

# Verify memory usage
make inspect
```

## Flash Methods

### Method 1: OpenOCD (Recommended)

```bash
# Flash using OpenOCD
make flash

# Manual OpenOCD command
openocd -f scripts/openocd.cfg \
        -c "program build/coffee-fw.elf verify reset exit"
```

### Method 2: ST-Link Utility

```bash
# Convert to hex format
arm-none-eabi-objcopy -O ihex build/coffee-fw.elf build/coffee-fw.hex

# Flash with ST-Link utility
st-flash --format ihex write build/coffee-fw.hex
```

### Method 3: DFU (Device Firmware Update)

```bash
# Put device in DFU mode (BOOT0=1, reset)
# Flash via DFU
dfu-util -a 0 -s 0x08000000 -D build/coffee-fw.bin
```

## Verification

After flashing, verify the firmware:

```bash
# Check reset vector
arm-none-eabi-objdump -d build/coffee-fw.elf | head -20

# Verify flash contents
st-flash read flash_dump.bin 0x08000000 0x40000
hexdump -C flash_dump.bin | head
```

## Troubleshooting

**Flash fails**: Check ST-Link connection and drivers
**Verification error**: Ensure proper power supply
**Device not found**: Try different USB port or cable
**Permission denied**: Run with sudo (Linux) or as Administrator (Windows)
