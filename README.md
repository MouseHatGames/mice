# mice

[![PkgGoDev](https://pkg.go.dev/badge/github.com/MouseHatGames/mice)](https://pkg.go.dev/github.com/MouseHatGames/mice)

Opinionated microservice framework designed for deployment on Kubernetes.

Minimal example:

```go
func main() {
	svc := mice.NewService(
		options.Name("my-service"),
		options.ListenAddr(":8080"),
		json.Codec(),
		grpc.Transport(),
	)

	if err := svc.Start(); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}
```
