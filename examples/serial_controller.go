package main

import (
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

type CoffeeStatus struct {
	State      uint8
	Temp       uint8
	WaterLevel uint8
	HeaterOn   uint8
	PumpOn     uint8
}

type SerialController struct {
	port serial.Port
}

func NewSerialController(portName string) (*SerialController, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}

	return &SerialController{port: port}, nil
}

func (sc *SerialController) ReadStatus() (*CoffeeStatus, error) {
	buf := make([]byte, 7) // AA + 5 data bytes + 55
	n, err := sc.port.Read(buf)
	if err != nil || n != 7 {
		return nil, fmt.Errorf("read error: %v", err)
	}

	if buf[0] != 0xAA || buf[6] != 0x55 {
		return nil, fmt.Errorf("invalid frame")
	}

	return &CoffeeStatus{
		State:      buf[1],
		Temp:       buf[2],
		WaterLevel: buf[3],
		HeaterOn:   buf[4],
		PumpOn:     buf[5],
	}, nil
}

func (sc *SerialController) SendCommand(cmd byte) error {
	_, err := sc.port.Write([]byte{cmd})
	return err
}

func (sc *SerialController) Close() {
	sc.port.Close()
}

// Simulation mode when no hardware available
func simulationMode() {
	fmt.Println("Running in simulation mode...")
	coffee := NewCoffeeMaker()

	for i := 0; i < 50; i++ {
		coffee.tick()
		coffee.status()
		time.Sleep(300 * time.Millisecond)
	}
}

func main() {
	// Try to connect to hardware first
	controller, err := NewSerialController("/dev/tty.usbserial") // Adjust port
	if err != nil {
		fmt.Printf("No hardware found: %v\n", err)
		simulationMode()
		return
	}
	defer controller.Close()

	fmt.Println("Connected to coffee maker hardware!")

	// Monitor and control loop
	for {
		status, err := controller.ReadStatus()
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}

		fmt.Printf("HW State: %d | Temp: %d°C | Water: %d | H:%d P:%d\n",
			status.State, status.Temp, status.WaterLevel,
			status.HeaterOn, status.PumpOn)

		// Auto-control logic
		if status.Temp < 80 && status.WaterLevel > 100 {
			controller.SendCommand('H') // Turn on heater
		}

		time.Sleep(500 * time.Millisecond)
	}
}
