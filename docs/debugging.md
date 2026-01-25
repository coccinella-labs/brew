# Debugging Tips

## Debug Setup

### GDB with OpenOCD

```bash
# Start debug session
make debug

# Manual setup
openocd -f scripts/openocd.cfg &
arm-none-eabi-gdb build/coffee-fw.elf
(gdb) target remote :3333
(gdb) monitor reset halt
```

## Common Debug Commands

```bash
# Load symbols and reset
(gdb) file build/coffee-fw.elf
(gdb) load
(gdb) monitor reset halt

# Set breakpoints
(gdb) break main
(gdb) break Reset_Handler
(gdb) break coffee_state_machine

# Examine memory
(gdb) x/10x 0x20000000    # RAM start
(gdb) x/10x 0x08000000    # Flash start

# Check stack
(gdb) info registers
(gdb) bt                  # Backtrace
```

## Debugging Startup Issues

### Reset Handler Problems

```bash
# Check vector table
(gdb) x/8x 0x08000000
# Should show: stack_top, reset_handler, nmi_handler, hardfault_handler

# Step through startup
(gdb) break Reset_Handler
(gdb) continue
(gdb) stepi               # Step instruction by instruction
```

### Stack Issues

```bash
# Check stack pointer
(gdb) info registers sp
(gdb) x/20x $sp          # Examine stack contents

# Check for stack overflow
(gdb) print _stack_top
(gdb) print $sp
```

## Runtime Debugging

### State Machine Issues

```bash
# Monitor coffee state
(gdb) print state
(gdb) print target_temp
(gdb) watch state        # Break when state changes
```

### Hardware Issues

```bash
# Check GPIO registers
(gdb) x/x 0x40020000     # GPIOA base
(gdb) x/x 0x40020014     # GPIOA ODR (output data)

# Monitor ADC readings
(gdb) print temp
(gdb) print water
```

## Common Issues

**HardFault**: Check for null pointer dereference or stack overflow
**Infinite loop**: Use Ctrl+C to break and check location
**No output**: Verify UART pins and baud rate
**Erratic behavior**: Check for uninitialized variables
