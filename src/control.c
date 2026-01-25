#include "control.h"
#include "board.h"
#include "heater.h"
#include "pump.h"

static coffee_state_t state = IDLE;
static uint32_t target_temp = 85;

void coffee_state_machine(void) {
    uint32_t temp = read_temperature();
    uint32_t water = read_water_level();

    switch(state) {
        case IDLE:
            if (water > 100) {
                heater_on();
                state = HEATING;
            }
            break;

        case HEATING:
            if (temp >= target_temp) {
                pump_on();
                state = BREWING;
            }
            break;

        case BREWING:
            if (water < 50) {
                pump_off();
                heater_off();
                state = DONE;
            }
            break;

        case DONE:
            if (water > 150) state = IDLE;
            break;
    }
}
