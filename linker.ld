/* Linker script for ARM Cortex-M4 */
MEMORY
{
    FLASH (rx) : ORIGIN = 0x08000000, LENGTH = 256K
    RAM (rwx)  : ORIGIN = 0x20000000, LENGTH = 64K
}

/* Stack size */
_stack_size = 4K;

SECTIONS
{
    .text : {
        KEEP(*(.vector_table))
        *(.text*)
    } > FLASH

    .data : {
        _sdata = .;
        *(.data*)
        _edata = .;
    } > RAM AT > FLASH

    _sidata = LOADADDR(.data);

    .bss : {
        _sbss = .;
        *(.bss*)
        _ebss = .;
    } > RAM

    /* Stack at end of RAM */
    _stack_top = ORIGIN(RAM) + LENGTH(RAM);
}
