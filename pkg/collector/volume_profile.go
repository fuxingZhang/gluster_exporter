package collector

import (
	"errors"
	"gluster_exporter/pkg/conf"
	"gluster_exporter/pkg/consts"
	"gluster_exporter/pkg/gluster"
	"gluster_exporter/pkg/logger"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(NewVolumeProfileCollector(), true)
}

func NewVolumeProfileCollector() Collector {
	const subsystem = "volume_profile"

	var (
		volumeProfileInfoLabels    = []string{"cluster_id", "volume", "brick"}
		volumeProfileFopInfoLabels = []string{
			"cluster_id",
			"volume",
			"brick",
			"host",
			"fop",
		}
	)

	return &VolumeProfileCollector{
		volumeProfileTotalReadsVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "total_reads",
			Help:      "Total no of reads",
		}, volumeProfileInfoLabels),
		volumeProfileTotalWritesVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "total_writes",
			Help:      "Total no of writes",
		}, volumeProfileInfoLabels),
		volumeProfileDurationVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "duration_secs",
			Help:      "Duration",
		}, volumeProfileInfoLabels),
		volumeProfileTotalReadsIntVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "total_reads_interval",
			Help:      "Total no of reads for interval stats",
		}, volumeProfileInfoLabels),
		volumeProfileTotalWritesIntVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "total_writes_interval",
			Help:      "Total no of writes for interval stats",
		}, volumeProfileInfoLabels),
		volumeProfileDurationIntVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "duration_secs_interval",
			Help:      "Duration for interval stats",
		}, volumeProfileInfoLabels),
		volumeProfileFopHitsVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "fop_hits",
			Help:      "Cumulative FOP hits",
		}, volumeProfileFopInfoLabels),
		volumeProfileFopAvgLatencyVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "fop_avg_latency",
			Help:      "Cumulative FOP avergae latency",
		}, volumeProfileFopInfoLabels),
		volumeProfileFopMinLatencyVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "fop_min_latency",
			Help:      "Cumulative FOP min latency",
		}, volumeProfileFopInfoLabels),
		volumeProfileFopMaxLatencyVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "fop_max_latency",
			Help:      "Cumulative FOP max latency",
		}, volumeProfileFopInfoLabels),
		volumeProfileFopHitsIntVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "fop_hits_interval",
			Help:      "Interval based FOP hits",
		}, volumeProfileFopInfoLabels),
		volumeProfileFopAvgLatencyIntVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "fop_avg_latency_interval",
			Help:      "Interval based FOP average latency",
		}, volumeProfileFopInfoLabels),
		volumeProfileFopMinLatencyIntVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "fop_min_latency_interval",
			Help:      "Interval based FOP min latency",
		}, volumeProfileFopInfoLabels),
		volumeProfileFopMaxLatencyIntVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "fop_max_latency_interval",
			Help:      "Interval based FOP max latency",
		}, volumeProfileFopInfoLabels),
		volumeProfileFopTotalHitsAggregatedOpsVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "fop_total_hits_on_aggregated_fops",
			Help:      "Cumulative total hits on aggregated FOPs like READ_WRIET_OPS, LOCK_OPS, INODE_OPS etc",
		}, volumeProfileFopInfoLabels),
		volumeProfileFopTotalHitsAggregatedOpsIntVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "fop_total_hits_on_aggregated_fops_interval",
			Help:      "Interval based total hits on aggregated FOPs like READ_WRIET_OPS, LOCK_OPS, INODE_OPS etc",
		}, volumeProfileFopInfoLabels),
	}
}

type VolumeProfileCollector struct {
	volumeProfileTotalReadsVec                   *prometheus.GaugeVec
	volumeProfileTotalWritesVec                  *prometheus.GaugeVec
	volumeProfileDurationVec                     *prometheus.GaugeVec
	volumeProfileTotalReadsIntVec                *prometheus.GaugeVec
	volumeProfileTotalWritesIntVec               *prometheus.GaugeVec
	volumeProfileDurationIntVec                  *prometheus.GaugeVec
	volumeProfileFopTotalHitsAggregatedOpsVec    *prometheus.GaugeVec
	volumeProfileFopTotalHitsAggregatedOpsIntVec *prometheus.GaugeVec
	volumeProfileFopHitsVec                      *prometheus.GaugeVec
	volumeProfileFopAvgLatencyVec                *prometheus.GaugeVec
	volumeProfileFopMinLatencyVec                *prometheus.GaugeVec
	volumeProfileFopMaxLatencyVec                *prometheus.GaugeVec
	volumeProfileFopHitsIntVec                   *prometheus.GaugeVec
	volumeProfileFopAvgLatencyIntVec             *prometheus.GaugeVec
	volumeProfileFopMinLatencyIntVec             *prometheus.GaugeVec
	volumeProfileFopMaxLatencyIntVec             *prometheus.GaugeVec
}

