package collector

import (
	"gluster_exporter/pkg/conf"
	"gluster_exporter/pkg/consts"
	"gluster_exporter/pkg/gluster"
	"gluster_exporter/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(NewBrickCollector(), true)
}

func NewBrickCollector() Collector {
	const subsystem = ""

	var (
		brickLabels = []string{
			"cluster_id",
			"host",
			"id",
			"brick_path",
			"volume",
			"subvolume",
		}

		subvolLabels = []string{
			"cluster_id",
			"volume",
			"subvolume",
		}

		lvmLbls = []string{
			"cluster_id",
			"host",
			"id",
			"brick_path",
			"volume",
			"subvolume",
			"vg_name",
			"lv_path",
			"lv_uuid",
		}

		thinLvmLbls = []string{
			"cluster_id",
			"host",
			"thinpool_name",
			"vg_name",
			"volume",
			"subvolume",
			"brick_path",
		}
	)

	return &BrickCollector{
		brickCapacityUsedVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "brick_capacity_used_bytes",
			Help:      "Used capacity of gluster bricks in bytes",
		}, brickLabels),
		brickCapacityFreeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "brick_capacity_free_bytes",
			Help:      "Free capacity of gluster bricks in bytes",
		}, brickLabels),
		brickCapacityTotalVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "brick_capacity_bytes_total",
			Help:      "Total capacity of gluster bricks in bytes",
		}, brickLabels),
		brickInodesTotalVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "brick_inodes_total",
			Help:      "Total no of inodes of gluster brick disk",
		}, brickLabels),
		brickInodesFreeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "brick_inodes_free",
			Help:      "Free no of inodes of gluster brick disk",
		}, brickLabels),
		brickInodesUsedVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "brick_inodes_used",
			Help:      "Used no of inodes of gluster brick disk",
		}, brickLabels),
		subvolCapacityUsedVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "subvol_capacity_used_bytes",
			Help:      "Effective used capacity of gluster subvolume in bytes",
		}, subvolLabels),
		subvolCapacityTotalVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "subvol_capacity_total_bytes",
			Help:      "Effective total capacity of gluster subvolume in bytes",
		}, subvolLabels),
		brickLVSizeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "brick_lv_size_bytes",
			Help:      "Bricks LV size Bytes",
		}, lvmLbls),
		brickLVPercentVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "brick_lv_percent",
			Help:      "Bricks LV usage percent",
		}, lvmLbls),
		brickLVMetadataSizeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "brick_lv_metadata_size_bytes",
			Help:      "Bricks LV metadata size Bytes",
		}, lvmLbls),
		brickLVMetadataPercentVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "brick_lv_metadata_percent",
			Help:      "Bricks LV metadata usage percent",
		}, lvmLbls),
		glusterVGExtentTotalVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "vg_extent_total_count",
			Help:      "VG extent total count ",
		}, lvmLbls),
		glusterVGExtentAllocVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "vg_extent_alloc_count",
			Help:      "VG extent allocated count ",
		}, lvmLbls),
		glusterThinPoolDataTotalVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "thinpool_data_total_bytes",
			Help:      "Thin pool size Bytes",
		}, thinLvmLbls),
		glusterThinPoolDataUsedVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "thinpool_data_used_bytes",
			Help:      "Thin pool data used Bytes",
		}, thinLvmLbls),
		glusterThinPoolMetadataTotalVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "thinpool_metadata_total_bytes",
			Help:      "Thin pool metadata size Bytes",
		}, thinLvmLbls),
		glusterThinPoolMetadataUsedVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "thinpool_metadata_used_bytes",
			Help:      "Thin pool metadata used Bytes",
		}, thinLvmLbls),
	}
}

type BrickCollector struct {
	brickCapacityUsedVec            *prometheus.GaugeVec
	brickCapacityFreeVec            *prometheus.GaugeVec
	brickCapacityTotalVec           *prometheus.GaugeVec
	brickInodesTotalVec             *prometheus.GaugeVec
	brickInodesFreeVec              *prometheus.GaugeVec
	brickInodesUsedVec              *prometheus.GaugeVec
	subvolCapacityUsedVec           *prometheus.GaugeVec
	subvolCapacityTotalVec          *prometheus.GaugeVec
	brickLVSizeVec                  *prometheus.GaugeVec
	brickLVPercentVec               *prometheus.GaugeVec
	brickLVMetadataSizeVec          *prometheus.GaugeVec
	brickLVMetadataPercentVec       *prometheus.GaugeVec
	glusterVGExtentTotalVec         *prometheus.GaugeVec
	glusterVGExtentAllocVec         *prometheus.GaugeVec
	glusterThinPoolDataTotalVec     *prometheus.GaugeVec
	glusterThinPoolDataUsedVec      *prometheus.GaugeVec
	glusterThinPoolMetadataTotalVec *prometheus.GaugeVec
	glusterThinPoolMetadataUsedVec  *prometheus.GaugeVec
}

func (c *BrickCollector) Name() string {
	return "brick"
}

