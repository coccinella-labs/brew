#ifndef PUMP_H
#define PUMP_H

// Basic types
typedef unsigned int uint32_t;

void pump_init(void);
void pump_on(void);
void pump_off(void);
uint32_t read_water_level(void);

#endif
