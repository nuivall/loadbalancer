package balancer

type Strategy string
const (
	StrategyRandom Strategy = "random"
	StrategyRoundRobin Strategy = "round-robin"
)

type LoadBalancer interface {
	Register(p Provider) error
	Get() (string, error)
	Include(id string)
	Exclude(id string)
	HealthCheck()
}

type Provider interface {
	Get() string
	Check() bool
	ID() string
}