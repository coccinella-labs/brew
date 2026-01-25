#ifndef CONTROL_H
#define CONTROL_H

// Basic types
typedef unsigned int uint32_t;

typedef enum {
    IDLE,
    HEATING,
    BREWING,
    DONE
} coffee_state_t;

void coffee_state_machine(void);

#endif
