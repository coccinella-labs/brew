// Minimal coffee maker firmware - bare metal
#include <stdint.h>

// Hardware registers (example addresses)
#define HEATER_CTRL     *((volatile uint32_t*)0x40000000)
#define PUMP_CTRL       *((volatile uint32_t*)0x40000004)
#define TEMP_SENSOR     *((volatile uint32_t*)0x40000008)
#define WATER_LEVEL     *((volatile uint32_t*)0x4000000C)

// Coffee maker states
typedef enum {
    IDLE,
    HEATING,
    BREWING,
    DONE
} coffee_state_t;

static coffee_state_t state = IDLE;
static uint32_t target_temp = 85; // Celsius

void coffee_state_machine(void) {
    uint32_t current_temp = TEMP_SENSOR;
    uint32_t water = WATER_LEVEL;

    switch(state) {
        case IDLE:
            if (water > 100) { // Enough water
                HEATER_CTRL = 1;
                state = HEATING;
            }
            break;

        case HEATING:
            if (current_temp >= target_temp) {
                PUMP_CTRL = 1;
                state = BREWING;
            }
            break;

        case BREWING:
            if (water < 50) { // Low water
                PUMP_CTRL = 0;
                HEATER_CTRL = 0;
                state = DONE;
            }
            break;

        case DONE:
            // Wait for reset
            break;
    }
}

int main(void) {
    while(1) {
        coffee_state_machine();
        // Add delay here
    }
}
