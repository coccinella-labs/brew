package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CoffeeMetrics struct {
	coffee *CoffeeMaker

	// Prometheus metrics
	temperature    prometheus.Gauge
	waterLevel     prometheus.Gauge
	brewsTotal     prometheus.Counter
	brewDuration   prometheus.Histogram
	systemUptime   prometheus.Counter
}

func NewCoffeeMetrics() *CoffeeMetrics {
	cm := &CoffeeMetrics{
		coffee: NewCoffeeMaker(),

		temperature: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "coffee_temperature_celsius",
			Help: "Current coffee maker temperature",
		}),

		waterLevel: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "coffee_water_level_ml",
			Help: "Current water level in ml",
		}),

		brewsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "coffee_brews_total",
			Help: "Total number of coffee brews",
		}),

		brewDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "coffee_brew_duration_seconds",
			Help:    "Time taken to brew coffee",
			Buckets: prometheus.LinearBuckets(60, 30, 10), // 60s to 330s
		}),

		systemUptime: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "coffee_system_uptime_seconds",
			Help: "System uptime in seconds",
		}),
	}

	// Register metrics
	prometheus.MustRegister(cm.temperature)
	prometheus.MustRegister(cm.waterLevel)
	prometheus.MustRegister(cm.brewsTotal)
	prometheus.MustRegister(cm.brewDuration)
	prometheus.MustRegister(cm.systemUptime)

	return cm
}

func (cm *CoffeeMetrics) updateMetrics() {
	cm.temperature.Set(float64(cm.coffee.temperature))
	cm.waterLevel.Set(float64(cm.coffee.waterLevel))
	cm.systemUptime.Inc()
}

func (cm *CoffeeMetrics) run() {
	// Start metrics server
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":9090", nil))
	}()

	fmt.Println("📊 Metrics server: http://localhost:9090/metrics")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastState := cm.coffee.state
	brewStart := time.Now()

	for {
		select {
		case <-ticker.C:
			cm.coffee.tick()
			cm.updateMetrics()

			// Track brew cycles
			if lastState != "BREWING" && cm.coffee.state == "BREWING" {
				brewStart = time.Now()
			}

			if lastState == "BREWING" && cm.coffee.state == "DONE" {
				duration := time.Since(brewStart).Seconds()
				cm.brewDuration.Observe(duration)
				cm.brewsTotal.Inc()
				fmt.Printf("☕ Brew completed in %.1fs\n", duration)
			}

			lastState = cm.coffee.state
		}
	}
}

func main() {
	metrics := NewCoffeeMetrics()
	metrics.run()
}
