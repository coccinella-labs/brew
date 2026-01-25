target extended-remote :3333
file build/coffee-fw.elf
load
monitor reset halt
break main
continue
