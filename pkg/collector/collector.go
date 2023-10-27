package collector

import (
	"fmt"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// Collectors Collectors
var Collectors = map[Collector]*bool{}

// Collector is the interface a collector has to implement.
type Collector interface {
	Name() string
	// Get new metrics and expose them via prometheus registry.
	Collect(ch chan<- prometheus.Metric) error
}

func registerCollector(collector Collector, isDefaultEnabled bool) {
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("collector.%s", collector.Name())
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector.Name(), helpDefaultState)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)

	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Bool()
	Collectors[collector] = flag
}

func collectGaugeVec(gaugeVec *prometheus.GaugeVec, labels prometheus.Labels, val float64, ch chan<- prometheus.Metric) {
	metric := gaugeVec.With(labels)
	metric.Set(val)
	metric.Collect(ch)
}
