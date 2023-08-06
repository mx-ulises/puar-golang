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
	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "puar_golang_http_request_total",
			Help: "Total number of HTTP requests received in this Golang App with labels",
		},
		[]string{"path"},
	)

	requestError = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "puar_golang_http_request_error",
			Help: "Number of  HTTP requests errors received in this Golang App with labels",
		},
		[]string{"path"},
	)

	transtactionCoinTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "puar_golang_transaction_coin_total",
			Help: "Total number of coins transactioned",
		},
		[]string{"path"},
	)

	coinBalance = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "puar_golang_coin_balance",
			Help: "Current balance of coins in this app (Saved - Spent)",
		},
	)

	averageLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "puar_golang_http_latency_seconds",
			Help: "Average latency in seconds",
		},
		[]string{"path"},
	)

	percentileLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "puar_golang_http_latency_percentiles_seconds",
			Help:    "Percentile latency in seconds",
			Buckets: prometheus.DefBuckets, // Default linear buckets
		},
		[]string{"path"},
	)

	spendingAllowStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "puar_golang_account_lock_status",
			Help: "Status of the spending lock",
		},
		[]string{"status"},
	)

	appVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "puar_golang_version",
			Help: "App Version running",
		},
		[]string{"version"},
	)
)

// Create the children to reduce lookup time:
var (
	// Metric children for '/v1/save' path
	savePathLabel          = prometheus.Labels{"path": "/v1/save"}
	saveRequestCount       = requestTotal.With(savePathLabel)
	saveRequestErrorCount  = requestError.With(savePathLabel)
	savedCoinCount         = transtactionCoinTotal.With(savePathLabel)
	savedAverageLatency    = averageLatency.With(savePathLabel)
	savedPercentileLatency = percentileLatency.With(savePathLabel)

	// Metric children for '/v1/spend' path
	spendPathLabel         = prometheus.Labels{"path": "/v1/spend"}
	spendRequestCount      = requestTotal.With(spendPathLabel)
	spendRequestErrorCount = requestError.With(spendPathLabel)
	spentCoinCount         = transtactionCoinTotal.With(spendPathLabel)
	spendAverageLatency    = averageLatency.With(spendPathLabel)
	spendPercentileLatency = percentileLatency.With(spendPathLabel)

	// Spend block status Metric
	spendPathLockLabel   = prometheus.Labels{"status": "BLOCK"}
	spendStatusLock      = spendingAllowStatus.With(spendPathLockLabel)
	spendPathUnlockLabel = prometheus.Labels{"status": "UNBLOCK"}
	spendStatusUnlock    = spendingAllowStatus.With(spendPathUnlockLabel)
)

var allowSpending = true

func init() {
	prometheus.MustRegister(requestTotal)
	prometheus.MustRegister(requestError)
	prometheus.MustRegister(transtactionCoinTotal)
	prometheus.MustRegister(coinBalance)
	prometheus.MustRegister(averageLatency)
	prometheus.MustRegister(percentileLatency)
	prometheus.MustRegister(spendingAllowStatus)
	prometheus.MustRegister(appVersion)

	appVersionLabel := prometheus.Labels{"version": os.Getenv("VERSION")}
	appVersion.With(appVersionLabel).Set(1.0)
	spendStatusLock.Set(0.0)
	spendStatusUnlock.Set(1.0)
}

func GeneralHandler(
	requestCount *prometheus.Counter,
	requestErrorCount *prometheus.Counter,
	coinCount *prometheus.Counter,
	averageLatency *prometheus.Observer,
	percentileLatency *prometheus.Observer,
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
	(*averageLatency).Observe(diff)
	(*percentileLatency).Observe(diff)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	GeneralHandler(&saveRequestCount, &saveRequestErrorCount, &savedCoinCount, &savedAverageLatency, &savedPercentileLatency, "Save", &w, r, 1)
}

func SpendHandler(w http.ResponseWriter, r *http.Request) {
	if allowSpending {
		GeneralHandler(&spendRequestCount, &spendRequestErrorCount, &spentCoinCount, &spendAverageLatency, &spendPercentileLatency, "Spend", &w, r, -1)
	}
}

func LockSpending(w http.ResponseWriter, r *http.Request) {
	allowSpending = false
	spendStatusLock.Set(1.0)
	spendStatusUnlock.Set(0.0)
	fmt.Fprint(w, "Successful lock operation")
}

func UnlockSpending(w http.ResponseWriter, r *http.Request) {
	allowSpending = true
	spendStatusLock.Set(0.0)
	spendStatusUnlock.Set(1.0)
	fmt.Fprint(w, "Successful unlock operation")
}

func main() {
	http.HandleFunc("/v1/save", SaveHandler)
	http.HandleFunc("/v1/spend", SpendHandler)
	http.HandleFunc("/v1/lock", LockSpending)
	http.HandleFunc("/v1/unlock", UnlockSpending)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
