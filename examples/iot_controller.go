package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type IoTController struct {
	client mqtt.Client
	coffee *CoffeeMaker
}

func NewIoTController(broker, clientID string) *IoTController {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	return &IoTController{
		client: client,
		coffee: NewCoffeeMaker(),
	}
}

func (iot *IoTController) publishStatus() {
	status := map[string]interface{}{
		"timestamp":   time.Now().Unix(),
		"state":       iot.coffee.state,
		"temperature": iot.coffee.temperature,
		"water_level": iot.coffee.waterLevel,
		"heater_on":   iot.coffee.heaterOn,
		"pump_on":     iot.coffee.pumpOn,
	}

	data, _ := json.Marshal(status)
	iot.client.Publish("coffee/status", 0, false, data)
}

func (iot *IoTController) subscribeCommands() {
	iot.client.Subscribe("coffee/command", 0, func(client mqtt.Client, msg mqtt.Message) {
		var cmd map[string]interface{}
		json.Unmarshal(msg.Payload(), &cmd)

		switch cmd["action"] {
		case "start":
			iot.coffee.waterLevel = 200
		case "stop":
			iot.coffee.heaterOn = false
			iot.coffee.pumpOn = false
		}

		fmt.Printf("Command received: %s\n", cmd["action"])
	})
}

func (iot *IoTController) run() {
	iot.subscribeCommands()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			iot.coffee.tick()
			iot.publishStatus()
		}
	}
}

func main() {
	// Use AWS IoT Core or local MQTT broker
	broker := "ssl://your-iot-endpoint.amazonaws.com:8883"
	// broker := "tcp://localhost:1883" // Local testing

	iot := NewIoTController(broker, "coffee-maker-001")
	fmt.Println("☁️ Connected to IoT cloud")
	iot.run()
}
