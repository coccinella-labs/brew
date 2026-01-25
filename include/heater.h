#ifndef HEATER_H
#define HEATER_H

// Basic types
typedef unsigned int uint32_t;

void heater_init(void);
void heater_on(void);
void heater_off(void);
uint32_t read_temperature(void);

#endif
