package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/emmanuel971/ovs-exporter/pkg/collector"

)

func main() {
	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		collector.CollectOvsMetrics()
		promhttp.Handler().ServeHTTP(w, r)
	}))

	fmt.Println("OVS exporter running on :9101")
	http.ListenAndServe(":9101", nil)
}
