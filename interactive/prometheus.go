package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}
