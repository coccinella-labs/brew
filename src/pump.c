#include "pump.h"
#include "board.h"

void pump_init(void) {
    gpio_init(PUMP_PIN, GPIO_OUTPUT);
}

void pump_on(void) {
    gpio_set(PUMP_PIN);
}

void pump_off(void) {
    gpio_clear(PUMP_PIN);
}

uint32_t read_water_level(void) {
    return adc_read(WATER_SENSOR_PIN);
}
