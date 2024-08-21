package simpleloadbalancer

import (
	"math/rand"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
)

var name = "simple_balancer"

var logger = grpclog.Component("simple_balancer")

type simplePickerBuilder struct {
	base.PickerBuilder
}

func (b *simplePickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	logger.Infof("simplePickerBuilder: Build called with info: %+v", info)

	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &simplePicker{
		subConns: scs,
		length:   len(scs),
	}
}

func (b *simplePickerBuilder) Name() string {
	return name
}

type simplePicker struct {
	balancer.Picker
	subConns []balancer.SubConn
	length   int
}

func (p *simplePicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	logger.Infof("Picking a new address. len: %v", len(p.subConns))
	index := rand.Intn(p.length)
	return balancer.PickResult{SubConn: p.subConns[index]}, nil
}

func init() {
	balancer.Register(base.NewBalancerBuilder(name, &simplePickerBuilder{}, base.Config{}))
}
