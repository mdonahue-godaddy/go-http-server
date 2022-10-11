package gometrics_test

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"

	"github.com/mdonahue-godaddy/go-http-server/log"
	"github.com/mdonahue-godaddy/go-http-server/metrics/gometrics"
)

var (
	writeBuffer  bytes.Buffer
	metricPrefix = "app.unit.test"
)

func createTestLogger() *log.Logger {
	logger := log.NewLogger()
	writer := bufio.NewWriter(&writeBuffer)
	logger.Logger = logger.Logger.Output(writer)
	return &logger
}

// Test_StartMetricsLogger test starting loggerr
func Test_StartMetricsLogger(t *testing.T) {
	logger := createTestLogger()
	interval := time.Duration(time.Second * 10)

	loggingTestRegistry := metrics.NewRegistry()

	gm := gometrics.NewGoMetrics(loggingTestRegistry, metricPrefix)

	before := writeBuffer.Len()

	gm.StartMetricsLogger(logger, interval)

	gm.TrackedMetrics.ServiceRequest.Update(time.Duration(1 * time.Second))

	time.Sleep(interval)

	after := writeBuffer.Len()

	assert.Greater(t, after, before)
}

// Test_EnableExpHandler verify ExpHandler has been setup
func Test_EnableExpHandler(t *testing.T) {
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.EnableExpHandler()

	assert.NotNil(t, gm.ExpHandler)
}

// Test_GetMetricsPrefix verify metrics are setup with correct type and work as expected.
func Test_GetMetricsPrefix(t *testing.T) {
	expected := metricPrefix

	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, expected)

	actual := gm.GetMetricsPrefix()

	assert.Equal(t, expected, actual)
}

// Test_ServiceRequest verify metrics are setup with correct type and work as expected.
func Test_ServiceRequest(t *testing.T) {
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	assert.IsType(t, &metrics.StandardTimer{}, gm.TrackedMetrics.ServiceRequest)
	assert.Equal(t, int64(0), gm.TrackedMetrics.ServiceRequest.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.ServiceRequest.Max())
	gm.TrackedMetrics.ServiceRequest.Update(time.Duration(1 * time.Second))
	assert.Equal(t, int64(1), gm.TrackedMetrics.ServiceRequest.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.ServiceRequest.Max())
	gm.IncServiceRequest(time.Since(time.Now().UTC().Add(-2 * time.Second)))
	assert.Equal(t, int64(2), gm.TrackedMetrics.ServiceRequest.Count())
	assert.GreaterOrEqual(t, gm.TrackedMetrics.ServiceRequest.Min(), int64(time.Duration(1*time.Second)))
	assert.GreaterOrEqual(t, gm.TrackedMetrics.ServiceRequest.Max(), int64(time.Duration(2*time.Second)))
}

// Test_HealthRequest verify metrics are setup with correct type and work as expected.
func Test_HealthRequest(t *testing.T) {
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	assert.IsType(t, &metrics.StandardTimer{}, gm.TrackedMetrics.HealthRequest)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HealthRequest.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HealthRequest.Max())
	gm.TrackedMetrics.HealthRequest.Update(time.Duration(1 * time.Second))
	assert.Equal(t, int64(1), gm.TrackedMetrics.HealthRequest.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.HealthRequest.Max())
	gm.IncHealthRequest(time.Since(time.Now().UTC().Add(-2 * time.Second)))
	assert.Equal(t, int64(2), gm.TrackedMetrics.HealthRequest.Count())
	assert.GreaterOrEqual(t, gm.TrackedMetrics.HealthRequest.Min(), int64(time.Duration(1*time.Second)))
	assert.GreaterOrEqual(t, gm.TrackedMetrics.HealthRequest.Max(), int64(time.Duration(2*time.Second)))
}

