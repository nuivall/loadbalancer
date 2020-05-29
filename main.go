package main

import (
	"example.com/iptiq/balancer"
	"example.com/iptiq/provider"
	"fmt"
	"time"
)

func main() {
	runBalancerWithStrategy(balancer.StrategyRoundRobin)
	runBalancerWithStrategy(balancer.StrategyRandom)
	runEmptyBalancer()
}

func runBalancerWithStrategy(s balancer.Strategy) {
	fmt.Println("-----------------------------------------")
	fmt.Printf("Running balancer with strategy %s...\n", s)
	b := balancer.New(s)
	registerProviders(b)
	exerciseBalancer(b)
	exerciseFaultyProvider(b)
}

func runEmptyBalancer() {
	fmt.Println("-----------------------------------------")
	fmt.Printf("Running empty balancer...\n")
	runGetRange(balancer.New(balancer.StrategyRoundRobin))
}

func registerProviders(b balancer.LoadBalancer) {
	b.Register(&provider.SimpleProvider{Id: "simple1"})
	b.Register(&provider.SimpleProvider{Id: "simple2"})
	b.Register(&provider.SimpleProvider{Id: "simple3"})

	b.Register(&provider.FaultyProvider{Id: "faulty1_permanent",
		FaultDuration: 3000})
	b.Register(&provider.FaultyProvider{Id: "faulty2_temporary",
		FaultDuration: 2})
}

func exerciseBalancer(b balancer.LoadBalancer) {
	b.Exclude("simple1")
	runGetRange(b)
	b.Include("simple1")
	runGetRange(b)
}

func exerciseFaultyProvider(b balancer.LoadBalancer) {
	fmt.Println("Make sure some health checks have passed...")
	time.Sleep((2*balancer.CheckIntervalSeconds+1)*time.Second)
	runGetRange(b)
}

func runGetRange(b balancer.LoadBalancer) {
	for i:=0; i<5; i++ {
		g, err := b.Get()
		if err != nil {
			fmt.Printf("Get() error = %v \n", err)
		} else {
			fmt.Printf("Get() = %s \n", g)
		}
	}
}