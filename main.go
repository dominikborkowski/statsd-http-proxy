package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/GoMetric/statsd-http-proxy/proxy"
)

var (
	// Version is a current git commit hash and tag
	// Injected by compilation flag
	Version = "Unknown"

	// BuildNumber is a current commit hash
	// Injected by compilation flag
	BuildNumber = "Unknown"

	// BuildDate is a date of build
	// Injected by compilation flag
	BuildDate = "Unknown"
)

func main() {
	// Declare command line options
	var httpHost = flag.String("http-host", getEnv("HTTP_HOST", "127.0.0.1"), "HTTP Host")
	var httpPort = flag.Int("http-port", getEnvInt("HTTP_PORT", 8825), "HTTP Port")
	var httpReadTimeout = flag.Int("http-timeout-read", getEnvInt("HTTP_TIMEOUT_READ", 1), "Read timeout in seconds")
	var httpWriteTimeout = flag.Int("http-timeout-write", getEnvInt("HTTP_TIMEOUT_WRITE", 1), "Write timeout in seconds")
	var httpIdleTimeout = flag.Int("http-timeout-idle", getEnvInt("HTTP_TIMEOUT_IDLE", 1), "Idle timeout in seconds")
	var tlsCert = flag.String("tls-cert", getEnv("TLS_CERT", ""), "TLS certificate for HTTPS")
	var tlsKey = flag.String("tls-key", getEnv("TLS_KEY", ""), "TLS private key for HTTPS")
	var statsdHost = flag.String("statsd-host", getEnv("STATSD_HOST", "127.0.0.1"), "StatsD Host")
	var statsdPort = flag.Int("statsd-port", getEnvInt("STATSD_PORT", 8125), "StatsD Port")
	var metricPrefix = flag.String("metric-prefix", getEnv("METRIC_PREFIX", ""), "Metric name prefix")
	var tokenSecret = flag.String("jwt-secret", getEnv("JWT_SECRET", ""), "Secret to encrypt JWT")
	var verbose = flag.Bool("verbose", getEnvBool("VERBOSE", false), "Verbose")
	var version = flag.Bool("version", false, "Show version")
	var httpRouterName = flag.String("http-router-name", getEnv("HTTP_ROUTER_NAME", "HttpRouter"), "Type of HTTP router")
	var statsdClientName = flag.String("statsd-client-name", getEnv("STATSD_CLIENT_NAME", "GoMetric"), "Type of StatsD client")
	var profilerHTTPort = flag.Int("profiler-http-port", getEnvInt("PROFILER_HTTP_PORT", 0), "Start profiler localhost")

	// Parse flags
	flag.Parse()

	// Show version and exit
	if *version {
		fmt.Printf(
			"StatsD HTTP Proxy v.%s, build %s from %s\n",
			Version,
			BuildNumber,
			BuildDate,
		)
		os.Exit(0)
	}

	// Log build version
	log.Printf(
		"Starting StatsD HTTP Proxy v.%s, build %s from %s\n",
		Version,
		BuildNumber,
		BuildDate,
	)

	// Start profiler
	if *profilerHTTPort > 0 {
		startProfiler(*profilerHTTPort)
	}

	// Start proxy server
	proxyServer := proxy.NewServer(
		*httpHost,
		*httpPort,
		*httpReadTimeout,
		*httpWriteTimeout,
		*httpIdleTimeout,
		*statsdHost,
		*statsdPort,
		*tlsCert,
		*tlsKey,
		*metricPrefix,
		*tokenSecret,
		*verbose,
		*httpRouterName,
		*statsdClientName,
	)

	proxyServer.Listen()
}

// Helper function to get an environment variable or return a default value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Helper function to get an integer environment variable or return a default value
func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// Helper function to get a boolean environment variable or return a default value
func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}

// Helper function to start the profiler
func startProfiler(port int) {
	// Enable block profiling
	runtime.SetBlockProfileRate(1)

	// Start debug server
	profilerHTTPAddress := fmt.Sprintf("localhost:%d", port)
	go func() {
		log.Println("Profiler started at " + profilerHTTPAddress)
		log.Println("Open 'http://" + profilerHTTPAddress + "/debug/pprof/' in your browser or use 'go tool pprof http://" + profilerHTTPAddress + "/debug/pprof/heap' from the console")
		log.Println("See details about pprof in https://golang.org/pkg/net/http/pprof/")
		log.Println(http.ListenAndServe(profilerHTTPAddress, nil))
	}()
}
