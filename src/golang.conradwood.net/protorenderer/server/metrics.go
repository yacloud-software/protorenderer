package main

import (
	"golang.conradwood.net/go-easyops/prometheus"
	"time"
)

var (
	failedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "protorenderer_failed_protos",
			Help: "V=1 UNIT=none DESC=number of proto files failed to compile",
		},
	)
)

func start_metrics() {
	prometheus.MustRegister(failedGauge)
	go start_metrics_loop()
}
func start_metrics_loop() {
	t := time.Duration(3) * time.Second
	for {
		time.Sleep(t)
		t = time.Duration(30) * time.Second
		metrics_update()
	}
}
func metrics_update() {
	v := completeVersion
	if v == nil {
		return
	}
	failedGauge.Set(float64(len(v.failures.Failures())))
}






