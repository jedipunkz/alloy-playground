// Package tiny is copied into the Alloy source tree as
// internal/component/prometheus/exporter/tiny and registers
// the prometheus.exporter.tiny component.
package tiny

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/grafana/alloy/internal/component"
	"github.com/grafana/alloy/internal/component/prometheus/exporter"
	"github.com/grafana/alloy/internal/component/prometheus/exporter/common"
	"github.com/grafana/alloy/internal/featuregate"
	"github.com/grafana/alloy/internal/static/integrations"
)

func init() {
	component.Register(component.Registration{
		Name:      "prometheus.exporter.tiny",
		Stability: featuregate.StabilityExperimental,
		Args:      Arguments{},
		Exports:   exporter.Exports{},

		Build: exporter.New(createExporter, "tiny"),
	})
}

// Arguments configures the prometheus.exporter.tiny component.
type Arguments struct {
	LoadavgPath string `alloy:"loadavg_path,attr,optional"`
}

// DefaultArguments holds the default arguments for prometheus.exporter.tiny.
var DefaultArguments = Arguments{LoadavgPath: "/proc/loadavg"}

// SetToDefault implements syntax.Defaulter.
func (a *Arguments) SetToDefault() { *a = DefaultArguments }

func createExporter(opts component.Options, args component.Arguments) (integrations.Integration, string, error) {
	a := args.(Arguments)
	i := integrations.NewCollectorIntegration("tiny",
		integrations.WithLogger(opts.Logger),
		integrations.WithCollectors(newLoadCollector(a.LoadavgPath)),
	)
	// loadavg is host-level, so use the hostname as the instance label.
	return i, common.HostNameInstanceKey(), nil
}

// loadCollector is the same collector as standalone-exporter/collector.go.
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
