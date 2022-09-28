package gometrics

import (
	"net/http"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	//prometheusmetrics "github.com/deathowl/go-metrics-prometheus"
	//"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
)

type TrackedMetrics struct {
	ServiceRequestTimer   metrics.Timer
	LivenessRequestTimer  metrics.Timer
	ReadinessRequestTimer metrics.Timer
	DBGetTimer            metrics.Timer
	DBPutTimer            metrics.Timer
	HTTPStatus1xx         metrics.Counter
	HTTPStatus2xx         metrics.Counter
	HTTPStatus3xx         metrics.Counter
	HTTPStatus4xx         metrics.Counter
	HTTPStatus400         metrics.Counter
	HTTPStatus401         metrics.Counter
	HTTPStatus402         metrics.Counter
	HTTPStatus403         metrics.Counter
	HTTPStatus404         metrics.Counter
	HTTPStatus5xx         metrics.Counter
	HTTPStatus500         metrics.Counter
	HTTPStatus501         metrics.Counter
	HTTPStatus502         metrics.Counter
	HTTPStatus503         metrics.Counter
}

type GoMetrics struct {
	sync.Mutex     // embbed sync.Mutex to add Lock() and Unlock() for thread safety
	Registry       metrics.Registry
	ExpHandler     http.Handler
	TrackedMetrics TrackedMetrics
}

func NewGoMetrics() *GoMetrics {
	gm := &GoMetrics{}

	registry := metrics.NewRegistry()

	gm.SetMetricsRegistry(registry)
	gm.CreateMetrics()

	return gm
}

func (gm *GoMetrics) SetMetricsRegistry(registry metrics.Registry) {
	gm.Lock()
	defer gm.Unlock()
	gm.Registry = registry
	gm.ExpHandler = exp.ExpHandler(gm.Registry)
}

func (gm *GoMetrics) CreateMetrics() {
	gm.TrackedMetrics.ServiceRequestTimer = *gm.CreateTimer("ServiceRequestTimer")
	gm.TrackedMetrics.LivenessRequestTimer = *gm.CreateTimer("LivenessRequestTimer")
	gm.TrackedMetrics.ReadinessRequestTimer = *gm.CreateTimer("ReadinessRequestTimer")
	gm.TrackedMetrics.DBGetTimer = *gm.CreateTimer("DBGetTimer")
	gm.TrackedMetrics.DBPutTimer = *gm.CreateTimer("DBPutTimer")
	gm.TrackedMetrics.HTTPStatus1xx = *gm.CreateCounter("HTTPStatus1xx")
	gm.TrackedMetrics.HTTPStatus2xx = *gm.CreateCounter("HTTPStatus2xx")
	gm.TrackedMetrics.HTTPStatus3xx = *gm.CreateCounter("HTTPStatus3xx")
	gm.TrackedMetrics.HTTPStatus4xx = *gm.CreateCounter("HTTPStatus4xx")
	gm.TrackedMetrics.HTTPStatus400 = *gm.CreateCounter("HTTPStatus400")
	gm.TrackedMetrics.HTTPStatus401 = *gm.CreateCounter("HTTPStatus401")
	gm.TrackedMetrics.HTTPStatus402 = *gm.CreateCounter("HTTPStatus402")
	gm.TrackedMetrics.HTTPStatus403 = *gm.CreateCounter("HTTPStatus403")
	gm.TrackedMetrics.HTTPStatus404 = *gm.CreateCounter("HTTPStatus404")
	gm.TrackedMetrics.HTTPStatus5xx = *gm.CreateCounter("HTTPStatus5xx")
	gm.TrackedMetrics.HTTPStatus500 = *gm.CreateCounter("HTTPStatus500")
	gm.TrackedMetrics.HTTPStatus501 = *gm.CreateCounter("HTTPStatus501")
	gm.TrackedMetrics.HTTPStatus502 = *gm.CreateCounter("HTTPStatus502")
	gm.TrackedMetrics.HTTPStatus503 = *gm.CreateCounter("HTTPStatus503")
}

func (gm *GoMetrics) CreateCounter(name string) *metrics.Counter {
	ctr := metrics.GetOrRegisterCounter(name, gm.Registry)
	metrics.Register(name, ctr)
	return &ctr
}

