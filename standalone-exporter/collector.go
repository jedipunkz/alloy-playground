package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// loadCollector exposes /proc/loadavg as tiny_load_average{period="1m|5m|15m"}.
type loadCollector struct {
	path string
	desc *prometheus.Desc
}

func newLoadCollector(path string) *loadCollector {
	return &loadCollector{
		path: path,
		desc: prometheus.NewDesc("tiny_load_average", "System load average.", []string{"period"}, nil),
	}
}

func (c *loadCollector) Describe(ch chan<- *prometheus.Desc) { ch <- c.desc }

func (c *loadCollector) Collect(ch chan<- prometheus.Metric) {
	b, err := os.ReadFile(c.path)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.desc, err)
		return
	}
	loads, err := parseLoadavg(string(b))
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.desc, err)
		return
	}
	for i, period := range []string{"1m", "5m", "15m"} {
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, loads[i], period)
	}
}

// parseLoadavg parses the first three fields of /proc/loadavg, e.g. "0.52 0.58 0.59 1/389 12345".
func parseLoadavg(s string) ([3]float64, error) {
	var loads [3]float64
	f := strings.Fields(s)
	if len(f) < 3 {
		return loads, fmt.Errorf("unexpected loadavg format: %q", s)
	}
	for i := range loads {
		v, err := strconv.ParseFloat(f[i], 64)
		if err != nil {
			return loads, err
		}
		loads[i] = v
	}
	return loads, nil
}
