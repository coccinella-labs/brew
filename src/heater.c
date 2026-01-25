#include "heater.h"
#include "board.h"

void heater_init(void) {
    gpio_init(HEATER_PIN, GPIO_OUTPUT);
}

void heater_on(void) {
    gpio_set(HEATER_PIN);
}

void heater_off(void) {
    gpio_clear(HEATER_PIN);
}

uint32_t read_temperature(void) {
    return adc_read(TEMP_SENSOR_PIN);
}