// Test_MetricsRequest verify metrics are setup with correct type and work as expected.
func Test_MetricsRequest(t *testing.T) {
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	assert.IsType(t, &metrics.StandardTimer{}, gm.TrackedMetrics.MetricsRequest)
	assert.Equal(t, int64(0), gm.TrackedMetrics.MetricsRequest.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.MetricsRequest.Max())
	gm.TrackedMetrics.MetricsRequest.Update(time.Duration(1 * time.Second))
	assert.Equal(t, int64(1), gm.TrackedMetrics.MetricsRequest.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.MetricsRequest.Max())
	gm.IncMetricRequest(time.Since(time.Now().UTC().Add(-2 * time.Second)))
	assert.Equal(t, int64(2), gm.TrackedMetrics.MetricsRequest.Count())
	assert.GreaterOrEqual(t, gm.TrackedMetrics.MetricsRequest.Min(), int64(time.Duration(1*time.Second)))
	assert.GreaterOrEqual(t, gm.TrackedMetrics.MetricsRequest.Max(), int64(time.Duration(2*time.Second)))
}

// Test_HTTPHealthStatus1xx verify metrics are setup with correct type and work as expected.
func Test_HTTPHealthStatus1xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPHealth.Status1xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPHealth.Status1xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.Status1xx.Count())
	gm.TrackedMetrics.HTTPHealth.Status1xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPHealth.Status1xx.Count())
	gm.IncHTTPHealth(logger, 100, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPHealth.Status1xx.Count())
}

// Test_HTTPHealthStatus2xx verify metrics are setup with correct type and work as expected.
func Test_HTTPHealthStatus2xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPHealth.Status2xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPHealth.Status2xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.Status2xx.Count())
	gm.TrackedMetrics.HTTPHealth.Status2xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPHealth.Status2xx.Count())
	gm.IncHTTPHealth(logger, 200, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPHealth.Status2xx.Count())
}

// Test_HTTPHealthStatus3xx verify metrics are setup with correct type and work as expected.
func Test_HTTPHealthStatus3xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPHealth.Status3xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPHealth.Status3xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.Status3xx.Count())
	gm.TrackedMetrics.HTTPHealth.Status3xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPHealth.Status3xx.Count())
	gm.IncHTTPHealth(logger, 300, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPHealth.Status3xx.Count())
}

// Test_HTTPHealthStatus4xx verify metrics are setup with correct type and work as expected.
func Test_HTTPHealthStatus4xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPHealth.Status4xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPHealth.Status4xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.Status4xx.Count())
	gm.TrackedMetrics.HTTPHealth.Status4xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPHealth.Status4xx.Count())
	gm.IncHTTPHealth(logger, 422, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPHealth.Status4xx.Count())
}

// Test_HTTPHealthStatus5xx verify metrics are setup with correct type and work as expected.
func Test_HTTPHealthStatus5xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPHealth.Status5xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPHealth.Status5xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.Status5xx.Count())
	gm.TrackedMetrics.HTTPHealth.Status5xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPHealth.Status5xx.Count())
	gm.IncHTTPHealth(logger, 504, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPHealth.Status5xx.Count())
}

func Test_HTTPHealthStatus000(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPHealth.StatusOOR.Clear()

	gm.IncHTTPHealth(logger, 0, duration)

	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPHealth.StatusOOR.Count())
}

func Test_HTTPHealthStatus600(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPHealth.StatusOOR.Clear()

	gm.IncHTTPHealth(logger, 600, duration)

	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPHealth.StatusOOR.Count())
}

// Test_HTTPMetricStatus1xx verify metrics are setup with correct type and work as expected.
func Test_HTTPMetricStatus1xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPMetric.Status1xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPMetric.Status1xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.Status1xx.Count())
	gm.TrackedMetrics.HTTPMetric.Status1xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPMetric.Status1xx.Count())
	gm.IncHTTPMetric(logger, 100, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPMetric.Status1xx.Count())
}

// Test_HTTPMetricStatus2xx verify metrics are setup with correct type and work as expected.
func Test_HTTPMetricStatus2xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPMetric.Status2xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPMetric.Status2xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.Status2xx.Count())
	gm.TrackedMetrics.HTTPMetric.Status2xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPMetric.Status2xx.Count())
	gm.IncHTTPMetric(logger, 200, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPMetric.Status2xx.Count())
}

// Test_HTTPMetricStatus3xx verify metrics are setup with correct type and work as expected.
func Test_HTTPMetricStatus3xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPMetric.Status3xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPMetric.Status3xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.Status3xx.Count())
	gm.TrackedMetrics.HTTPMetric.Status3xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPMetric.Status3xx.Count())
	gm.IncHTTPMetric(logger, 300, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPMetric.Status3xx.Count())
}

