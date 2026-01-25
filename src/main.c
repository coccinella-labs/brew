#include "heater.h"
#include "pump.h"
#include "control.h"

int main(void) {
    heater_init();
    pump_init();

    while(1) {
        coffee_state_machine();
    }

    return 0;
}
