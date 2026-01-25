# Production bare-metal Makefile for coffee maker firmware
TARGET = coffee-fw
MCU = cortex-m4

# Toolchain
CC = arm-none-eabi-gcc
AS = arm-none-eabi-as
LD = arm-none-eabi-ld
OBJCOPY = arm-none-eabi-objcopy
OBJDUMP = arm-none-eabi-objdump
SIZE = arm-none-eabi-size
GDB = arm-none-eabi-gdb

# Directories
BUILD_DIR = build
SRC_DIR = src
STARTUP_DIR = startup
DRIVERS_DIR = drivers
INCLUDE_DIR = include

# Sources
STARTUP_SRC = $(STARTUP_DIR)/startup.s $(STARTUP_DIR)/system_init.c
APP_SRC = $(wildcard $(SRC_DIR)/*.c)
DRIVER_SRC = $(wildcard $(DRIVERS_DIR)/*.c)
ALL_SRC = $(STARTUP_SRC) $(APP_SRC) $(DRIVER_SRC)

# Objects
OBJS = $(ALL_SRC:%.c=$(BUILD_DIR)/%.o)
OBJS := $(OBJS:%.s=$(BUILD_DIR)/%.o)

# Flags
CFLAGS = -mcpu=$(MCU) -mthumb -mfloat-abi=soft
CFLAGS += -Wall -Wextra -Werror -O2 -g3
CFLAGS += -ffunction-sections -fdata-sections
CFLAGS += -I$(INCLUDE_DIR) -I$(STARTUP_DIR) -Iboard/

LDFLAGS = -T linker.ld -Wl,--gc-sections -Wl,--print-memory-usage -nostdlib

# Memory limits (bytes)
MAX_FLASH = 262144  # 256KB
MAX_RAM = 65536     # 64KB

.PHONY: all clean flash debug size inspect quality

all: $(BUILD_DIR)/$(TARGET).elf $(BUILD_DIR)/$(TARGET).bin

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)/$(SRC_DIR) $(BUILD_DIR)/$(STARTUP_DIR) $(BUILD_DIR)/$(DRIVERS_DIR)

# Compile C files
$(BUILD_DIR)/%.o: %.c | $(BUILD_DIR)
	$(CC) $(CFLAGS) -c $< -o $@

# Assemble files
$(BUILD_DIR)/%.o: %.s | $(BUILD_DIR)
	$(AS) -mcpu=$(MCU) -mthumb $< -o $@

# Link ELF
$(BUILD_DIR)/$(TARGET).elf: $(OBJS)
	$(CC) $(CFLAGS) $(LDFLAGS) $^ -o $@
	@echo "=== Memory Usage ==="
	$(SIZE) $@
	@echo "=== Checking limits ==="
	@FLASH=$$($(SIZE) $@ | awk 'NR==2 {print $$1}'); \
	 RAM=$$($(SIZE) $@ | awk 'NR==2 {print $$2}'); \
	 echo "Flash: $$FLASH/$(MAX_FLASH) bytes"; \
	 echo "RAM: $$RAM/$(MAX_RAM) bytes"; \
	 test $$FLASH -le $(MAX_FLASH) || (echo "FLASH OVERFLOW!" && exit 1); \
	 test $$RAM -le $(MAX_RAM) || (echo "RAM OVERFLOW!" && exit 1)

# Create binary
$(BUILD_DIR)/$(TARGET).bin: $(BUILD_DIR)/$(TARGET).elf
	$(OBJCOPY) -O binary $< $@

# Quality checks
quality:
	./scripts/quality_check.sh

# Flash to hardware
flash: $(BUILD_DIR)/$(TARGET).elf
	openocd -f scripts/openocd.cfg \
	        -c "program $< verify reset exit"

# Debug session
debug: $(BUILD_DIR)/$(TARGET).elf
	openocd -f scripts/openocd.cfg &
	sleep 2
	$(GDB) -x scripts/debug.gdb $<
	pkill openocd

# Memory inspection
inspect: $(BUILD_DIR)/$(TARGET).elf
	@echo "=== Sections ==="
	$(OBJDUMP) -h $<
	@echo "=== Reset Handler ==="
	$(OBJDUMP) -d $< | grep -A 10 "Reset_Handler"
	@echo "=== Symbols ==="
	arm-none-eabi-nm $< | grep -E "(Reset_Handler|main|_stack)"

# Size analysis
size: $(BUILD_DIR)/$(TARGET).elf
	$(SIZE) -A $<
	$(SIZE) -B $<

clean:
	rm -rf $(BUILD_DIR)

# CI target
ci: all quality inspect
	@echo "✅ Firmware build successful"
