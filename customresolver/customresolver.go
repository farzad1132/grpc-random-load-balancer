package customresolver

import (
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

const customScheme = "custom"

var customServiceName = "custom.service.name"
var backendAddr = []string{"localhost:50051", "localhost:50052"}

type CustomResolverBuilder struct{}

var logger = grpclog.Component("CustomResolver")

func (*CustomResolverBuilder) Build(target resolver.Target,
	cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	logger.Info("Using 'custom' scheme resolver")
	r := &customResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			customServiceName: backendAddr,
		},
	}
	r.start()
	return r, nil
}
func (*CustomResolverBuilder) Scheme() string { return customScheme }

// customResolver is a
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type customResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *customResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}

	state := resolver.State{Addresses: addrs}
	r.cc.UpdateState(state)
}
func (*customResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*customResolver) Close()                                  {}

func init() {
	// Register the example ResolverBuilder. This is usually done in a package's
	// init() function.
	resolver.Register(&CustomResolverBuilder{})
}
