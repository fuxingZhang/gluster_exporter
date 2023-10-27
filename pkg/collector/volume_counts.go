package collector

import (
	"gluster_exporter/pkg/conf"
	"gluster_exporter/pkg/consts"
	"gluster_exporter/pkg/gluster"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(NewVolumeCountsCollector(), true)
}

func NewVolumeCountsCollector() Collector {
	const subsystem = "volume"

	var (
		volumeLabels   = []string{"cluster_id", "volume"}
		clusterIDLabel = []string{"cluster_id"}
	)

	return &VolumeCountsCollector{
		volumeTotalCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "total_count"),
			"Total no of volumes",
			clusterIDLabel,
			nil,
		),
		volumeCreatedCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "created_count"),
			"Freshly created no of volumes",
			clusterIDLabel,
			nil,
		),
		volumeStartedCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "started_count"),
			"Total no of started volumes",
			clusterIDLabel,
			nil,
		),
		volumeBrickCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "brick_count"),
			"Total no of bricks in volume",
			volumeLabels,
			nil,
		),
		volumeSnapshotBrickCountTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "snapshot_brick_count_total"),
			"Total count of snapshots bricks for volume",
			volumeLabels,
			nil,
		),
		volumeSnapshotBrickCountActive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "snapshot_brick_count_active"),
			"Total active count of snapshots bricks for volume",
			volumeLabels,
			nil,
		),
		volumeUp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"Volume is started or not (1-started, 0-not started)",
			volumeLabels,
			nil,
		),
	}
}

type VolumeCountsCollector struct {
	volumeTotalCount               *prometheus.Desc
	volumeCreatedCount             *prometheus.Desc
	volumeStartedCount             *prometheus.Desc
	volumeBrickCount               *prometheus.Desc
	volumeSnapshotBrickCountTotal  *prometheus.Desc
	volumeSnapshotBrickCountActive *prometheus.Desc
	volumeUp                       *prometheus.Desc
}

func (c *VolumeCountsCollector) Name() string {
	return "volume_counts"
}

func (c *VolumeCountsCollector) Collect(ch chan<- prometheus.Metric) (err error) {
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
	snapshots, err := gluster.Gluster().Snapshots()
	if err != nil {
		return
	}

	var volCount, volStartCount, volCreatedCount int

	volCount = len(volumes)
	for _, volume := range volumes {
		up := 0
		switch volume.State {
		case consts.VolumeStateStarted:
			up = 1
			volStartCount++
		case consts.VolumeStateCreated:
			volCreatedCount++
		default:
			// Volume is stopped, nothing to do as the stopped count
			// could be derived using total - started - created
		}
		volumeLabels := []string{conf.ClusterID(), volume.Name}
		ch <- prometheus.MustNewConstMetric(
			c.volumeUp,
			prometheus.GaugeValue,
			float64(up),
			volumeLabels...,
		)
		volBrickCount := 0
		for _, subvol := range volume.SubVolumes {
			volBrickCount += len(subvol.Bricks)
		}
		ch <- prometheus.MustNewConstMetric(
			c.volumeBrickCount,
			prometheus.GaugeValue,
			float64(volBrickCount),
			volumeLabels...,
		)
		volSnapBrickCountTotal := 0
		volSnapBrickCountActive := 0
		for _, snap := range snapshots {
			if volume.Name == snap.VolumeName {
				volSnapBrickCountTotal += volBrickCount
				if snap.Started {
					volSnapBrickCountActive += volBrickCount
				}
			}
		}
		ch <- prometheus.MustNewConstMetric(
			c.volumeSnapshotBrickCountTotal,
			prometheus.GaugeValue,
			float64(volSnapBrickCountTotal),
			volumeLabels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.volumeSnapshotBrickCountActive,
			prometheus.GaugeValue,
			float64(volSnapBrickCountActive),
			volumeLabels...,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		c.volumeTotalCount,
		prometheus.GaugeValue,
		float64(volCount),
		conf.ClusterID(),
	)
	ch <- prometheus.MustNewConstMetric(
		c.volumeStartedCount,
		prometheus.GaugeValue,
		float64(volStartCount),
		conf.ClusterID(),
	)
	ch <- prometheus.MustNewConstMetric(
		c.volumeCreatedCount,
		prometheus.GaugeValue,
		float64(volCreatedCount),
		conf.ClusterID(),
	)

	return
}
