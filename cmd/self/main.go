package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	pyroscope "github.com/grafana/pyroscope-go"
)

func smokeTestMagicFunction() {
	x := 0
	for i := 0; i < 1_000_000; i++ {
		x += i
	}
	runtime.KeepAlive(x)
}

func main() {
	serverAddr := os.Getenv("PYROSCOPE_SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "http://localhost:4040"
	}
	serviceName := os.Getenv("PYROSCOPE_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "smoke-test-app"
	}

	profiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: serviceName,
		ServerAddress:   serverAddr,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "pyroscope start error: %v\n", err)
		os.Exit(1)
	}
	defer profiler.Stop()

	fmt.Printf("profiling %s → %s\n", serviceName, serverAddr)

	deadline := time.Now().Add(90 * time.Second)
	for time.Now().Before(deadline) {
		smokeTestMagicFunction()
	}
}
