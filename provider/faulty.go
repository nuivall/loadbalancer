package provider

import (
	"example.com/iptiq/balancer"
	"fmt"
)

type FaultyProvider struct {
	Id string
	FaultDuration int
}

var _ balancer.Provider = &FaultyProvider{}

func (p *FaultyProvider) Get() string {
	if p.FaultDuration > 0 {
		panic("can't execute get on a faulty provider")
	}
	return p.Id
}

func (p *FaultyProvider) ID() string {
	return p.Id
}

func (p *FaultyProvider) Check() bool {
	if p.FaultDuration > 0 {
		p.FaultDuration--
	}
	r := p.FaultDuration <= 0
	fmt.Println("Health check on", p.Id, "result", r)
	return r
}