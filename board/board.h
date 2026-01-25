#ifndef BOARD_H
#define BOARD_H

// Basic types for bare metal
typedef unsigned int uint32_t;
typedef unsigned char uint8_t;

// Pin definitions
#define HEATER_PIN        12
#define PUMP_PIN          13
#define TEMP_SENSOR_PIN   0
#define WATER_SENSOR_PIN  1

// GPIO modes
#define GPIO_OUTPUT       1
#define GPIO_INPUT        0

// Function prototypes
void gpio_init(uint32_t pin, uint32_t mode);
void gpio_set(uint32_t pin);
void gpio_clear(uint32_t pin);
uint32_t adc_read(uint32_t channel);

#endif
