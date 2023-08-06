package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	saveRequestCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "puar_golang_save_http_request_count",
			Help: "Total number of save HTTP requests received in this Golang App",
		},
	)

	saveRequestErrorCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "puar_golang_save_http_request_error_count",
			Help: "Total number of save HTTP requests errors received in this Golang App",
		},
	)

	savedCoinCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "puar_golang_saved_coin_count",
			Help: "Total number of coins saved",
		},
	)

	spendRequestCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "puar_golang_spend_http_request_count",
			Help: "Total number of spend HTTP requests received in this Golang App",
		},
	)

	spendRequestErrorCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "puar_golang_spend_http_request_error_count",
			Help: "Total number of spend HTTP requests errors received in this Golang App",
		},
	)

	spentCoinCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "puar_golang_spent_coin_count",
			Help: "Total number of coins spent",
		},
	)

	coinBalance = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "puar_golang_coin_balance",
			Help: "Current balance of coins in this app (Saved - Spent)",
		},
	)

	averageLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "puar_golang_latency_seconds",
			Help: "Average latency in seconds",
		},
	)

	percentileLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "puar_golang_latency_percentiles_seconds",
			Help:    "Percentile latency in seconds",
			Buckets: prometheus.DefBuckets, // Default linear buckets
		},
	)
)

func init() {
	prometheus.MustRegister(saveRequestCount)
	prometheus.MustRegister(saveRequestErrorCount)
	prometheus.MustRegister(savedCoinCount)
	prometheus.MustRegister(spendRequestCount)
	prometheus.MustRegister(spendRequestErrorCount)
	prometheus.MustRegister(spentCoinCount)
	prometheus.MustRegister(coinBalance)
	prometheus.MustRegister(averageLatency)
	prometheus.MustRegister(percentileLatency)
}

func GeneralHandler(
	requestCount *prometheus.Counter,
	requestErrorCount *prometheus.Counter,
	coinCount *prometheus.Counter,
	operation string,
	w *http.ResponseWriter,
	r *http.Request,
	sign float64,
) {
	// Get start time
	startTime := time.Now()
	// Increase request count
	(*requestCount).Inc()

	// Exception Handler
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprint(*w, "Request Error: ", err)
			(*requestErrorCount).Inc()
		}
	}()

	// Evaluate and increase coin count
	queryParams := r.URL.Query()
	coinsString := queryParams.Get("coins")
	if coinsString == "" {
		coinsString = "1"
	}

	coins, err := strconv.ParseFloat(coinsString, 64)
	if err != nil {
		panic(err)
	}
	if coins <= 0 {
		panic("Negative number of coins provided")
	}

	// Evaluate Random success rate at 20%
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Float64()
	if randomNumber < 0.2 {
		panic("Random Error")
	}

	// Generate random sleep
	duration := time.Duration(1000*randomNumber) * time.Millisecond
	time.Sleep(duration)

	// Increase number of coins
	(*coinCount).Add(coins)
	coinBalance.Add(coins * sign)

	// Successful request
	fmt.Fprint(*w, "Successful ", operation, " Request, coins: ", coins)
	endTime := time.Now()
	diff := float64(endTime.Sub(startTime).Seconds())
	averageLatency.Observe(diff)
	percentileLatency.Observe(diff)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	GeneralHandler(&saveRequestCount, &saveRequestErrorCount, &savedCoinCount, "Save", &w, r, 1)
}

func SpendHandler(w http.ResponseWriter, r *http.Request) {
	GeneralHandler(&spendRequestCount, &spendRequestErrorCount, &spentCoinCount, "Spend", &w, r, -1)
}

func main() {
	http.HandleFunc("/v1/save", SaveHandler)
	http.HandleFunc("/v1/spend", SpendHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