func (c *VolumeProfileCollector) Name() string {
	return "volume_profile"
}

func (c *VolumeProfileCollector) Collect(ch chan<- prometheus.Metric) (err error) {
	isLeader, err := gluster.Gluster().IsLeader()
	if err != nil {
		return
	}
	if !isLeader {
		return
	}

	volumes, err := gluster.Gluster().VolumeInfo()
	if err != nil {
		return
	}
	volOption := consts.CountFOPHitsGD1

	if conf.Conf().GlusterMgmt == consts.MgmtGlusterd2 {
		volOption = consts.CountFOPHitsGD2
	}
	var (
		// supported aggregated operations are,
		// READ_WRITE_OPS, LOCK_OPS, ENTRY_OPS, INODE_OPS
		aggregatedOps = []opType{
			newReadWriteOpType(),
			newLockOpType(),
			newEntryOpType(),
			newINodeOpType(),
		}
	)
	for _, volume := range volumes {
		err := gluster.Gluster().EnableVolumeProfiling(volume)
		if err != nil {
			logger.Error(
				"msg", "Error enabling profiling for volume",
				"volume", volume.Name,
				"err", err,
			)
			continue
		}
		if value, exists := volume.Options[volOption]; !exists || value == "off" {
			// Volume profiling is explicitly switched off for volume, dont collect profile metrics
			continue
		}
		name := volume.Name
		profileinfo, err := gluster.Gluster().VolumeProfileInfo(name)
		if err != nil {
			logger.Error(
				"msg", "Error getting profile info",
				"volume", name,
				"err", err,
			)
			continue
		}
		for _, entry := range profileinfo {
			labels := prometheus.Labels{
				"cluster_id": conf.ClusterID(),
				"volume":     name,
				"brick":      entry.BrickName,
			}

			collectGaugeVec(c.volumeProfileTotalReadsVec, labels, float64(entry.TotalReads), ch)
			collectGaugeVec(c.volumeProfileTotalWritesVec, labels, float64(entry.TotalWrites), ch)
			collectGaugeVec(c.volumeProfileDurationVec, labels, float64(entry.Duration), ch)
			collectGaugeVec(c.volumeProfileTotalReadsIntVec, labels, float64(entry.TotalReadsInt), ch)
			collectGaugeVec(c.volumeProfileTotalWritesIntVec, labels, float64(entry.TotalWritesInt), ch)
			collectGaugeVec(c.volumeProfileDurationIntVec, labels, float64(entry.DurationInt), ch)

			brickhost := getBrickHost(volume, entry.BrickName)
			for _, eachOp := range aggregatedOps {
				labels := getVolumeProfileFopInfoLabels(name, entry.BrickName, brickhost, eachOp.String())
				collectGaugeVec(c.volumeProfileFopTotalHitsAggregatedOpsVec, labels, eachOp.opHits(entry.FopStats), ch)
				collectGaugeVec(c.volumeProfileFopTotalHitsAggregatedOpsIntVec, labels, eachOp.opHits(entry.FopStatsInt), ch)
			}
			for _, fopInfo := range entry.FopStats {
				labels := getVolumeProfileFopInfoLabels(name, entry.BrickName, brickhost, fopInfo.Name)
				collectGaugeVec(c.volumeProfileFopHitsVec, labels, float64(fopInfo.Hits), ch)
				collectGaugeVec(c.volumeProfileFopAvgLatencyVec, labels, fopInfo.AvgLatency, ch)
				collectGaugeVec(c.volumeProfileFopMinLatencyVec, labels, fopInfo.MinLatency, ch)
				collectGaugeVec(c.volumeProfileFopMaxLatencyVec, labels, fopInfo.MaxLatency, ch)
			}
			for _, fopInfo := range entry.FopStatsInt {
				labels := getVolumeProfileFopInfoLabels(name, entry.BrickName, brickhost, fopInfo.Name)
				collectGaugeVec(c.volumeProfileFopHitsIntVec, labels, float64(fopInfo.Hits), ch)
				collectGaugeVec(c.volumeProfileFopAvgLatencyIntVec, labels, fopInfo.AvgLatency, ch)
				collectGaugeVec(c.volumeProfileFopMinLatencyIntVec, labels, fopInfo.MinLatency, ch)
				collectGaugeVec(c.volumeProfileFopMaxLatencyIntVec, labels, fopInfo.MaxLatency, ch)
			}
		}
	}

	return
}

