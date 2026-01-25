// Reset handler and vector table
.syntax unified
.cpu cortex-m4
.thumb

.section .vector_table,"a",%progbits
.type vector_table, %object
.size vector_table, .-vector_table

vector_table:
    .word _stack_top
    .word Reset_Handler
    .word NMI_Handler
    .word HardFault_Handler

.section .text.Reset_Handler
.weak Reset_Handler
.type Reset_Handler, %function
Reset_Handler:
    // Set stack pointer
    ldr sp, =_stack_top

    // Copy .data from flash to RAM
    ldr r0, =_sdata
    ldr r1, =_edata
    ldr r2, =_sidata
    bl copy_data

    // Zero .bss
    ldr r0, =_sbss
    ldr r1, =_ebss
    bl zero_bss

    // Jump to main
    bl main

    // Infinite loop if main returns
hang:
    b hang

copy_data:
    cmp r0, r1
    bge copy_done
    ldr r3, [r2], #4
    str r3, [r0], #4
    b copy_data
copy_done:
    bx lr

zero_bss:
    mov r2, #0
zero_loop:
    cmp r0, r1
    bge zero_done
    str r2, [r0], #4
    b zero_loop
zero_done:
    bx lr

.weak NMI_Handler
.weak HardFault_Handler
NMI_Handler:
HardFault_Handler:
    b .
