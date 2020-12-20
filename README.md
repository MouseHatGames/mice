# mice

[![PkgGoDev](https://pkg.go.dev/badge/github.com/MouseHatGames/mice)](https://pkg.go.dev/github.com/MouseHatGames/mice)

Opinionated microservice framework designed for deployment on Kubernetes.

Minimal example:

```go
package main

import (
	"github.com/MouseHatGames/mice"
	"github.com/MouseHatGames/mice-plugins/codec/json"
	"github.com/MouseHatGames/mice-plugins/transport/grpc"
)

func main() {
	svc := mice.NewService(
		options.Name("my-service"),
		options.RPCPort(8080),
		json.Codec(),
		grpc.Transport(),
	)

	if err := svc.Start(); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}
```
