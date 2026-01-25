// Serial communication between firmware and Go controller
#include <stdint.h>

// UART registers
#define UART_DATA   *((volatile uint32_t*)0x40001000)
#define UART_STATUS *((volatile uint32_t*)0x40001004)
#define UART_CTRL   *((volatile uint32_t*)0x40001008)

// Coffee maker with serial interface
typedef struct {
    uint8_t state;
    uint8_t temp;
    uint8_t water_level;
    uint8_t heater_on;
    uint8_t pump_on;
} coffee_status_t;

void uart_send_byte(uint8_t data) {
    while (!(UART_STATUS & 0x01)); // Wait for TX ready
    UART_DATA = data;
}

void send_status(coffee_status_t *status) {
    uart_send_byte(0xAA); // Start marker
    uart_send_byte(status->state);
    uart_send_byte(status->temp);
    uart_send_byte(status->water_level);
    uart_send_byte(status->heater_on);
    uart_send_byte(status->pump_on);
    uart_send_byte(0x55); // End marker
}

uint8_t uart_receive_byte(void) {
    while (!(UART_STATUS & 0x02)); // Wait for RX ready
    return UART_DATA & 0xFF;
}

// Command protocol: 'H' = heater on/off, 'P' = pump on/off
void process_commands(coffee_status_t *status) {
    if (UART_STATUS & 0x02) { // Data available
        uint8_t cmd = uart_receive_byte();
        switch(cmd) {
            case 'H': status->heater_on = !status->heater_on; break;
            case 'P': status->pump_on = !status->pump_on; break;
        }
    }
}
