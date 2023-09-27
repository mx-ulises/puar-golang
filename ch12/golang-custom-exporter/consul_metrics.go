package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	up = prometheus.NewDesc(
		"consul_up",
		"Was the last query to Consul successful.",
		nil, nil,
	)
	invalidChars = regexp.MustCompile("[^a-zA-Z0-9_:]")
)

type ConsulCollector struct {
}

// Implements prometheus.Collector.
func (c ConsulCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

// Implements prometheus.Collector.
func (c ConsulCollector) Collect(ch chan<- prometheus.Metric) {
	config := api.DefaultConfig()
	config.Address = "consul:8500"
	consul, err := api.NewClient(config)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		return
	}

	metrics, err := consul.Agent().Metrics()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		return
	}
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)

	for _, gauge := range metrics.Gauges {
		name := invalidChars.ReplaceAllLiteralString(gauge.Name, "_")
		desc := prometheus.NewDesc(name, "Consul metric "+gauge.Name, nil, gauge.Labels)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(gauge.Value))
	}

	for _, counter := range metrics.Counters {
		name := invalidChars.ReplaceAllLiteralString(counter.Name, "_")
		desc := prometheus.NewDesc(name+"_total", "Consul metric "+counter.Name, nil, counter.Labels)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, float64(counter.Count))
	}

	for _, sample := range metrics.Samples {
		name := invalidChars.ReplaceAllLiteralString(sample.Name, "_") + "_seconds"
		countDesc := prometheus.NewDesc(name+"_count", "Consul metric "+sample.Name, nil, sample.Labels)
		ch <- prometheus.MustNewConstMetric(countDesc, prometheus.CounterValue, float64(sample.Count))
		sumDesc := prometheus.NewDesc(name+"_sum", "Consul metric "+sample.Name, nil, sample.Labels)
		ch <- prometheus.MustNewConstMetric(sumDesc, prometheus.CounterValue, sample.Sum/1000)
	}
}

func main() {
	consulCollector := ConsulCollector{}
	prometheus.MustRegister(consulCollector)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8000", nil))
}
