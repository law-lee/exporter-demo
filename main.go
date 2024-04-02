package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/law-lee/exporter-demo/collect"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "myapp",
		Name:      "processed_ops_total",
		Help:      "The total number of processed events",
	})
)

func main() {
	prometheus.MustRegister(opsProcessed)
	recordMetrics()
	prometheus.MustRegister(collect.NewloadavgCollector())

	http.Handle("/metrics", promhttp.Handler())
	log.Print("export /metrics on port :8085")
	http.ListenAndServe(":8085", nil)
}
