package metrics

import (
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type ServerRuntimeMetrics struct {
	requestsInCurrentSecond atomic.Int64
	requestsPerSecond       prometheus.Gauge
	serverUp                prometheus.Gauge
}

func NewServerRuntimeMetrics() *ServerRuntimeMetrics {
	requestsPerSecond := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cuoph26_requests_per_second",
		Help: "Observed requests per second in the last 1-second window",
	})

	serverUp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cuoph26_server_up",
		Help: "Server status: 1 when process is running",
	})

	prometheus.MustRegister(requestsPerSecond, serverUp)

	m := &ServerRuntimeMetrics{
		requestsPerSecond: requestsPerSecond,
		serverUp:          serverUp,
	}

	m.serverUp.Set(1)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			count := m.requestsInCurrentSecond.Swap(0)
			m.requestsPerSecond.Set(float64(count))
		}
	}()

	return m
}

func (m *ServerRuntimeMetrics) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() != "/metrics" {
			m.requestsInCurrentSecond.Add(1)
		}
		return c.Next()
	}
}
