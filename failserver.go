package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	requestDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_microseconds",
			Help:       "Time spent on HTTP requests",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
	)
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"code"},
	)
	maxLatency = int64(getIntEnv("MAX_LATENCY_MS"))
	version    = strings.TrimSpace(runCommand(exec.Command("git", "rev-parse", "HEAD")))
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

func statusCodeLabel(status int) prometheus.Labels {
	return prometheus.Labels{"code": fmt.Sprintf("%d", status)}
}

func requestDurationTrack(start time.Time) {
	elapsed := float64(time.Since(start) / time.Microsecond)
	requestDuration.Observe(elapsed)
}

func runCommand(cmd *exec.Cmd) string {
	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}

func handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	defer requestDurationTrack(now)
	simulateLatency()

	switch n := rand.Intn(100); n {
	case 4:
		httpRequests.With(statusCodeLabel(http.StatusNotFound)).Inc()
		http.Error(w, "Could not find your lucky number!", http.StatusNotFound)
	case 5:
		httpRequests.With(statusCodeLabel(http.StatusInternalServerError)).Inc()
		http.Error(w, "Failed to compute your lucky number!", http.StatusInternalServerError)
	default:
		httpRequests.With(statusCodeLabel(http.StatusOK)).Inc()
		fmt.Fprintf(w, "Your lucky number is %d", n)
	}
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}

func main() {
	prometheus.MustRegister(httpRequests, requestDuration)

	http.HandleFunc("/", handler)
	http.HandleFunc("/version", versionHandler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