func getVolumeProfileFopInfoLabels(volname string, brick string, host string, fop string) prometheus.Labels {
	return prometheus.Labels{
		"cluster_id": conf.ClusterID(),
		"volume":     volname,
		"brick":      brick,
		"host":       host,
		"fop":        fop,
	}
}

func getBrickHost(vol gluster.Volume, brickname string) string {
	hostid := strings.Split(brickname, ":")[0]
	for _, subvol := range vol.SubVolumes {
		for _, brick := range subvol.Bricks {
			if brick.PeerID == hostid {
				return brick.Host
			}
		}
	}
	return ""
}

// opType represents aggregated operations like
// READ_WRITE_OPS, INODE_OPS, ENTRY_OPS, LOCK_OPS etc...
type opType struct {
	opName       string
	opsSupported map[string]struct{}
}

// String method returns the name of the common 'opType'
// makes it compatible with 'Stringer' interface
func (ot opType) String() string {
	return ot.opName
}

// opSupported method checks whether the given operation is supported in this opType
func (ot opType) opSupported(opLabel string) bool {
	if _, ok := ot.opsSupported[opLabel]; ok {
		return true
	}
	return false
}

// opHits calculates total number of 'ot' type operations in a list of 'FopStat's
// and returns the total number of hits
func (ot opType) opHits(fopStats []gluster.FopStat) float64 {
	var totalOpHits float64 // default ZERO value is assigned
	for _, eachFopS := range fopStats {
		if ot.opSupported(eachFopS.Name) {
			totalOpHits += float64(eachFopS.Hits)
		}
	}
	return totalOpHits
}

// newOperationType creates a 'opType' object
func newOperationType(commonOpName string, supportedOpLabels []string) (opType, error) {
	commonOpName = strings.TrimSpace(commonOpName)
	var opT = opType{opName: commonOpName}
	if len(supportedOpLabels) == 0 {
		return opT, errors.New("supported operation labels should not be empty")
	} else if commonOpName == "" {
		return opT, errors.New("empty common operation name is not allowed")
	}
	opT.opsSupported = make(map[string]struct{})
	var emtS struct{}
	for _, opLabel := range supportedOpLabels {
		opT.opsSupported[opLabel] = emtS
	}
	return opT, nil
}

// newReadWriteOpType creates a 'READ_WRITE_OPS' type,
// which aggregates all the read/write operations
func newReadWriteOpType() opType {
	var opsSupported = []string{"CREATE", "DISCARD", "FALLOCATE", "FLUSH", "FSYNC",
		"FSYNCDIR", "RCHECKSUM", "READ", "READDIR", "READDIRP", "READY",
		"WRITE", "ZEROFILL",
	}
	var rwOT, _ = newOperationType("READ_WRITE_OPS", opsSupported)
	return rwOT
}

// newLockOpType creates a 'LOCK_OPS' type,
// which aggregates all the lock operations
func newLockOpType() opType {
	var opsSupported = []string{"ENTRYLK", "FENTRYLK", "FINODELK", "INODELK", "LK"}
	var lockOT, _ = newOperationType("LOCK_OPS", opsSupported)
	return lockOT
}

// newINodeOpType creates a 'INODE_OPS' type,
// which aggregates all the iNode associated operations
func newINodeOpType() opType {
	var opsSupported = []string{"ACCESS", "FGETXATTR", "FREMOVEXATTR", "FSETATTR",
		"FSETXATTR", "FSTAT", "FTRUNCATE", "FXATTROP", "GETXATTR", "LOOKUP", "OPEN",
		"OPENDIR", "READLINK", "REMOVEXATTR", "SEEK", "SETATTR", "SETXATTR", "STAT",
		"STATFS", "TRUNCATE", "XATTROP"}
	var iNodeOT, _ = newOperationType("INODE_OPS", opsSupported)
	return iNodeOT
}

// newEntryOpType creates a 'ENTRY_OPS' type,
// which aggregates all the file entry related operations
func newEntryOpType() opType {
	var opsSupported = []string{"LINK", "MKDIR", "MKNOD", "RENAME",
		"RMDIR", "SYMLINK", "UNLINK"}
	var entryOT, _ = newOperationType("ENTRY_OPS", opsSupported)
	return entryOT
}
