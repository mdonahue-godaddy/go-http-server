package gometrics

//go:generate mockgen -source=gometrics.go -destination=mocks/gometrics_mocks.go -package=mocks

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/apmtest"
	"go.elastic.co/apm/v2/model"

	//prometheusmetrics "github.com/deathowl/go-metrics-prometheus"
	//"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mdonahue-godaddy/go-http-server/log"
)

var (
	//GlobalMetrics        IGoMetrics
	metricsLoggerRunning bool = false
)

func init() {
	//GlobalMetrics = NewGoMetrics(metrics.DefaultRegistry, "go-http-server")
}

// IGoMetrics is defined to allow easy future mocking of the package whenever needed in the future.
type IGoMetrics interface {
	SetMetricsRegistry(registry metrics.Registry)
	SetMetricsPrefix(prefix string)
	GetMetricsPrefix() string
	EnableExpHandler()
	EnableDebugGCStats(duration time.Duration)
	EnableRuntimeMemStats(duration time.Duration)
	EnableMetricsLogger(logger *log.Logger, duration time.Duration)
	StartMetricsLogger(logger *log.Logger, duration time.Duration)
	CreateMetrics()
	CreateCounter(name string) metrics.Counter
	CreateTimer(name string) metrics.Timer
	IncServiceRequest(duration time.Duration)
	IncHealthRequest(duration time.Duration)
	IncMetricRequest(duration time.Duration)
	IncHTTPHealth(logger *log.Logger, httpStatusCode int, duration time.Duration)
	IncHTTPMetric(logger *log.Logger, httpStatusCode int, duration time.Duration)
	IncHTTPService(logger *log.Logger, httpStatusCode int, duration time.Duration)
}

type HTTPMetrics struct {
	Status1xx metrics.Counter
	Status2xx metrics.Counter
	Status3xx metrics.Counter
	Status4xx metrics.Counter
	Status400 metrics.Counter
	Status401 metrics.Counter
	Status402 metrics.Counter
	Status403 metrics.Counter
	Status404 metrics.Counter
	Status5xx metrics.Counter
	Status500 metrics.Counter
	Status501 metrics.Counter
	Status502 metrics.Counter
	Status503 metrics.Counter
	StatusOOR metrics.Counter // Out-Of-Range (n < 0 or n > 599)
}

type HTTPBasicMetrics struct {
	Status1xx metrics.Counter
	Status2xx metrics.Counter
	Status3xx metrics.Counter
	Status4xx metrics.Counter
	Status5xx metrics.Counter
	StatusOOR metrics.Counter // Out-Of-Range (n < 0 or n > 599)
}

type TrackedMetrics struct {
	ServiceRequest metrics.Timer
	HealthRequest  metrics.Timer
	MetricsRequest metrics.Timer
	HTTPService    HTTPMetrics
	HTTPHealth     HTTPBasicMetrics
	HTTPMetric     HTTPBasicMetrics
}

type GoMetrics struct {
	// Internals
	sync.Mutex // embbed sync.Mutex to add Lock() and Unlock() for thread safety
	registry   metrics.Registry
	prefix     string
	// Exported
	ExpHandler     http.Handler // http server will need access to this
	TrackedMetrics TrackedMetrics
}

func NewGoMetrics(registry metrics.Registry, prefix string) *GoMetrics {
	gm := &GoMetrics{}
	gm.SetMetricsRegistry(registry)
	gm.SetMetricsPrefix(prefix)
	gm.CreateMetrics()

	return gm
}

func (gm *GoMetrics) SetMetricsRegistry(registry metrics.Registry) {
	gm.Lock()
	defer gm.Unlock()
	gm.registry = registry
}

func (gm *GoMetrics) SetMetricsPrefix(prefix string) {
	gm.Lock()
	defer gm.Unlock()
	gm.prefix = prefix
}

func (gm *GoMetrics) GetMetricsPrefix() string {
	return gm.prefix
}

