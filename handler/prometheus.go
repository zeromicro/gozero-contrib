package handler

import (
	"net/http"
	"strconv"

	"github.com/wqyjh/zero-contrib/handler/response"
	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/timex"
)

const serverNamespace = "http_server"

// PrometheusHandler returns a middleware that reports stats to prometheus.
func PrometheusHandler(path, method string, opts ...PrometheusOption) func(http.Handler) http.Handler {
	options := prometheusOptions{
		buckets: []float64{5, 10, 25, 50, 100, 250, 500, 750, 1000},
	}
	for _, o := range opts {
		o(&options)
	}

	metricServerReqDur := metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path", "method", "code"},
		Buckets:   options.buckets,
	})

	metricServerReqCodeTotal := metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "method", "code"},
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := timex.Now()
			cw := response.NewWithCodeResponseWriter(w)
			defer func() {
				code := strconv.Itoa(cw.Code)
				metricServerReqDur.Observe(timex.Since(startTime).Milliseconds(), path, method, code)
				metricServerReqCodeTotal.Inc(path, method, code)
			}()

			next.ServeHTTP(cw, r)
		})
	}
}

type prometheusOptions struct {
	buckets []float64
}

// PrometheusOption allows for managing prometheus options.
type PrometheusOption func(*prometheusOptions)

// WithBuckets sets the buckets for the prometheus metrics.
func WithBuckets(buckets []float64) PrometheusOption {
	return func(o *prometheusOptions) {
		o.buckets = buckets
	}
}
