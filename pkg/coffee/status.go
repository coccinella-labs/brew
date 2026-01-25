package coffee

// CoffeeStatus represents the current status of the coffee maker
type CoffeeStatus struct {
	State      uint8
	Temp       uint8
	WaterLevel uint8
	HeaterOn   uint8
	PumpOn     uint8
}
