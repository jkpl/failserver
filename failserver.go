package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	requestDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "http_request_time_ms",
			Help:       "Time spent on requests",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
	)
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"status"},
	)
	maxLatency = int64(getIntEnv("MAX_LATENCY_MS"))
)

func getIntEnv(envKey string) int {
	envStr := os.Getenv(envKey)
	i, _ := strconv.Atoi(envStr)
	return i
}

func simulateLatency() {
	if maxLatency > 0 {
		time.Sleep(time.Duration(rand.Int63n(maxLatency)) * time.Millisecond)
	}
}

func statusLabel(status int) prometheus.Labels {
	return prometheus.Labels{"status": fmt.Sprintf("%d", status)}
}

func requestDurationTrack(start time.Time) {
	elapsed := time.Since(start)
	elapsedMillis := float64(elapsed / time.Millisecond)
	requestDuration.Observe(elapsedMillis)
}

func handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	defer requestDurationTrack(now)
	simulateLatency()

	switch n := rand.Intn(100); n {
	case 4:
		httpRequests.With(statusLabel(http.StatusNotFound)).Inc()
		http.Error(w, "Could not find your lucky number!", http.StatusNotFound)
	case 5:
		httpRequests.With(statusLabel(http.StatusInternalServerError)).Inc()
		http.Error(w, "Failed to compute your lucky number!", http.StatusInternalServerError)
	default:
		httpRequests.With(statusLabel(http.StatusOK)).Inc()
		fmt.Fprintf(w, "Your lucky number is %d", n)
	}
}

func main() {
	prometheus.MustRegister(httpRequests, requestDuration)

	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
