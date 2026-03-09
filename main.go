package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

// ──────────────────────────────────────────────
// Global state — request counter and server start time
// ──────────────────────────────────────────────

var (
	requestCount int64
	startTime    = time.Now()
)

// ──────────────────────────────────────────────
// Build metadata — injected at compile time or set as defaults
// ──────────────────────────────────────────────

var (
	appName   = "ip-app-golang"
	version   = "1.4.2"
	commit    = "a3f22b9"
	buildTime = "2026-03-05"
)

// ──────────────────────────────────────────────
// Middleware — counts every inbound HTTP request using an atomic counter
// ──────────────────────────────────────────────

func withMetrics(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&requestCount, 1)
		next(w, r)
	}
}

// ──────────────────────────────────────────────
// Helper — writes a JSON response with the given status code
// ──────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ──────────────────────────────────────────────
// getOutboundIP — discovers the host's preferred outbound IP
// by opening a UDP connection to a public DNS server.
// No actual traffic is sent; the OS simply picks the right interface.
// ──────────────────────────────────────────────

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// ──────────────────────────────────────────────
// Root — returns the host's IP address, hostname, and a greeting
// ──────────────────────────────────────────────

func rootHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	ip := getOutboundIP()

	writeJSON(w, http.StatusOK, map[string]string{
		"hostname": hostname,
		"ip":       ip,
		"message":  "Welcome to ip-app-golang — a simple IP discovery service.",
	})
}

// ──────────────────────────────────────────────
// Health probes — a single handler covers /healthz, /ready, /live, /status
// ──────────────────────────────────────────────

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// ──────────────────────────────────────────────
// Ping — minimal liveness check, returns plain text "pong"
// ──────────────────────────────────────────────

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "pong\n")
}

// ──────────────────────────────────────────────
// Metrics — Prometheus-style plain-text metrics for observability
// ──────────────────────────────────────────────

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	uptime := time.Since(startTime).Seconds()
	count := atomic.LoadInt64(&requestCount)

	// cpu_usage_percent is a placeholder; Go's standard library does not
	// expose per-process CPU usage without calling into OS-specific APIs.
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "http_requests_total %d\n", count)
	fmt.Fprintf(w, "cpu_usage_percent 0.0\n")
	fmt.Fprintf(w, "memory_usage_bytes %d\n", memStats.Alloc)
	fmt.Fprintf(w, "uptime_seconds %.2f\n", uptime)
}

// ──────────────────────────────────────────────
// Info — build metadata and runtime environment details
// ──────────────────────────────────────────────

func infoHandler(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"app":       appName,
		"version":   version,
		"commit":    commit,
		"buildTime": buildTime,
		"go":        runtime.Version(),
		"platform":  runtime.GOOS,
		"arch":      runtime.GOARCH,
		"env":       env,
	})
}

// ──────────────────────────────────────────────
// Main — registers routes and starts the HTTP server
// ──────────────────────────────────────────────

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", withMetrics(rootHandler))
	http.HandleFunc("/healthz", withMetrics(healthHandler))
	http.HandleFunc("/ready", withMetrics(healthHandler))
	http.HandleFunc("/live", withMetrics(healthHandler))
	http.HandleFunc("/status", withMetrics(healthHandler))
	http.HandleFunc("/ping", withMetrics(pingHandler))
	http.HandleFunc("/metrics", withMetrics(metricsHandler))
	http.HandleFunc("/info", withMetrics(infoHandler))

	addr := ":" + port
	log.Printf("ip-app-golang listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