// Test_HTTPMetricStatus4xx verify metrics are setup with correct type and work as expected.
func Test_HTTPMetricStatus4xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPMetric.Status4xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPMetric.Status4xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.Status4xx.Count())
	gm.TrackedMetrics.HTTPMetric.Status4xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPMetric.Status4xx.Count())
	gm.IncHTTPMetric(logger, 422, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPMetric.Status4xx.Count())
}

// Test_HTTPMetricStatus5xx verify metrics are setup with correct type and work as expected.
func Test_HTTPMetricStatus5xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPMetric.Status5xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPMetric.Status5xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.Status5xx.Count())
	gm.TrackedMetrics.HTTPMetric.Status5xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPMetric.Status5xx.Count())
	gm.IncHTTPMetric(logger, 504, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPMetric.Status5xx.Count())
}

func Test_HTTPMetricStatus000(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPMetric.StatusOOR.Clear()

	gm.IncHTTPMetric(logger, 0, duration)

	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPMetric.StatusOOR.Count())
}

func Test_HTTPMetricStatus600(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPMetric.StatusOOR.Clear()

	gm.IncHTTPMetric(logger, 600, duration)

	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPMetric.StatusOOR.Count())
}

// Test_HTTPServiceStatus1xx verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus1xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status1xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status1xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status1xx.Count())
	gm.TrackedMetrics.HTTPService.Status1xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status1xx.Count())
	gm.IncHTTPService(logger, 100, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status1xx.Count())
}

// Test_HTTPServiceStatus2xx verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus2xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status2xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status2xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status2xx.Count())
	gm.TrackedMetrics.HTTPService.Status2xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status2xx.Count())
	gm.IncHTTPService(logger, 200, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status2xx.Count())
}

// Test_HTTPServiceStatus3xx verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus3xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status3xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status3xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status3xx.Count())
	gm.TrackedMetrics.HTTPService.Status3xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status3xx.Count())
	gm.IncHTTPService(logger, 300, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status3xx.Count())
}

// Test_HTTPServiceStatus4xx verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus4xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status4xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status4xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status4xx.Count())
	gm.TrackedMetrics.HTTPService.Status4xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status4xx.Count())
	gm.IncHTTPService(logger, 422, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status4xx.Count())
}

// Test_HTTPServiceStatus400 verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus400(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status400.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status400)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status400.Count())
	gm.TrackedMetrics.HTTPService.Status400.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status400.Count())
	gm.IncHTTPService(logger, 400, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status400.Count())
}

// Test_HTTPServiceStatus401 verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus401(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status401.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status401)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status401.Count())
	gm.TrackedMetrics.HTTPService.Status401.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status401.Count())
	gm.IncHTTPService(logger, 401, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status401.Count())
}

// Test_HTTPServiceStatus402 verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus402(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status402.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status402)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status402.Count())
	gm.TrackedMetrics.HTTPService.Status402.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status402.Count())
	gm.IncHTTPService(logger, 402, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status402.Count())
}

// Test_HTTPServiceStatus403 verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus403(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status403.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status403)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status403.Count())
	gm.TrackedMetrics.HTTPService.Status403.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status403.Count())
	gm.IncHTTPService(logger, 403, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status403.Count())
}

// Test_HTTPServiceStatus404 verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus404(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status404.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status404)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status404.Count())
	gm.TrackedMetrics.HTTPService.Status404.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status404.Count())
	gm.IncHTTPService(logger, 404, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status404.Count())
}

// Test_HTTPServiceStatus5xx verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus5xx(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status5xx.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status5xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status5xx.Count())
	gm.TrackedMetrics.HTTPService.Status5xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status5xx.Count())
	gm.IncHTTPService(logger, 504, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status5xx.Count())
}

// Test_HTTPServiceStatus500 verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus500(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status500.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status500)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status500.Count())
	gm.TrackedMetrics.HTTPService.Status500.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status500.Count())
	gm.IncHTTPService(logger, 500, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status500.Count())
}

