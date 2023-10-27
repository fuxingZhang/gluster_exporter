package collector

import (
	"gluster_exporter/pkg/conf"
	"gluster_exporter/pkg/gluster"
	"gluster_exporter/pkg/logger"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(NewVolumeHealCollector(), true)
}

func NewVolumeHealCollector() Collector {
	const subsystem = "volume"

	var volumeHealLabels = []string{"cluster_id", "volume", "brick_path", "host"}

	return &VolumeHealCollector{
		volumeHealCountVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "heal_count",
			Help:      "self heal count for volume",
		}, volumeHealLabels),
		volumeSplitBrainHealCountVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "split_brain_heal_count",
			Help:      "self heal count for volume in split brain",
		}, volumeHealLabels),
	}
}

type VolumeHealCollector struct {
	volumeHealCountVec           *prometheus.GaugeVec
	volumeSplitBrainHealCountVec *prometheus.GaugeVec
}

func (c *VolumeHealCollector) Name() string {
	return "volume_heal"
}

func (c *VolumeHealCollector) Collect(ch chan<- prometheus.Metric) (err error) {
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

	// locHealInfoFunc is a function literal, which takes
	// arg1: f1 a function which takes a string and returns ([]HealEntry, error)
	// (can be 'HealInfo' or 'SplitBrainHealInfo')
	// arg2: gVect a pointer to GaugeVec
	// (can be either 'glusterVolumeHealCount' or 'glusterVolumeSplitBrainHealCount')
	// arg3: volName a string representing the volume name
	// arg4: errStr the error string in case of error
	locHealInfoFunc := func(f1 func(string) ([]gluster.HealEntry, error), gaugeVec *prometheus.GaugeVec, volName string, errStr string) {
		// Get the heal count
		heals, err := f1(volName)
		if err != nil {
			logger.Error("msg", errStr, "err", err, "volume", volName)
			return
		}
		for _, healinfo := range heals {
			labels := prometheus.Labels{
				"cluster_id": conf.ClusterID(),
				"volume":     volName,
				"brick_path": healinfo.Brick,
				"host":       healinfo.Hostname,
			}
			collectGaugeVec(gaugeVec, labels, float64(healinfo.NumHealEntries), ch)
		}
	}

	for _, volume := range volumes {
		name := volume.Name
		if strings.Contains(volume.Type, "Replicate") {
			locHealInfoFunc(
				gluster.Gluster().HealInfo,
				c.volumeHealCountVec,
				name,
				"Error getting heal info",
			)
			locHealInfoFunc(
				gluster.Gluster().SplitBrainHealInfo,
				c.volumeSplitBrainHealCountVec,
				name,
				"Error getting split brain heal info",
			)
		}
		// volume.Type, "Distribute" error "exit status 234" volume=myvolume
		// if strings.Contains(volume.Type, "Disperse") || strings.Contains(volume.Type, "Distribute") {
		if strings.Contains(volume.Type, "Disperse") {
			locHealInfoFunc(
				gluster.Gluster().HealInfo,
				c.volumeHealCountVec,
				name,
				"Error getting heal info",
			)
		}
	}
	return
}
