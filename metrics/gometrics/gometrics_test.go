package gometrics_test

import (
	"context"
	"testing"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"

	"github.com/mdonahue-godaddy/go-http-server/metrics/gometrics"
)

// Test_NewGoMetrics verify metrics are setup with correct type and work as expected.
func Test_NewGoMetrics(t *testing.T) {
	ctx := context.Background()
	gm := gometrics.NewGoMetrics()

	assert.IsType(t, &metrics.StandardTimer{}, gm.TrackedMetrics.DBGetTimer)
	assert.Equal(t, int64(0), gm.TrackedMetrics.DBGetTimer.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.DBGetTimer.Max())
	gm.TrackedMetrics.DBGetTimer.Update(time.Duration(1 * time.Second))
	assert.Equal(t, int64(1), gm.TrackedMetrics.DBGetTimer.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.DBGetTimer.Max())
	gm.IncDBGetTimer(time.Now().UTC().Add(-2 * time.Second))
	assert.Equal(t, int64(2), gm.TrackedMetrics.DBGetTimer.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.DBGetTimer.Min())
	assert.Equal(t, int64(time.Duration(2*time.Second)), gm.TrackedMetrics.DBGetTimer.Max())

	assert.IsType(t, &metrics.StandardTimer{}, gm.TrackedMetrics.DBPutTimer)
	assert.Equal(t, int64(0), gm.TrackedMetrics.DBPutTimer.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.DBPutTimer.Max())
	gm.TrackedMetrics.DBPutTimer.Update(time.Duration(1 * time.Second))
	assert.Equal(t, int64(1), gm.TrackedMetrics.DBPutTimer.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.DBPutTimer.Max())
	gm.IncDBPutTimer(time.Now().UTC().Add(-2 * time.Second))
	assert.Equal(t, int64(2), gm.TrackedMetrics.DBPutTimer.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.DBPutTimer.Min())
	assert.Equal(t, int64(time.Duration(2*time.Second)), gm.TrackedMetrics.DBPutTimer.Max())

	assert.IsType(t, &metrics.StandardTimer{}, gm.TrackedMetrics.LivenessRequestTimer)
	assert.Equal(t, int64(0), gm.TrackedMetrics.LivenessRequestTimer.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.LivenessRequestTimer.Max())
	gm.TrackedMetrics.LivenessRequestTimer.Update(time.Duration(1 * time.Second))
	assert.Equal(t, int64(1), gm.TrackedMetrics.LivenessRequestTimer.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.LivenessRequestTimer.Max())
	gm.IncLivenessRequestTimer(time.Now().UTC().Add(-2 * time.Second))
	assert.Equal(t, int64(2), gm.TrackedMetrics.LivenessRequestTimer.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.LivenessRequestTimer.Min())
	assert.Equal(t, int64(time.Duration(2*time.Second)), gm.TrackedMetrics.LivenessRequestTimer.Max())

	assert.IsType(t, &metrics.StandardTimer{}, gm.TrackedMetrics.ReadinessRequestTimer)
	assert.Equal(t, int64(0), gm.TrackedMetrics.ReadinessRequestTimer.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.ReadinessRequestTimer.Max())
	gm.TrackedMetrics.ReadinessRequestTimer.Update(time.Duration(1 * time.Second))
	assert.Equal(t, int64(1), gm.TrackedMetrics.ReadinessRequestTimer.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.ReadinessRequestTimer.Max())
	gm.IncReadinessRequestTimer(time.Now().UTC().Add(-2 * time.Second))
	assert.Equal(t, int64(2), gm.TrackedMetrics.ReadinessRequestTimer.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.ReadinessRequestTimer.Min())
	assert.Equal(t, int64(time.Duration(2*time.Second)), gm.TrackedMetrics.ReadinessRequestTimer.Max())

	assert.IsType(t, &metrics.StandardTimer{}, gm.TrackedMetrics.ServiceRequestTimer)
	assert.Equal(t, int64(0), gm.TrackedMetrics.ServiceRequestTimer.Count())
	assert.Equal(t, int64(0), gm.TrackedMetrics.ServiceRequestTimer.Max())
	gm.TrackedMetrics.ServiceRequestTimer.Update(time.Duration(1 * time.Second))
	assert.Equal(t, int64(1), gm.TrackedMetrics.ServiceRequestTimer.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.ServiceRequestTimer.Max())
	gm.IncServiceRequestTimer(time.Now().UTC().Add(-2 * time.Second))
	assert.Equal(t, int64(2), gm.TrackedMetrics.ServiceRequestTimer.Count())
	assert.Equal(t, int64(time.Duration(1*time.Second)), gm.TrackedMetrics.ServiceRequestTimer.Min())
	assert.Equal(t, int64(time.Duration(2*time.Second)), gm.TrackedMetrics.ServiceRequestTimer.Max())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus1xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus1xx.Count())
	gm.TrackedMetrics.HTTPStatus1xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus1xx.Count())
	gm.IncHTTPStatusCounters(ctx, 100)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus1xx.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus2xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus2xx.Count())
	gm.TrackedMetrics.HTTPStatus2xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus2xx.Count())
	gm.IncHTTPStatusCounters(ctx, 200)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus2xx.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus3xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus3xx.Count())
	gm.TrackedMetrics.HTTPStatus3xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus3xx.Count())
	gm.IncHTTPStatusCounters(ctx, 300)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus3xx.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus4xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus4xx.Count())
	gm.TrackedMetrics.HTTPStatus4xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus4xx.Count())
	gm.IncHTTPStatusCounters(ctx, 422)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus4xx.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus400)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus400.Count())
	gm.TrackedMetrics.HTTPStatus400.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus400.Count())
	gm.IncHTTPStatusCounters(ctx, 400)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus400.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus401)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus401.Count())
	gm.TrackedMetrics.HTTPStatus401.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus401.Count())
	gm.IncHTTPStatusCounters(ctx, 401)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus401.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus402)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus402.Count())
	gm.TrackedMetrics.HTTPStatus402.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus402.Count())
	gm.IncHTTPStatusCounters(ctx, 402)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus402.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus403)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus403.Count())
	gm.TrackedMetrics.HTTPStatus403.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus403.Count())
	gm.IncHTTPStatusCounters(ctx, 403)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus403.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus404)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus404.Count())
	gm.TrackedMetrics.HTTPStatus404.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus404.Count())
	gm.IncHTTPStatusCounters(ctx, 404)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus404.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus5xx)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus5xx.Count())
	gm.TrackedMetrics.HTTPStatus5xx.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus5xx.Count())
	gm.IncHTTPStatusCounters(ctx, 504)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus5xx.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus500)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus500.Count())
	gm.TrackedMetrics.HTTPStatus500.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus500.Count())
	gm.IncHTTPStatusCounters(ctx, 500)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus500.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus501)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus501.Count())
	gm.TrackedMetrics.HTTPStatus501.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus501.Count())
	gm.IncHTTPStatusCounters(ctx, 501)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus501.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus502)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus502.Count())
	gm.TrackedMetrics.HTTPStatus502.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus502.Count())
	gm.IncHTTPStatusCounters(ctx, 502)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus502.Count())

	assert.IsType(t, &metrics.StandardCounter{}, gm.TrackedMetrics.HTTPStatus503)
	assert.Equal(t, int64(0), gm.TrackedMetrics.HTTPStatus503.Count())
	gm.TrackedMetrics.HTTPStatus503.Inc(1)
	assert.Equal(t, int64(1), gm.TrackedMetrics.HTTPStatus503.Count())
	gm.IncHTTPStatusCounters(ctx, 503)
	assert.Equal(t, int64(2), gm.TrackedMetrics.HTTPStatus503.Count())

	//gm.IncHTTPStatusCounters(ctx, 0)

	//gm.IncHTTPStatusCounters(ctx, 600)
}