// Test_HTTPServiceStatus501 verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus501(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status501.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status501)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status501.Count())
	gm.TrackedMetrics.HTTPService.Status501.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status501.Count())
	gm.IncHTTPService(logger, 501, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status501.Count())
}

// Test_HTTPServiceStatus502 verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus502(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status502.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status502)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status502.Count())
	gm.TrackedMetrics.HTTPService.Status502.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status502.Count())
	gm.IncHTTPService(logger, 502, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status502.Count())
}

// Test_HTTPServiceStatus503 verify metrics are setup with correct type and work as expected.
func Test_HTTPServiceStatus503(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.Status503.Clear()

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPService.Status503)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status503.Count())
	gm.TrackedMetrics.HTTPService.Status503.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.Status503.Count())
	gm.IncHTTPService(logger, 503, duration)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPService.Status503.Count())
}

func Test_HTTPServiceStatus000(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.StatusOOR.Clear()

	gm.IncHTTPService(logger, 0, duration)

	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.StatusOOR.Count())
}

func Test_HTTPServiceStatus600(t *testing.T) {
	logger := createTestLogger()
	duration := time.Duration(time.Second * 1)
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	gm.TrackedMetrics.HTTPService.StatusOOR.Clear()

	gm.IncHTTPService(logger, 600, duration)

	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPService.StatusOOR.Count())
}

func Test_ResetCounters(t *testing.T) {
	// create struct instance
	gm := gometrics.NewGoMetrics(metrics.DefaultRegistry, metricPrefix)

	// Ensure we start with a clean slate
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

	// Increment all of the counters by one
	gm.TrackedMetrics.HTTPHealth.Status1xx.Inc(1)
	gm.TrackedMetrics.HTTPHealth.Status2xx.Inc(1)
	gm.TrackedMetrics.HTTPHealth.Status3xx.Inc(1)
	gm.TrackedMetrics.HTTPHealth.Status4xx.Inc(1)
	gm.TrackedMetrics.HTTPHealth.Status5xx.Inc(1)
	gm.TrackedMetrics.HTTPHealth.StatusOOR.Inc(1)
	gm.TrackedMetrics.HTTPMetric.Status1xx.Inc(1)
	gm.TrackedMetrics.HTTPMetric.Status2xx.Inc(1)
	gm.TrackedMetrics.HTTPMetric.Status3xx.Inc(1)
	gm.TrackedMetrics.HTTPMetric.Status4xx.Inc(1)
	gm.TrackedMetrics.HTTPMetric.Status5xx.Inc(1)
	gm.TrackedMetrics.HTTPMetric.StatusOOR.Inc(1)
	gm.TrackedMetrics.HTTPService.Status1xx.Inc(1)
	gm.TrackedMetrics.HTTPService.Status2xx.Inc(1)
	gm.TrackedMetrics.HTTPService.Status3xx.Inc(1)
	gm.TrackedMetrics.HTTPService.Status4xx.Inc(1)
	gm.TrackedMetrics.HTTPService.Status400.Inc(1)
	gm.TrackedMetrics.HTTPService.Status401.Inc(1)
	gm.TrackedMetrics.HTTPService.Status402.Inc(1)
	gm.TrackedMetrics.HTTPService.Status403.Inc(1)
	gm.TrackedMetrics.HTTPService.Status404.Inc(1)
	gm.TrackedMetrics.HTTPService.Status5xx.Inc(1)
	gm.TrackedMetrics.HTTPService.Status500.Inc(1)
	gm.TrackedMetrics.HTTPService.Status501.Inc(1)
	gm.TrackedMetrics.HTTPService.Status502.Inc(1)
	gm.TrackedMetrics.HTTPService.Status503.Inc(1)
	gm.TrackedMetrics.HTTPService.StatusOOR.Inc(1)

	// call reset
	gm.ResetCounters()

	// verify each counter was reset
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.Status1xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.Status2xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.Status3xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.Status4xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.Status5xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPHealth.StatusOOR.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.Status1xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.Status2xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.Status3xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.Status4xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.Status5xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPMetric.StatusOOR.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status1xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status2xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status3xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status4xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status400.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status401.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status402.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status403.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status404.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status5xx.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status500.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status501.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status502.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.Status503.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPService.StatusOOR.Count())
}
