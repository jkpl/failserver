package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	requestTime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "failserver_request_time_ms",
			Help: "Time spent on requests",
		},
	)
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "failserver_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"status"},
	)
	maxLatency = int64(200)
)

func simulateLatency() {
	time.Sleep(time.Duration(rand.Int63n(maxLatency)) * time.Millisecond)
}

func statusLabel(status int) prometheus.Labels {
	return prometheus.Labels{"status": fmt.Sprintf("%d", status)}
}

func requestTimeTrack(start time.Time) {
	elapsed := time.Since(start)
	elapsedMillis := float64(elapsed / time.Millisecond)
	requestTime.Set(elapsedMillis)
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer requestTimeTrack(time.Now())
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
	prometheus.MustRegister(httpRequests, requestTime)

	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