func (gm *GoMetrics) EnableExpHandler() {
	gm.Lock()
	defer gm.Unlock()
	gm.ExpHandler = exp.ExpHandler(gm.registry)
}

func (gm *GoMetrics) EnableDebugGCStats(duration time.Duration) {
	gm.Lock()
	defer gm.Unlock()

	metrics.RegisterDebugGCStats(gm.registry)
	go metrics.CaptureDebugGCStats(gm.registry, duration)
}

func (gm *GoMetrics) EnableRuntimeMemStats(duration time.Duration) {
	gm.Lock()
	defer gm.Unlock()

	metrics.RegisterRuntimeMemStats(gm.registry)
	go metrics.CaptureRuntimeMemStats(gm.registry, duration)
}

func (gm *GoMetrics) EnableMetricsLogger(logger *log.Logger, duration time.Duration) {
	gm.Lock()
	defer gm.Unlock()

	go metrics.Log(gm.registry, duration, logger)
}

func (gm *GoMetrics) StartMetricsLogger(logger *log.Logger, duration time.Duration) {
	if !metricsLoggerRunning {
		go MetricsLogger(gm.registry, logger, duration)
	}
}

func MetricsLogger(registry metrics.Registry, logger *log.Logger, duration time.Duration) {
	metricsLoggerRunning = true

	gatherer := WrapRegistry(registry)

	cpuPercentArr := make([]string, 0)
	memPercent := ""

	// cpu - usage
	cpuPercents, cerr := cpu.Percent(0, false)
	if cerr == nil {
		// convert to string array for dictionary
		cpuPercentArr = make([]string, len(cpuPercents))
		for idx, cpuPercent := range cpuPercents {
			cpuPercentArr[idx] = strconv.FormatFloat(cpuPercent, 'f', 2, 64)
		}
	}

	// memory - usage
	vm, merr := mem.VirtualMemory()
	if merr == nil {
		// convert to string for dictionary
		memPercent = strconv.FormatFloat(vm.UsedPercent, 'f', 2, 64)
	}

	for {
		gms := gatherMetrics(gatherer)

		for _, gm := range gms {
			ctnr := zerolog.Dict()

			cpu := zerolog.Dict()
			cpu.Strs("usage", cpuPercentArr)

			ctnr.Dict("cpu", cpu)

			mem := zerolog.Dict()
			mem.Str("usage", memPercent)

			ctnr.Dict("memory", mem)

			mets := zerolog.Dict()

			keys := make([]string, 0, len(gm.Samples))

			for k := range gm.Samples {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, key := range keys {
				mets.Str(key, fmt.Sprintf("%f", gm.Samples[key].Value))
			}

			info := logger.Info().
				Tags(log.ApplicationTag).
				ECSMetric(log.InformationType).
				Dict("container", ctnr).
				Dict("metric", mets)

			info.Msg("log pstore metrics")
		}

		time.Sleep(duration)
	}
}

func gatherMetrics(mg apm.MetricsGatherer) []model.Metrics {
	tracer := apmtest.NewRecordingTracer()
	defer tracer.Close()

	tracer.RegisterMetricsGatherer(mg)
	tracer.SendMetrics(nil)

	metrics := tracer.Payloads().Metrics

	for idx := range metrics {
		metrics[idx].Timestamp = model.Time{}
	}

	return metrics
}

func (gm *GoMetrics) CreateMetrics() {
	gm.TrackedMetrics.ServiceRequest = gm.CreateTimer(gm.CreateMetricName("http.service.request"))
	gm.TrackedMetrics.HealthRequest = gm.CreateTimer(gm.CreateMetricName("http.health.request"))
	gm.TrackedMetrics.MetricsRequest = gm.CreateTimer(gm.CreateMetricName("http.metric.request"))
	gm.TrackedMetrics.HTTPHealth.Status1xx = gm.CreateCounter(gm.CreateMetricName("http.health.response.status.1xx"))
	gm.TrackedMetrics.HTTPHealth.Status2xx = gm.CreateCounter(gm.CreateMetricName("http.health.response.status.2xx"))
	gm.TrackedMetrics.HTTPHealth.Status3xx = gm.CreateCounter(gm.CreateMetricName("http.health.response.status.3xx"))
	gm.TrackedMetrics.HTTPHealth.Status4xx = gm.CreateCounter(gm.CreateMetricName("http.health.response.status.4xx"))
	gm.TrackedMetrics.HTTPHealth.Status5xx = gm.CreateCounter(gm.CreateMetricName("http.health.response.status.5xx"))
	gm.TrackedMetrics.HTTPHealth.StatusOOR = gm.CreateCounter(gm.CreateMetricName("http.health.response.status.oor"))
	gm.TrackedMetrics.HTTPMetric.Status1xx = gm.CreateCounter(gm.CreateMetricName("http.metric.response.status.1xx"))
	gm.TrackedMetrics.HTTPMetric.Status2xx = gm.CreateCounter(gm.CreateMetricName("http.metric.response.status.2xx"))
	gm.TrackedMetrics.HTTPMetric.Status3xx = gm.CreateCounter(gm.CreateMetricName("http.metric.response.status.3xx"))
	gm.TrackedMetrics.HTTPMetric.Status4xx = gm.CreateCounter(gm.CreateMetricName("http.metric.response.status.4xx"))
	gm.TrackedMetrics.HTTPMetric.Status5xx = gm.CreateCounter(gm.CreateMetricName("http.metric.response.status.5xx"))
	gm.TrackedMetrics.HTTPMetric.StatusOOR = gm.CreateCounter(gm.CreateMetricName("http.metric.response.status.oor"))
	gm.TrackedMetrics.HTTPService.Status1xx = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.1xx"))
	gm.TrackedMetrics.HTTPService.Status2xx = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.2xx"))
	gm.TrackedMetrics.HTTPService.Status3xx = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.3xx"))
	gm.TrackedMetrics.HTTPService.Status4xx = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.4xx"))
	gm.TrackedMetrics.HTTPService.Status400 = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.400"))
	gm.TrackedMetrics.HTTPService.Status401 = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.401"))
	gm.TrackedMetrics.HTTPService.Status402 = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.402"))
	gm.TrackedMetrics.HTTPService.Status403 = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.403"))
	gm.TrackedMetrics.HTTPService.Status404 = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.404"))
	gm.TrackedMetrics.HTTPService.Status5xx = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.5xx"))
	gm.TrackedMetrics.HTTPService.Status500 = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.500"))
	gm.TrackedMetrics.HTTPService.Status501 = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.501"))
	gm.TrackedMetrics.HTTPService.Status502 = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.502"))
	gm.TrackedMetrics.HTTPService.Status503 = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.503"))
	gm.TrackedMetrics.HTTPService.StatusOOR = gm.CreateCounter(gm.CreateMetricName("http.service.response.status.oor"))
}

func (gm *GoMetrics) ResetCounters() {
	// Timers are histograms with a NewExpDecaySample(1028, 0.015) and cannot be cleared.
	gm.TrackedMetrics.HTTPHealth.Status1xx.Clear()
	gm.TrackedMetrics.HTTPHealth.Status2xx.Clear()
	gm.TrackedMetrics.HTTPHealth.Status3xx.Clear()
	gm.TrackedMetrics.HTTPHealth.Status4xx.Clear()
	gm.TrackedMetrics.HTTPHealth.Status5xx.Clear()
	gm.TrackedMetrics.HTTPHealth.StatusOOR.Clear()
	gm.TrackedMetrics.HTTPMetric.Status1xx.Clear()
	gm.TrackedMetrics.HTTPMetric.Status2xx.Clear()
	gm.TrackedMetrics.HTTPMetric.Status3xx.Clear()
	gm.TrackedMetrics.HTTPMetric.Status4xx.Clear()
	gm.TrackedMetrics.HTTPMetric.Status5xx.Clear()
	gm.TrackedMetrics.HTTPMetric.StatusOOR.Clear()
	gm.TrackedMetrics.HTTPService.Status1xx.Clear()
	gm.TrackedMetrics.HTTPService.Status2xx.Clear()
	gm.TrackedMetrics.HTTPService.Status3xx.Clear()
	gm.TrackedMetrics.HTTPService.Status4xx.Clear()
	gm.TrackedMetrics.HTTPService.Status400.Clear()
	gm.TrackedMetrics.HTTPService.Status401.Clear()
	gm.TrackedMetrics.HTTPService.Status402.Clear()
	gm.TrackedMetrics.HTTPService.Status403.Clear()
	gm.TrackedMetrics.HTTPService.Status404.Clear()
	gm.TrackedMetrics.HTTPService.Status5xx.Clear()
	gm.TrackedMetrics.HTTPService.Status500.Clear()
	gm.TrackedMetrics.HTTPService.Status501.Clear()
	gm.TrackedMetrics.HTTPService.Status502.Clear()
	gm.TrackedMetrics.HTTPService.Status503.Clear()
	gm.TrackedMetrics.HTTPService.StatusOOR.Clear()
}

func (gm *GoMetrics) CreateMetricName(detail string) string {
	return fmt.Sprintf("%s.%s", gm.prefix, detail)
}

func (gm *GoMetrics) CreateCounter(name string) metrics.Counter {
	return metrics.GetOrRegisterCounter(name, gm.registry)
}

func (gm *GoMetrics) CreateTimer(name string) metrics.Timer {
	return metrics.GetOrRegisterTimer(name, gm.registry)
}

func (gm *GoMetrics) IncServiceRequest(duration time.Duration) {
	gm.TrackedMetrics.ServiceRequest.Update(duration)
}

func (gm *GoMetrics) IncHealthRequest(duration time.Duration) {
	gm.TrackedMetrics.HealthRequest.Update(duration)
}

func (gm *GoMetrics) IncMetricRequest(duration time.Duration) {
	gm.TrackedMetrics.MetricsRequest.Update(duration)
}

func (gm *GoMetrics) IncHTTPHealth(logger *log.Logger, httpStatusCode int, duration time.Duration) {
	gm.IncHealthRequest(duration)

	// Translate Status Code to counter(s)
	switch {
	case httpStatusCode >= 100 && httpStatusCode < 200:
		gm.TrackedMetrics.HTTPHealth.Status1xx.Inc(1)
	case httpStatusCode >= 200 && httpStatusCode < 300:
		gm.TrackedMetrics.HTTPHealth.Status2xx.Inc(1)
	case httpStatusCode >= 300 && httpStatusCode < 400:
		gm.TrackedMetrics.HTTPHealth.Status3xx.Inc(1)
	case httpStatusCode >= 400 && httpStatusCode < 500:
		gm.TrackedMetrics.HTTPHealth.Status4xx.Inc(1)
	case httpStatusCode >= 500 && httpStatusCode < 600:
		gm.TrackedMetrics.HTTPHealth.Status5xx.Inc(1)
	default:
		gm.TrackedMetrics.HTTPHealth.StatusOOR.Inc(1)
		warning := logger.Warn().
			Dict("warn", zerolog.Dict().Str("httpStatusCode", strconv.Itoa(httpStatusCode))).
			ECSTimedEvent(log.Web, log.NotApplicable, duration, log.InformationType)

		warning.Msg("unexpected health http status code")
	}
}

func (gm *GoMetrics) IncHTTPMetric(logger *log.Logger, httpStatusCode int, duration time.Duration) {
	gm.IncMetricRequest(duration)

	// Translate Status Code to counter(s)
	switch {
	case httpStatusCode >= 100 && httpStatusCode < 200:
		gm.TrackedMetrics.HTTPMetric.Status1xx.Inc(1)
	case httpStatusCode >= 200 && httpStatusCode < 300:
		gm.TrackedMetrics.HTTPMetric.Status2xx.Inc(1)
	case httpStatusCode >= 300 && httpStatusCode < 400:
		gm.TrackedMetrics.HTTPMetric.Status3xx.Inc(1)
	case httpStatusCode >= 400 && httpStatusCode < 500:
		gm.TrackedMetrics.HTTPMetric.Status4xx.Inc(1)
	case httpStatusCode >= 500 && httpStatusCode < 600:
		gm.TrackedMetrics.HTTPMetric.Status5xx.Inc(1)
	default:
		gm.TrackedMetrics.HTTPMetric.StatusOOR.Inc(1)
		warning := logger.Warn().
			Dict("warn", zerolog.Dict().Str("httpStatusCode", strconv.Itoa(httpStatusCode))).
			ECSTimedEvent(log.Web, log.NotApplicable, duration, log.InformationType)

		warning.Msg("unexpected metric http status code")
	}
}

func (gm *GoMetrics) IncHTTPService(logger *log.Logger, httpStatusCode int, duration time.Duration) {
	gm.IncServiceRequest(duration)

	// Translate Status Code to counter(s)
	switch {
	case httpStatusCode >= 100 && httpStatusCode < 200:
		gm.TrackedMetrics.HTTPService.Status1xx.Inc(1)
	case httpStatusCode >= 200 && httpStatusCode < 300:
		gm.TrackedMetrics.HTTPService.Status2xx.Inc(1)
	case httpStatusCode >= 300 && httpStatusCode < 400:
		gm.TrackedMetrics.HTTPService.Status3xx.Inc(1)
	case httpStatusCode == 400:
		gm.TrackedMetrics.HTTPService.Status4xx.Inc(1)
		gm.TrackedMetrics.HTTPService.Status400.Inc(1)
	case httpStatusCode == 401:
		gm.TrackedMetrics.HTTPService.Status4xx.Inc(1)
		gm.TrackedMetrics.HTTPService.Status401.Inc(1)
	case httpStatusCode == 402:
		gm.TrackedMetrics.HTTPService.Status4xx.Inc(1)
		gm.TrackedMetrics.HTTPService.Status402.Inc(1)
	case httpStatusCode == 403:
		gm.TrackedMetrics.HTTPService.Status4xx.Inc(1)
		gm.TrackedMetrics.HTTPService.Status403.Inc(1)
	case httpStatusCode == 404:
		gm.TrackedMetrics.HTTPService.Status4xx.Inc(1)
		gm.TrackedMetrics.HTTPService.Status404.Inc(1)
	case httpStatusCode > 404 && httpStatusCode < 500:
		gm.TrackedMetrics.HTTPService.Status4xx.Inc(1)
	case httpStatusCode == 500:
		gm.TrackedMetrics.HTTPService.Status5xx.Inc(1)
		gm.TrackedMetrics.HTTPService.Status500.Inc(1)
	case httpStatusCode == 501:
		gm.TrackedMetrics.HTTPService.Status5xx.Inc(1)
		gm.TrackedMetrics.HTTPService.Status501.Inc(1)
	case httpStatusCode == 502:
		gm.TrackedMetrics.HTTPService.Status5xx.Inc(1)
		gm.TrackedMetrics.HTTPService.Status502.Inc(1)
	case httpStatusCode == 503:
		gm.TrackedMetrics.HTTPService.Status5xx.Inc(1)
		gm.TrackedMetrics.HTTPService.Status503.Inc(1)
	case httpStatusCode > 503 && httpStatusCode < 600:
		gm.TrackedMetrics.HTTPService.Status5xx.Inc(1)
	default:
		gm.TrackedMetrics.HTTPService.StatusOOR.Inc(1)
		warning := logger.Warn().
			Dict("warn", zerolog.Dict().Str("httpStatusCode", strconv.Itoa(httpStatusCode))).
			ECSTimedEvent(log.Web, log.NotApplicable, duration, log.InformationType)

		warning.Msg("unexpected service http status code")
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
