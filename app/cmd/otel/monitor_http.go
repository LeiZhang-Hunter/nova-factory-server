package main

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"nova-factory-server/app/business/iot/metric/monitor/monitorcontroller"
	"nova-factory-server/app/business/iot/metric/monitor/monitorservice/monitorserviceimpl"
)

func newMonitorHTTPServer() *http.Server {
	service := monitorserviceimpl.NewOntologyMetricReportService(zap.L())
	handler := monitorcontroller.NewOntologyMetricHTTPHandler(service, zap.L())

	mux := http.NewServeMux()
	mux.Handle("/otel/monitor", handler)

	return &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
