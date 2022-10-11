package gometrics

import (
	"context"
	"fmt"

	metrics "github.com/rcrowley/go-metrics"

	"go.elastic.co/apm/v2"
)

type gatherer struct {
	registry    metrics.Registry
	percentiles []float64
}

// WrapRegistry() wraps metrics.Registry so that it can be used as an apm.MetricsGatherer.
func WrapRegistry(reg metrics.Registry) apm.MetricsGatherer {
	return gatherer{
		registry:    reg,
		percentiles: []float64{0.5, 0.75, 0.95, 0.99, 0.999},
	}
}

// GatherMetrics gathers metrics into m.
func (g gatherer) GatherMetrics(ctx context.Context, met *apm.Metrics) error {
	g.registry.Each(func(name string, v interface{}) {
		switch v := v.(type) {
		case metrics.Counter: // results and values should be consistant with https://github.com/rcrowley/go-metrics/blob/cf1acfcdf4751e0554ffa765d03e479ec491cad6/exp/exp.go#L81
			met.Add(name, nil, float64(v.Count()))
			/* Not needed yet, but I don't wanto just delete working code
			case metrics.Gauge: // results and values should be consistant with https://github.com/rcrowley/go-metrics/blob/cf1acfcdf4751e0554ffa765d03e479ec491cad6/exp/exp.go#L86
				met.Add(name, nil, float64(v.Value()))
			case metrics.GaugeFloat64: // results and values should be consistant with https://github.com/rcrowley/go-metrics/blob/cf1acfcdf4751e0554ffa765d03e479ec491cad6/exp/exp.go#L90
				met.Add(name, nil, v.Value())
			case metrics.Histogram: // results and values should be consistant with https://github.com/rcrowley/go-metrics/blob/cf1acfcdf4751e0554ffa765d03e479ec491cad6/exp/exp.go#L94
				h := v.Snapshot()
				ps := h.Percentiles(g.percentiles)
				met.Add(name+".count", nil, float64(h.Count()))
				met.Add(name+".sum", nil, float64(h.Sum()))
				met.Add(name+".min", nil, float64(h.Min()))
				met.Add(name+".max", nil, float64(h.Max()))
				met.Add(name+".std_dev", nil, h.StdDev())
				for idx := 0; idx < len(g.percentiles); idx++ {
					displayPercent := int(g.percentiles[idx] * 100)
					if g.percentiles[idx] > 0.99 {
						displayPercent = int(g.percentiles[idx] * 1000)
					}
					met.Add(fmt.Sprintf("%s.percentile.%d", name, displayPercent), nil, h.Percentile(ps[idx]))
				}
			case metrics.Meter: // results and values should be consistant with https://github.com/rcrowley/go-metrics/blob/cf1acfcdf4751e0554ffa765d03e479ec491cad6/exp/exp.go#L109
				m := v.Snapshot()
				met.Add(name+".count", nil, float64(m.Count()))
				met.Add(name+".one_minute", nil, m.Rate1())
				met.Add(name+".five_minute", nil, m.Rate5())
				met.Add(name+".fifteen_minute", nil, m.Rate15())
				met.Add(name+".mean_rate", nil, m.RateMean())
			*/
		case metrics.Timer: // results and values should be consistant with https://github.com/rcrowley/go-metrics/blob/cf1acfcdf4751e0554ffa765d03e479ec491cad6/exp/exp.go#L118
			t := v.Snapshot()
			ps := t.Percentiles(g.percentiles)
			met.Add(name+".count", nil, float64(t.Count()))
			met.Add(name+".max", nil, float64(t.Max()))
			met.Add(name+".mean", nil, t.Mean())
			met.Add(name+".min", nil, float64(t.Min()))
			met.Add(name+".std_dev", nil, float64(t.StdDev()))
			met.Add(name+".mean_rate", nil, t.RateMean())
			met.Add(name+".sum", nil, float64(t.Sum()))
			met.Add(name+".variance", nil, t.Variance())
			met.Add(name+".one_minute", nil, t.Rate1())
			met.Add(name+".five_minute", nil, t.Rate5())
			met.Add(name+".fifteen_minute", nil, t.Rate15())
			for idx := 0; idx < len(g.percentiles); idx++ {
				displayPercent := int(g.percentiles[idx] * 100)
				if g.percentiles[idx] > 0.99 {
					displayPercent = int(g.percentiles[idx] * 1000)
				}
				met.Add(fmt.Sprintf("%s.percentile.%d", name, displayPercent), nil, t.Percentile(ps[idx]))
			}
		default:
			// TODO: other types (metrics.EWMA)
		}
	})

	return nil
}
