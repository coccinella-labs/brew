# GDB script for debugging
target extended-remote :3333

# Load symbols
file coffee_fw.elf

# Set breakpoints
break main
break coffee_task
break systick_handler

# Flash and run
load
monitor reset halt
continue
