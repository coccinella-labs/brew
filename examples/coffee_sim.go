package main

import (
	"fmt"
	"time"
)

// Hardware simulation
type CoffeeMaker struct {
	heaterOn    bool
	pumpOn      bool
	temperature int
	waterLevel  int
	state       string
}

func NewCoffeeMaker() *CoffeeMaker {
	return &CoffeeMaker{
		temperature: 20, // Room temp
		waterLevel:  200, // Full tank
		state:      "IDLE",
	}
}

func (c *CoffeeMaker) tick() {
	// Simulate physics
	if c.heaterOn && c.temperature < 100 {
		c.temperature += 2
	} else if !c.heaterOn && c.temperature > 20 {
		c.temperature--
	}

	if c.pumpOn && c.waterLevel > 0 {
		c.waterLevel -= 5
	}

	// State machine (mirrors firmware)
	switch c.state {
	case "IDLE":
		if c.waterLevel > 100 {
			c.heaterOn = true
			c.state = "HEATING"
		}

	case "HEATING":
		if c.temperature >= 85 {
			c.pumpOn = true
			c.state = "BREWING"
		}

	case "BREWING":
		if c.waterLevel < 50 {
			c.heaterOn = false
			c.pumpOn = false
			c.state = "DONE"
		}

	case "DONE":
		// Wait for refill
		if c.waterLevel > 150 {
			c.state = "IDLE"
		}
	}
}

func (c *CoffeeMaker) status() {
	fmt.Printf("State: %s | Temp: %d°C | Water: %d | Heater: %t | Pump: %t\n",
		c.state, c.temperature, c.waterLevel, c.heaterOn, c.pumpOn)
}

func main() {
	coffee := NewCoffeeMaker()

	for i := 0; i < 100; i++ {
		coffee.tick()
		coffee.status()
		time.Sleep(200 * time.Millisecond)

		// Simulate refill at step 80
		if i == 80 {
			coffee.waterLevel = 200
			fmt.Println("--- REFILLED ---")
		}
	}
}
