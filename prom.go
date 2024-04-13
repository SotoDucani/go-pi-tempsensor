package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	temperature prometheus.Gauge
	tempUnits   string
	pressure    prometheus.Gauge
	humidity    prometheus.Gauge
}

func (em *PrometheusMetrics) Init(TempUnits string) {
	switch em.tempUnits {
	case "Fahrenheit":
		em.temperature = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "temperature_fahrenheit",
			Help: "Current ambient temperature in fahrenheit",
		})
	case "Celsius":
		em.temperature = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "temperature_celsius",
			Help: "Current ambient temperature in celsius",
		})
	default:
		em.temperature = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "temperature_fahrenheit",
			Help: "Current ambient temperature in fahrenheit",
		})
	}

	em.pressure = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pressure_kilopascals",
		Help: "Current ambient pressure in kilopascals",
	})

	em.humidity = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "relative_humidity_percentage",
		Help: "The current relative humitidy as a percentage",
	})
}

func ServePromServer(em *PrometheusMetrics) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(em.temperature)
	registry.MustRegister(em.pressure)
	registry.MustRegister(em.humidity)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
