package provider

import "example.com/iptiq/balancer"

type SimpleProvider struct {
	Id string
}

var _ balancer.Provider = &SimpleProvider{}

func (p *SimpleProvider) Get() string {
	return p.Id
}

func (p *SimpleProvider) ID() string {
	return p.Id
}

func (p *SimpleProvider) Check() bool {
	return true
}