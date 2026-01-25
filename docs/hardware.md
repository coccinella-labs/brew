# Hardware Setup

## Required Components

- **STM32F4 Discovery Board** (STM32F407VGT6)
- **ST-Link V2 Debugger** (usually integrated)
- **Coffee maker interface board** (custom PCB)
- **Temperature sensor** (DS18B20 or similar)
- **Water level sensor** (capacitive or ultrasonic)
- **Relay modules** for heater and pump control

## Pin Connections

| Function | STM32 Pin | Description |
|----------|-----------|-------------|
| Heater Control | PA12 | GPIO output to relay |
| Pump Control | PA13 | GPIO output to relay |
| Temperature Sensor | PA0 | ADC input |
| Water Level Sensor | PA1 | ADC input |
| Status LED | PD12 | GPIO output |
| Debug UART | PA2/PA3 | USART2 TX/RX |

## Power Requirements

- **5V supply** for relay modules
- **3.3V supply** for STM32 (from ST-Link or external)
- **12V supply** for coffee maker heater/pump

## Safety Considerations

⚠️ **High voltage warning**: Coffee maker operates at mains voltage
⚠️ **Isolation required**: Use optocouplers for mains control
⚠️ **Thermal protection**: Implement hardware thermal cutoff
⚠️ **Water protection**: Ensure electronics are waterproof

## Schematic

```
STM32F4 -----> Optocoupler -----> Relay -----> Coffee Maker
   |              |                  |
   |              +-- 5V Supply      +-- Mains Voltage
   |
   +-- Temperature/Water Sensors (3.3V)
```
