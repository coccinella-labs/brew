package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"coffee-sim/pkg/coffee"
)

type WebDashboard struct {
	coffee *coffee.CoffeeMaker
	status chan coffee.CoffeeStatus
}

func NewWebDashboard() *WebDashboard {
	return &WebDashboard{
		coffee: coffee.NewCoffeeMaker(),
		status: make(chan coffee.CoffeeStatus, 10),
	}
}

func (wd *WebDashboard) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	status := coffee.CoffeeStatus{
		State:      stateToInt(wd.coffee.State),
		Temp:       uint8(wd.coffee.Temperature),
		WaterLevel: uint8(wd.coffee.WaterLevel),
		HeaterOn:   boolToInt(wd.coffee.HeaterOn),
		PumpOn:     boolToInt(wd.coffee.PumpOn),
	}

	json.NewEncoder(w).Encode(status)
}

func (wd *WebDashboard) handleControl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST only", 405)
		return
	}

	action := r.URL.Query().Get("action")
	switch action {
	case "start":
		wd.coffee.WaterLevel = 200 // Refill
	case "stop":
		wd.coffee.HeaterOn = false
		wd.coffee.PumpOn = false
	}

	fmt.Fprintf(w, "OK")
}

func (wd *WebDashboard) handleDashboard(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head><title>Coffee Maker Dashboard</title></head>
<body>
<h1>☕ Coffee Maker Control</h1>
<div id="status"></div>
<button onclick="control('start')">Start Brew</button>
<button onclick="control('stop')">Stop</button>

<script>
function updateStatus() {
    fetch('/api/status')
        .then(r => r.json())
        .then(data => {
            document.getElementById('status').innerHTML =
                'State: ' + ['IDLE','HEATING','BREWING','DONE'][data.State] +
                ' | Temp: ' + data.Temp + '°C' +
                ' | Water: ' + data.WaterLevel +
                ' | Heater: ' + (data.HeaterOn ? 'ON' : 'OFF') +
                ' | Pump: ' + (data.PumpOn ? 'ON' : 'OFF');
        });
}

function control(action) {
    fetch('/api/control?action=' + action, {method: 'POST'});
}

setInterval(updateStatus, 1000);
updateStatus();
</script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func (wd *WebDashboard) run() {
	// Simulate coffee maker
	go func() {
		for {
			wd.coffee.Tick()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Web server
	http.HandleFunc("/", wd.handleDashboard)
	http.HandleFunc("/api/status", wd.handleStatus)
	http.HandleFunc("/api/control", wd.handleControl)

	fmt.Println("Dashboard: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func stateToInt(state string) uint8 {
	states := map[string]uint8{"IDLE": 0, "HEATING": 1, "BREWING": 2, "DONE": 3}
	return states[state]
}

func boolToInt(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func main() {
	dashboard := NewWebDashboard()
	dashboard.run()
}
