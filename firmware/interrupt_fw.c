// Interrupt-driven coffee maker firmware
#include <stdint.h>

// System timer for 1ms ticks
#define SYSTICK_CTRL  *((volatile uint32_t*)0xE000E010)
#define SYSTICK_LOAD  *((volatile uint32_t*)0xE000E014)
#define SYSTICK_VAL   *((volatile uint32_t*)0xE000E018)

// Global state
volatile uint32_t system_ticks = 0;
volatile uint8_t brew_timer = 0;

// SysTick interrupt handler - called every 1ms
void systick_handler(void) {
    system_ticks++;

    if (brew_timer > 0) {
        brew_timer--;
    }
}

// Temperature sensor interrupt - called when temp changes
void temp_sensor_handler(void) {
    uint32_t temp = TEMP_SENSOR;

    if (temp > 95) { // Overheat protection
        HEATER_CTRL = 0;
    }
}

void setup_interrupts(void) {
    // Configure SysTick for 1ms interrupts (assuming 72MHz clock)
    SYSTICK_LOAD = 72000 - 1;
    SYSTICK_CTRL = 0x07; // Enable, interrupt, use processor clock
}

// Non-blocking coffee state machine
void coffee_task(void) {
    static uint32_t last_check = 0;
    static coffee_state_t state = IDLE;

    // Run every 100ms
    if (system_ticks - last_check < 100) return;
    last_check = system_ticks;

    uint32_t temp = TEMP_SENSOR;
    uint32_t water = WATER_LEVEL;

    switch(state) {
        case IDLE:
            if (water > 100) {
                HEATER_CTRL = 1;
                state = HEATING;
            }
            break;

        case HEATING:
            if (temp >= 85) {
                PUMP_CTRL = 1;
                brew_timer = 30; // 30 second brew
                state = BREWING;
            }
            break;

        case BREWING:
            if (brew_timer == 0 || water < 50) {
                PUMP_CTRL = 0;
                HEATER_CTRL = 0;
                state = DONE;
            }
            break;

        case DONE:
            if (water > 150) state = IDLE;
            break;
    }
}

int main(void) {
    setup_interrupts();

    while(1) {
        coffee_task();
        // CPU can sleep here - interrupts will wake it
    }
}
