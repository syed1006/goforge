// Package metrics serves Prometheus metrics for demo.
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns an http.Handler that exposes Prometheus metrics. The default
// registry already includes Go runtime and process collectors.
func Handler() http.Handler {
	return promhttp.Handler()
}
