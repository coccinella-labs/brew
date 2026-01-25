package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type OTAManager struct {
	client   mqtt.Client
	version  uint32
	deviceID string
}

func NewOTAManager(broker, deviceID string) *OTAManager {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(deviceID + "-ota")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	return &OTAManager{
		client:   client,
		version:  100, // Current firmware version
		deviceID: deviceID,
	}
}

func (ota *OTAManager) checkForUpdates() {
	// Subscribe to OTA commands
	topic := fmt.Sprintf("ota/%s/command", ota.deviceID)
	ota.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		command := string(msg.Payload())

		switch command {
		case "check_update":
			ota.reportVersion()
		case "start_update":
			ota.downloadUpdate()
		}
	})
}

func (ota *OTAManager) reportVersion() {
	topic := fmt.Sprintf("ota/%s/version", ota.deviceID)
	message := fmt.Sprintf(`{"version": %d, "timestamp": %d}`,
		ota.version, time.Now().Unix())

	ota.client.Publish(topic, 0, false, message)
	fmt.Printf("📡 Reported version: %d\n", ota.version)
}

func (ota *OTAManager) downloadUpdate() {
	fmt.Println("⬇️ Downloading firmware update...")

	// Download from firmware server
	resp, err := http.Get("https://firmware-server.com/coffee-maker/latest.bin")
	if err != nil {
		log.Printf("Download failed: %v", err)
		return
	}
	defer resp.Body.Close()

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "firmware-*.bin")
	if err != nil {
		log.Printf("Temp file error: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Download and verify
	hasher := sha256.New()
	size, err := io.Copy(io.MultiWriter(tmpFile, hasher), resp.Body)
	if err != nil {
		log.Printf("Download error: %v", err)
		return
	}

	checksum := hasher.Sum(nil)
	fmt.Printf("✅ Downloaded %d bytes, SHA256: %x\n", size, checksum[:8])

	// Send to device via MQTT chunks
	ota.sendFirmwareToDevice(tmpFile.Name())
}

func (ota *OTAManager) sendFirmwareToDevice(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("File error: %v", err)
		return
	}
	defer file.Close()

	// Get file size
	stat, _ := file.Stat()
	size := stat.Size()

	// Send firmware header
	header := make([]byte, 16)
	binary.LittleEndian.PutUint32(header[0:4], 0xDEADBEEF) // Magic
	binary.LittleEndian.PutUint32(header[4:8], ota.version+1) // New version
	binary.LittleEndian.PutUint32(header[8:12], uint32(size)) // Size
	binary.LittleEndian.PutUint32(header[12:16], 0x12345678) // CRC placeholder

	topic := fmt.Sprintf("ota/%s/firmware", ota.deviceID)
	ota.client.Publish(topic+"/header", 0, false, header)

	// Send firmware in 1KB chunks
	buffer := make([]byte, 1024)
	offset := uint32(0)

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}

		// Send chunk with offset
		chunkTopic := fmt.Sprintf("%s/chunk/%d", topic, offset)
		ota.client.Publish(chunkTopic, 0, false, buffer[:n])

		offset += uint32(n)
		fmt.Printf("📤 Sent chunk at offset %d\n", offset)
		time.Sleep(100 * time.Millisecond) // Rate limit
	}

	// Signal completion
	ota.client.Publish(topic+"/complete", 0, false, "done")
	fmt.Println("🚀 Firmware update sent to device")
}

func (ota *OTAManager) run() {
	ota.checkForUpdates()
	ota.reportVersion()

	// Keep alive
	for {
		time.Sleep(30 * time.Second)
		ota.reportVersion()
	}
}

func main() {
	ota := NewOTAManager("tcp://localhost:1883", "coffee-maker-001")
	fmt.Println("🔄 OTA Manager started")
	ota.run()
}
