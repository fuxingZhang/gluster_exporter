package main

import (
	"net/http"
	"os"
	"regexp"
	"runtime"

	"gluster_exporter/pkg/collector"
	"gluster_exporter/pkg/conf"
	"gluster_exporter/pkg/gluster"
	"gluster_exporter/pkg/logger"

	stdlog "log"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

const exporterName = "gluster_exporter"

func main() {
	var (
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		includeExporterMetrics = kingpin.Flag(
			"web.include-exporter-metrics",
			"Include metrics about the exporter itself (promhttp_*, process_*, go_*) (default: false)",
		).Bool()
		maxRequests = kingpin.Flag(
			"web.max-requests",
			"Maximum number of parallel scrape requests. Use 0 to disable.",
		).Default("5").Int()
		maxProcs = kingpin.Flag(
			"runtime.gomaxprocs", "The target number of CPUs Go will run on (GOMAXPROCS)",
		).Envar("GOMAXPROCS").Default("1").Int()
		toolkitFlags      = kingpinflag.AddFlags(kingpin.CommandLine, ":9713")
		glusterd2Endpoint = kingpin.Flag(
			"gd2-rest-endpoint",
			"gd2-rest-endpoint",
		).Default("http://localhost:24007").String()
		glusterRemoteHost = kingpin.Flag(
			"gd1-remote-host",
			"Connect to a remote gd1 host. The following collectors won't work in remote mode : gluster_volume_counts, gluster_volume_profile",
		).String()
		glusterGlusterdSock = kingpin.Flag(
			"gd1-glusterd-sock",
			"gd1-glusterd-sock",
		).String()
		glusterMgmt = kingpin.Flag(
			"gluster-mgmt",
			"gluster-mgmt",
		).Default("glusterd").String()
		glusterCmd = kingpin.Flag(
			"gluster-binary-path",
			"gluster-binary-path",
		).Default("gluster").String()
		glusterdDir = kingpin.Flag(
			"glusterd-dir",
			"glusterd workdir",
		).Default("/var/lib/glusterd").String()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print(exporterName))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger.Init(promlogConfig)

	logger.Info("msg", "Starting gluster_exporter", "version", version.Info())
	logger.Info("msg", "Build context", "build_context", version.BuildContext())

	runtime.GOMAXPROCS(*maxProcs)
	logger.Debug("msg", "Go MAXPROCS", "procs", runtime.GOMAXPROCS(0))

	gluster.Init(conf.Init(&conf.ConfigCmdArgs{
		GlusterMgmt:         *glusterMgmt,
		Glusterd2Endpoint:   *glusterd2Endpoint,
		GlusterCmd:          *glusterCmd,
		GlusterRemoteHost:   *glusterRemoteHost,
		GlusterGlusterdSock: *glusterGlusterdSock,
		GlusterdWorkdir:     *glusterdDir,
	}))

	http.Handle(*metricsPath, handler(*includeExporterMetrics, *maxRequests))

	if *metricsPath != "/" {
		landingConfig := web.LandingConfig{
			Name:        "Gluster Exporter",
			Description: "Prometheus Gluster Exporter",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
					Text:    "Metrics",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			logger.Error("err", err)
			os.Exit(1)
		}
		http.Handle("/", landingPage)
	}

	server := &http.Server{}
	if err := web.ListenAndServe(server, toolkitFlags, logger.Logger()); err != nil {
		logger.Error("err", err)
		os.Exit(1)
	}
}

func handler(includeExporterMetrics bool, maxRequests int) http.Handler {
	registry := prometheus.NewRegistry()

	registry.MustRegister(version.NewCollector(exporterName))

	registry.MustRegister(collectors.NewBuildInfoCollector())

	if includeExporterMetrics {
		registry.MustRegister(
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(
				collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
			),
		)
	}

	c := collector.NewGlusterCollector()

	registry.MustRegister(c)

	handler := promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			ErrorLog:            stdlog.New(log.NewStdlibAdapter(level.Error(logger.Logger())), "", 0),
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: maxRequests,
			// promhttp_metric_handler_errors_total
			Registry: registry,
		},
	)

	if includeExporterMetrics {
		handler = promhttp.InstrumentMetricHandler(
			registry, handler,
		)
	}

	return handler
}
