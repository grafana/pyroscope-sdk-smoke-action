package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect"
	querierv1 "github.com/grafana/pyroscope/api/gen/proto/go/querier/v1"
	"github.com/grafana/pyroscope/api/gen/proto/go/querier/v1/querierv1connect"
)

func exitf(format string, a ...any) {
	format = "error: " + format + "\n"
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func errorf(format string, a ...any) {
	format = format + "\n"
	fmt.Fprintf(os.Stderr, format, a...)
}

func queryProfile(qc querierv1connect.QuerierServiceClient, term string, req *querierv1.SelectMergeProfileRequest) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := qc.SelectMergeProfile(ctx, connect.NewRequest(req))
	if err != nil {
		return false, err
	}

	for _, s := range res.Msg.StringTable {
		if strings.Contains(s, term) {
			return true, nil
		}
	}
	return false, nil
}

func main() {
	var (
		pyroscopeURL string
		timeout      time.Duration

		profileType string
		serviceName string
		searchTerm  string
	)

	flag.StringVar(&pyroscopeURL, "server", "http://localhost:4040", "Pyroscope server address")
	flag.DurationVar(&timeout, "timeout", 60*time.Second, "Timeout in seconds")
	flag.StringVar(&profileType, "profile-type", "process_cpu:cpu:nanoseconds:cpu:nanoseconds", "Profile type ID")
	flag.StringVar(&serviceName, "service-name", "", "service_name label")
	flag.StringVar(&searchTerm, "term", "", "Search term")
	flag.Parse()

	if profileType == "" {
		exitf("profile-type is required")
	}

	if serviceName == "" {
		exitf("service-name is required")
	}

	if searchTerm == "" {
		exitf("term is required")
	}

	qc := querierv1connect.NewQuerierServiceClient(http.DefaultClient, pyroscopeURL)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	n := 0
	done := false
	for !done {
		select {
		case <-ctx.Done():
			exitf("timeout: %q not found after %d attempts", searchTerm, n)
		case <-ticker.C:
			n++
			end := time.Now()
			start := end.Add(-time.Hour)

			found, err := queryProfile(qc, searchTerm, &querierv1.SelectMergeProfileRequest{
				ProfileTypeID: profileType,
				LabelSelector: fmt.Sprintf(`{service_name="%s"}`, serviceName),
				Start:         start.UnixMilli(),
				End:           end.UnixMilli(),
			})
			if err != nil {
				errorf("[%d] failed to query Pyroscope: %v", n, err)
				continue
			}

			if !found {
				errorf("[%d] %s not found", n, searchTerm)
				continue
			}

			fmt.Printf("[%d] %s is found\n", n, searchTerm)
			done = true
		}
	}
}