func (gm *GoMetrics) CreateTimer(name string) *metrics.Timer {
	tmr := metrics.GetOrRegisterTimer(name, gm.Registry)
	return &tmr
}

func (gm *GoMetrics) IncServiceRequestTimer(start time.Time) {
	end := time.Now().UTC()
	duration := end.Sub(start)
	gm.TrackedMetrics.ServiceRequestTimer.Update(duration)
}

func (gm *GoMetrics) IncLivenessRequestTimer(start time.Time) {
	end := time.Now().UTC()
	duration := end.Sub(start)
	gm.TrackedMetrics.LivenessRequestTimer.Update(duration)
}

func (gm *GoMetrics) IncReadinessRequestTimer(start time.Time) {
	end := time.Now().UTC()
	duration := end.Sub(start)
	gm.TrackedMetrics.ReadinessRequestTimer.Update(duration)
}

func (gm *GoMetrics) IncDBGetTimer(start time.Time) {
	end := time.Now().UTC()
	duration := end.Sub(start)
	gm.TrackedMetrics.DBGetTimer.Update(duration)
}

func (gm *GoMetrics) IncDBPutTimer(start time.Time) {
	end := time.Now().UTC()
	duration := end.Sub(start)
	gm.TrackedMetrics.DBPutTimer.Update(duration)
}

func (gm *GoMetrics) IncHTTPStatusCounters(httpStatusCode int) {
	// Translate Status Code to counter(s)
	switch {
	case httpStatusCode >= 100 && httpStatusCode < 200:
		gm.TrackedMetrics.HTTPStatus1xx.Inc(1)
	case httpStatusCode >= 200 && httpStatusCode < 300:
		gm.TrackedMetrics.HTTPStatus2xx.Inc(1)
	case httpStatusCode >= 300 && httpStatusCode < 400:
		gm.TrackedMetrics.HTTPStatus3xx.Inc(1)
	case httpStatusCode == 400:
		gm.TrackedMetrics.HTTPStatus4xx.Inc(1)
		gm.TrackedMetrics.HTTPStatus400.Inc(1)
	case httpStatusCode == 401:
		gm.TrackedMetrics.HTTPStatus4xx.Inc(1)
		gm.TrackedMetrics.HTTPStatus401.Inc(1)
	case httpStatusCode == 402:
		gm.TrackedMetrics.HTTPStatus4xx.Inc(1)
		gm.TrackedMetrics.HTTPStatus402.Inc(1)
	case httpStatusCode == 403:
		gm.TrackedMetrics.HTTPStatus4xx.Inc(1)
		gm.TrackedMetrics.HTTPStatus403.Inc(1)
	case httpStatusCode == 404:
		gm.TrackedMetrics.HTTPStatus4xx.Inc(1)
		gm.TrackedMetrics.HTTPStatus404.Inc(1)
	case httpStatusCode > 404 && httpStatusCode < 500:
		gm.TrackedMetrics.HTTPStatus4xx.Inc(1)
	case httpStatusCode == 500:
		gm.TrackedMetrics.HTTPStatus5xx.Inc(1)
		gm.TrackedMetrics.HTTPStatus500.Inc(1)
	case httpStatusCode == 501:
		gm.TrackedMetrics.HTTPStatus5xx.Inc(1)
		gm.TrackedMetrics.HTTPStatus501.Inc(1)
	case httpStatusCode == 502:
		gm.TrackedMetrics.HTTPStatus5xx.Inc(1)
		gm.TrackedMetrics.HTTPStatus502.Inc(1)
	case httpStatusCode == 503:
		gm.TrackedMetrics.HTTPStatus5xx.Inc(1)
		gm.TrackedMetrics.HTTPStatus503.Inc(1)
	case httpStatusCode > 503 && httpStatusCode < 600:
		gm.TrackedMetrics.HTTPStatus5xx.Inc(1)
	default:
		// ToDo: Log missed status.
	}
}

/*
func AddPrometheusClientRegistry(metricsRegistry metrics.Registry, nameSpace string, serviceName string) {
	flushInterval := time.Duration(1 * time.Second)

	prometheusClient := prometheusmetrics.NewPrometheusProvider(metricsRegistry, nameSpace, serviceName, prometheus.DefaultRegisterer, flushInterval)

	go prometheusClient.UpdatePrometheusMetrics()
}

func StartPrometheusMetricsEndpoint(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}
*/
