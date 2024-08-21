# grpc-random-load-balancer

This repository contains the code for a custom gRPC load balancer that randomly selects a backend to send the RPC. For more information about the implementation and design decision look at [this blog post](https://farzad1132.github.io/post/grpc-custom-lb/).


# Installation
- You will need Go 1.21 or higher
- In the project directory, run `go mod tidy` to install all the required packages automatically.


# Running the code

## Servers
We will run two servers listening on port `50051` and `50052`.

Running the first server:
```bash
go run greeter_server/main.go --port 50051
```

The second server:
```bash
go run greeter_server/main.go --port 50052
```


## Client

The client uses a service config file to config it's load balancer. In the root of the project, run the following to invoke the client:
```bash
CUSTOM_CONFIG_PATH=./service_config.json go run greeter_client/main.go
```

If you want to see the logs, run the following instead of the above:
```bash
GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info CUSTOM_CONFIG_PATH=./service_config.json go run greeter_client/main.go
```


# References

This code is based on the [helloworlld example](https://github.com/grpc/grpc-go/tree/master/examples/helloworld) provided in the gRPC Go repository.