func (c *BrickCollector) Collect(ch chan<- prometheus.Metric) (err error) {
	volumes, err := gluster.Gluster().VolumeInfo()

	if err != nil {
		// Return without exporting metric in this cycle
		return err
	}

	localPeerID, err := gluster.Gluster().LocalPeerID()
	if err != nil {
		// Return without exporting metric in this cycle
		return err
	}

	for _, volume := range volumes {
		if volume.State != consts.VolumeStateStarted {
			// Export brick metrics only if the Volume
			// is is in Started state
			continue
		}
		subvols := volume.SubVolumes
		for _, subvol := range subvols {
			bricks := subvol.Bricks
			var maxBrickUsed float64
			var leastBrickTotal float64
			for _, brick := range bricks {
				if brick.PeerID == localPeerID {
					usage, err := gluster.DiskUsage(brick.Path)
					if err != nil {
						logger.Error(
							"msg", "Error getting disk usage",
							"volume", volume.Name,
							"brick_path", brick.Path,
							"err", err,
						)
						continue
					}
					var labels = getGlusterBrickLabels(brick, subvol.Name)
					// Update the metrics
					collectGaugeVec(c.brickCapacityUsedVec, labels, usage.Used, ch)
					collectGaugeVec(c.brickCapacityFreeVec, labels, usage.Free, ch)
					collectGaugeVec(c.brickCapacityTotalVec, labels, usage.All, ch)
					collectGaugeVec(c.brickInodesTotalVec, labels, usage.InodesAll, ch)
					collectGaugeVec(c.brickInodesFreeVec, labels, usage.InodesFree, ch)
					collectGaugeVec(c.brickInodesUsedVec, labels, usage.InodesUsed, ch)

					// Skip exporting utilization data in case of arbiter
					// brick to avoid wrong values when both the data bricks
					// are down
					if brick.Type != consts.BrickTypeArbiter && usage.Used >= maxBrickUsed {
						maxBrickUsed = usage.Used
					}
					if brick.Type != consts.BrickTypeArbiter {
						if leastBrickTotal == 0 || usage.All <= leastBrickTotal {
							leastBrickTotal = usage.All
						}
					}
					// Get lvm usage details
					stats, thinStats, err := gluster.LvmUsage(brick.Path)
					if err != nil {
						logger.Error(
							"msg", "Error getting lvm usage",
							"volume", volume.Name,
							"brick_path", brick.Path,
							"err", err,
						)
						continue
					}
					// Add metrics
					for _, stat := range stats {
						var labels = getGlusterLVMLabels(brick, subvol.Name, stat)
						// Convert to bytes
						collectGaugeVec(c.brickLVSizeVec, labels, stat.Size*1024*1024, ch)
						collectGaugeVec(c.brickLVPercentVec, labels, stat.DataPercent, ch)
						// Convert to bytes
						collectGaugeVec(c.brickLVMetadataSizeVec, labels, stat.MetadataSize*1024*1024, ch)
						collectGaugeVec(c.brickLVMetadataPercentVec, labels, stat.MetadataPercent, ch)
						collectGaugeVec(c.glusterVGExtentTotalVec, labels, stat.VGExtentTotal, ch)
						collectGaugeVec(c.glusterVGExtentAllocVec, labels, stat.VGExtentAlloc, ch)

					}
					for _, thinStat := range thinStats {
						var labels = getGlusterThinPoolLabels(brick, volume.Name, subvol.Name, thinStat)
						collectGaugeVec(c.glusterThinPoolDataTotalVec, labels, thinStat.ThinPoolDataTotal*1024*1024, ch)
						collectGaugeVec(c.glusterThinPoolDataUsedVec, labels, thinStat.ThinPoolDataUsed*1024*1024, ch)
						collectGaugeVec(c.glusterThinPoolMetadataTotalVec, labels, thinStat.ThinPoolMetadataTotal*1024*1024, ch)
						collectGaugeVec(c.glusterThinPoolMetadataUsedVec, labels, thinStat.ThinPoolMetadataUsed*1024*1024, ch)
					}
				}
			}
			effectiveCapacity := maxBrickUsed
			effectiveTotalCapacity := leastBrickTotal
			var subvolLabels = prometheus.Labels{
				"cluster_id": conf.ClusterID(),
				"volume":     volume.Name,
				"subvolume":  subvol.Name,
			}
			if subvol.Type == consts.SubvolTypeDisperse {
				// In disperse volume data bricks contribute to the sub
				// volume size
				effectiveCapacity = maxBrickUsed * float64(subvol.DisperseDataCount)
				effectiveTotalCapacity = leastBrickTotal * float64(subvol.DisperseDataCount)
			}

			// Export the metric only if available. it will be zero if the subvolume
			// contains only arbiter brick on current node or no local bricks on
			// this node
			if effectiveCapacity > 0 {
				collectGaugeVec(c.subvolCapacityUsedVec, subvolLabels, effectiveCapacity, ch)
			}
			if effectiveTotalCapacity > 0 {
				collectGaugeVec(c.subvolCapacityTotalVec, subvolLabels, effectiveTotalCapacity, ch)
			}
		}
	}
	return
}

func getGlusterBrickLabels(brick gluster.Brick, subvol string) prometheus.Labels {
	return prometheus.Labels{
		"cluster_id": conf.ClusterID(),
		"host":       brick.Host,
		"id":         brick.ID,
		"brick_path": brick.Path,
		"volume":     brick.VolumeName,
		"subvolume":  subvol,
	}
}

func getGlusterThinPoolLabels(brick gluster.Brick, vol string, subvol string, thinStat gluster.ThinPoolStat) prometheus.Labels {
	return prometheus.Labels{
		"cluster_id":    conf.ClusterID(),
		"host":          brick.Host,
		"thinpool_name": thinStat.ThinPoolName,
		"vg_name":       thinStat.ThinPoolVGName,
		"volume":        vol,
		"subvolume":     subvol,
		"brick_path":    brick.Path,
	}
}

func getGlusterLVMLabels(brick gluster.Brick, subvol string, stat gluster.LVMStat) prometheus.Labels {
	return prometheus.Labels{
		"cluster_id": conf.ClusterID(),
		"host":       brick.Host,
		"id":         brick.ID,
		"brick_path": brick.Path,
		"volume":     brick.VolumeName,
		"subvolume":  subvol,
		"vg_name":    stat.VGName,
		"lv_path":    stat.Path,
		"lv_uuid":    stat.UUID,
	}
}
