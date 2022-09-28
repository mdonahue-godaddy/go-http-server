package gometrics

import (
	"net/http"
	"time"

	"github.com/rcrowley/go-metrics"

	prometheusmetrics "github.com/deathowl/go-metrics-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricLogger interface {
	RequestResponseTime()
	IncrementServiceRequests()
	IncrementHealthCheckRequests()
	IncrementHTTPStatus1xx()
	IncrementHTTPStatus2xx()
	IncrementHTTPStatus3xx()
	IncrementHTTPStatus4xx()
	IncrementHTTPStatus400()
	IncrementHTTPStatus401()
	IncrementHTTPStatus402()
	IncrementHTTPStatus403()
	IncrementHTTPStatus404()
	IncrementHTTPStatus5xx()
	IncrementHTTPStatus500()
	IncrementHTTPStatus501()
	IncrementHTTPStatus502()
	IncrementHTTPStatus503()
}

type GoMetrics struct {
	ServiceRequestResponseTimer     metrics.Timer
	HealthCheckRequestResponseTimer metrics.Timer
	DBGetTimer                      metrics.Timer
	DBPutTimer                      metrics.Timer
	HTTPStatus1xx                   metrics.Counter
	HTTPStatus2xx                   metrics.Counter
	HTTPStatus3xx                   metrics.Counter
	HTTPStatus4xx                   metrics.Counter
	HTTPStatus400                   metrics.Counter
	HTTPStatus401                   metrics.Counter
	HTTPStatus402                   metrics.Counter
	HTTPStatus403                   metrics.Counter
	HTTPStatus404                   metrics.Counter
	HTTPStatus5xx                   metrics.Counter
	HTTPStatus500                   metrics.Counter
	HTTPStatus501                   metrics.Counter
	HTTPStatus502                   metrics.Counter
	HTTPStatus503                   metrics.Counter
}

func AddPrometheusClientRegistry(metricsRegistry metrics.Registry, nameSpace string, serviceName string) {
	flushInterval := time.Duration(1 * time.Second)

	prometheusClient := prometheusmetrics.NewPrometheusProvider(metricsRegistry, nameSpace, serviceName, prometheus.DefaultRegisterer, flushInterval)

	go prometheusClient.UpdatePrometheusMetrics()
}

func StartPrometheusMetricsEndpoint(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}

func NewGoMetrics() *GoMetrics {
	metrics := GoMetrics{}

	return &metrics
}

func (m *GoMetrics) IncrementHTTPStatusCounters(httpStatusCode int) {
	// Translate Status Code to counter(s)
	switch {
	case httpStatusCode >= 100 && httpStatusCode < 200:
		m.HTTPStatus1xx.Inc(1)
	case httpStatusCode >= 200 && httpStatusCode < 300:
		m.HTTPStatus2xx.Inc(1)
	case httpStatusCode >= 300 && httpStatusCode < 400:
		m.HTTPStatus3xx.Inc(1)
	case httpStatusCode == 400:
		m.HTTPStatus4xx.Inc(1)
		m.HTTPStatus400.Inc(1)
	case httpStatusCode == 401:
		m.HTTPStatus4xx.Inc(1)
		m.HTTPStatus401.Inc(1)
	case httpStatusCode == 402:
		m.HTTPStatus4xx.Inc(1)
		m.HTTPStatus402.Inc(1)
	case httpStatusCode == 403:
		m.HTTPStatus4xx.Inc(1)
		m.HTTPStatus403.Inc(1)
	case httpStatusCode == 404:
		m.HTTPStatus4xx.Inc(1)
		m.HTTPStatus404.Inc(1)
	case httpStatusCode > 404 && httpStatusCode < 500:
		m.HTTPStatus4xx.Inc(1)
	case httpStatusCode == 500:
		m.HTTPStatus5xx.Inc(1)
		m.HTTPStatus500.Inc(1)
	case httpStatusCode == 501:
		m.HTTPStatus5xx.Inc(1)
		m.HTTPStatus501.Inc(1)
	case httpStatusCode == 502:
		m.HTTPStatus5xx.Inc(1)
		m.HTTPStatus502.Inc(1)
	case httpStatusCode == 503:
		m.HTTPStatus5xx.Inc(1)
		m.HTTPStatus503.Inc(1)
	case httpStatusCode > 503 && httpStatusCode < 600:
		m.HTTPStatus5xx.Inc(1)
	}
}
