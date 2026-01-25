// Enhanced startup with interrupt vectors
#include <stdint.h>

extern int main(void);
extern void systick_handler(void);
extern void temp_sensor_handler(void);

#define STACK_TOP 0x20010000

void reset_handler(void) {
    main();
}

// Default handler for unused interrupts
void default_handler(void) {
    while(1); // Trap
}

// Complete vector table
__attribute__((section(".vector_table")))
const uint32_t vector_table[] = {
    STACK_TOP,                    // 0: Initial stack pointer
    (uint32_t)reset_handler,      // 1: Reset
    (uint32_t)default_handler,    // 2: NMI
    (uint32_t)default_handler,    // 3: Hard Fault
    // ... more system exceptions ...
    0, 0, 0, 0, 0, 0, 0,         // 4-10: Reserved
    (uint32_t)default_handler,    // 11: SVCall
    0, 0,                        // 12-13: Reserved
    (uint32_t)default_handler,    // 14: PendSV
    (uint32_t)systick_handler,    // 15: SysTick

    // External interrupts start here
    (uint32_t)temp_sensor_handler, // 16: Temperature sensor
    // Add more as needed...
};
