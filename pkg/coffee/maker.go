package coffee

// CoffeeMaker represents the coffee maker hardware simulation
type CoffeeMaker struct {
	HeaterOn    bool
	PumpOn      bool
	Temperature int
	WaterLevel  int
	State       string
}

// NewCoffeeMaker creates a new coffee maker instance
func NewCoffeeMaker() *CoffeeMaker {
	return &CoffeeMaker{
		Temperature: 20,  // Room temp
		WaterLevel:  200, // Full tank
		State:       "IDLE",
	}
}

// Tick advances the coffee maker simulation by one step
func (c *CoffeeMaker) Tick() {
	// Simulate physics
	if c.HeaterOn && c.Temperature < 100 {
		c.Temperature += 2
	} else if !c.HeaterOn && c.Temperature > 20 {
		c.Temperature--
	}

	if c.PumpOn && c.WaterLevel > 0 {
		c.WaterLevel -= 5
	}

	// State machine (mirrors firmware)
	switch c.State {
	case "IDLE":
		if c.WaterLevel > 100 {
			c.HeaterOn = true
			c.State = "HEATING"
		}

	case "HEATING":
		if c.Temperature >= 85 {
			c.PumpOn = true
			c.State = "BREWING"
		}

	case "BREWING":
		if c.WaterLevel < 50 {
			c.HeaterOn = false
			c.PumpOn = false
			c.State = "DONE"
		}

	case "DONE":
		// Wait for refill
		if c.WaterLevel > 150 {
			c.State = "IDLE"
		}
	}
}
