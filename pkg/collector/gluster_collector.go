package collector

import (
	"gluster_exporter/pkg/logger"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "gluster" // For Prometheus metrics.
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"node_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_success"),
		"node_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)

// GlusterCollector implements the prometheus.Collector interface.
type GlusterCollector struct {
	collectors []Collector
}

// NewGlusterCollector creates a new GlusterCollector.
func NewGlusterCollector() *GlusterCollector {
	var collectors = []Collector{}

	for v, enable := range Collectors {
		if *enable {
			collectors = append(collectors, v)
		}
	}

	return &GlusterCollector{
		collectors: collectors,
	}
}

// Describe implements the prometheus.Collector interface.
func (gc GlusterCollector) Describe(ch chan<- *prometheus.Desc) {
	// prometheus.DescribeByCollect(gc, ch)
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

// Collect implements the prometheus.Collector interface.
func (gc GlusterCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(gc.collectors))
	for _, c := range gc.collectors {
		go func(c Collector) {
			execute(c, ch)
			wg.Done()
		}(c)
	}
	wg.Wait()
}

func execute(c Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Collect(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		logger.Error("msg", "collector failed", "name", c.Name(), "duration_seconds", duration.Seconds(), "err", err)
		success = 0
	} else {
		logger.Debug("msg", "collector succeeded", "name", c.Name(), "duration_seconds", duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), c.Name())
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, c.Name())
}
