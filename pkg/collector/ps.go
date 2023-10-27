package collector

import (
	"gluster_exporter/pkg/conf"
	"gluster_exporter/pkg/gluster"
	"gluster_exporter/pkg/logger"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(NewPsCollector(), true)
}

func NewPsCollector() Collector {
	const subsystem = ""

	var labels = []string{
		"cluster_id",
		"volume",
		"peer_id",
		"brick_path",
		"name",
	}

	return &PsCollector{
		cpuPercentageVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "cpu_percentage",
			Help:      "CPU Percentage used by Gluster processes",
		}, labels),
		memoryPercentageVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "memory_percentage",
			Help:      "Memory Percentage used by Gluster processes",
		}, labels),
		residentMemoryVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "resident_memory_bytes",
			Help:      "Resident Memory of Gluster processes in bytes",
		}, labels),
		virtualMemoryVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "virtual_memory_bytes",
			Help:      "Virtual Memory of Gluster processes in bytes",
		}, labels),
		elapsedTimeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "elapsed_time_seconds",
			Help:      "Elapsed Time of Gluster processes in seconds",
		}, labels),
	}
}

type PsCollector struct {
	cpuPercentageVec    *prometheus.GaugeVec
	memoryPercentageVec *prometheus.GaugeVec
	residentMemoryVec   *prometheus.GaugeVec
	virtualMemoryVec    *prometheus.GaugeVec
	elapsedTimeVec      *prometheus.GaugeVec
}

func (c *PsCollector) Name() string {
	return "ps"
}

func (c *PsCollector) Collect(ch chan<- prometheus.Metric) (err error) {
	var glusterProcs = []string{
		"glusterd",
		"glusterfsd",
		"glusterd2",
		// TODO: Add more processes
	}

	args := []string{
		"--no-header", // No header in the output
		"-ww",         // To set unlimited width to avoid crop
		"-o",          // Output Format
		"pid,pcpu,pmem,rsz,vsz,etimes,comm",
		"-C",
		strings.Join(glusterProcs, ","),
	}

	out, err := exec.Command("ps", args...).Output() // #nosec

	if err != nil {
		// Return without exporting metrics in this cycle
		return err
	}

	peerID, err := gluster.Gluster().LocalPeerID()
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(out), "\n") {
		// Sample data:
		// 6959  0.0  0.6 12840 713660  504076 glusterfs
		lineDataTmp := strings.Split(line, " ")
		lineData := []string{}
		for _, d := range lineDataTmp {
			if strings.Trim(d, " ") == "" {
				continue
			}
			lineData = append(lineData, d)
		}

		if len(lineData) < 7 {
			continue
		}
		cmdlineArgs, err := getCmdLine(lineData[0])
		if err != nil {
			logger.Error(
				"msg", "Error getting command line",
				"command", lineData[6],
				"pid", lineData[0],
				"err", err,
			)
			continue
		}

		if len(cmdlineArgs) == 0 {
			// No cmdline file, may be that process died
			continue
		}

		var lbls prometheus.Labels
		switch lineData[6] {
		case "glusterd", "glusterd2":
			lbls = getGlusterdLabels(peerID, lineData[6])
		case "glusterfsd":
			lbls = getGlusterFsdLabels(peerID, lineData[6], cmdlineArgs)
		default:
			lbls = getUnknownLabels(peerID, lineData[6])
		}

		pcpu, err := strconv.ParseFloat(lineData[1], 64)
		if err != nil {
			logger.Error(
				"msg", "Error getting command line",
				"value", lineData[1],
				"command", lineData[6],
				"pid", lineData[0],
				"err", err,
			)
			continue
		}

		pmem, err := strconv.ParseFloat(lineData[2], 64)
		if err != nil {
			logger.Error(
				"msg", "Error getting command line",
				"value", lineData[2],
				"command", lineData[6],
				"pid", lineData[0],
				"err", err,
			)
			continue
		}
		rsz, err := strconv.ParseFloat(lineData[3], 64)
		if err != nil {
			logger.Error(
				"msg", "Error getting command line",
				"value", lineData[3],
				"command", lineData[6],
				"pid", lineData[0],
				"err", err,
			)
			continue
		}

		vsz, err := strconv.ParseFloat(lineData[4], 64)
		if err != nil {
			logger.Error(
				"msg", "Unable to parse vsz value",
				"value", lineData[4],
				"command", lineData[6],
				"pid", lineData[0],
				"err", err,
			)
			continue
		}

		etimes, err := strconv.ParseFloat(lineData[5], 64)
		if err != nil {
			logger.Error(
				"msg", "Unable to parse etimes value",
				"value", lineData[5],
				"command", lineData[6],
				"pid", lineData[0],
				"err", err,
			)
			continue
		}

		// Convert to bytes from kilo bytes
		vsz = vsz * 1024
		rsz = rsz * 1024

		// Update the Metrics
		collectGaugeVec(c.cpuPercentageVec, lbls, pcpu, ch)
		collectGaugeVec(c.memoryPercentageVec, lbls, pmem, ch)
		collectGaugeVec(c.residentMemoryVec, lbls, rsz, ch)
		collectGaugeVec(c.virtualMemoryVec, lbls, vsz, ch)
		collectGaugeVec(c.elapsedTimeVec, lbls, etimes, ch)
	}
	return
}

func getCmdLine(pid string) ([]string, error) {
	var args []string

	out, err := os.ReadFile(filepath.Clean("/proc/" + pid + "/cmdline"))
	if err != nil {
		return args, err
	}

	return strings.Split(strings.Trim(string(out), "\x00"), "\x00"), nil
}

func getGlusterdLabels(peerID, cmd string) prometheus.Labels {
	return prometheus.Labels{
		"cluster_id": conf.ClusterID(),
		"name":       cmd,
		"volume":     "",
		"peer_id":    peerID,
		"brick_path": "",
	}
}

func getGlusterFsdLabels(peerID, cmd string, args []string) prometheus.Labels {
	bpath := ""
	volume := ""
	prevArg := ""

	for _, a := range args {
		if prevArg == "--brick-name" {
			bpath = a
		} else if prevArg == "--volfile-id" {
			volume = strings.Split(a, ".")[0]
		}
		prevArg = a
	}

	return prometheus.Labels{
		"cluster_id": conf.ClusterID(),
		"name":       cmd,
		"volume":     volume,
		"peer_id":    peerID,
		"brick_path": bpath,
	}
}

func getUnknownLabels(peerID, cmd string) prometheus.Labels {
	return prometheus.Labels{
		"cluster_id": conf.ClusterID(),
		"name":       cmd,
		"volume":     "",
		"peer_id":    peerID,
		"brick_path": "",
	}
}
