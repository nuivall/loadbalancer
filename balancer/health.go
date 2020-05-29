package balancer

import "time"

const CheckIntervalSeconds = 2

type HealthChecker struct {
	B LoadBalancer
	t *time.Ticker
}

func (h *HealthChecker) Start() {
	h.B.HealthCheck()
	h.t = time.NewTicker(CheckIntervalSeconds * time.Second)
	go func() {
		for range h.t.C {
			h.B.HealthCheck()
		}
	}()
}

func (h *HealthChecker) Stop() {
	h.t.Stop()
}