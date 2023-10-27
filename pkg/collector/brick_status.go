package collector

import (
	"gluster_exporter/pkg/conf"
	"gluster_exporter/pkg/consts"
	"gluster_exporter/pkg/gluster"
	"gluster_exporter/pkg/logger"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(NewBrickStatusCollector(), true)
}

func NewBrickStatusCollector() Collector {
	var labels = []string{
		"cluster_id",
		"volume",
		"hostname",
		"brick_path",
		"peer_id",
		"pid",
	}

	return &BrickStatusCollector{
		brickUpVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "brick_up",
			Help:      "Brick up (1-up, 0-down)",
		}, labels),
	}
}

type BrickStatusCollector struct {
	brickUpVec *prometheus.GaugeVec
}

func (c *BrickStatusCollector) Name() string {
	return "brick_status"
}

func (c *BrickStatusCollector) Collect(ch chan<- prometheus.Metric) (err error) {
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

	for _, volume := range volumes {
		// If volume is down, the bricks should be marked down
		var brickStatus []gluster.BrickStatus
		if volume.State != consts.VolumeStateStarted {
			for _, subvol := range volume.SubVolumes {
				for _, brick := range subvol.Bricks {
					status := gluster.BrickStatus{
						Hostname: brick.Host,
						PeerID:   brick.PeerID,
						Status:   0,
						PID:      0,
						Path:     brick.Path,
						Volume:   volume.Name,
					}
					brickStatus = append(brickStatus, status)
				}
			}
		} else {
			brickStatus, err = gluster.Gluster().VolumeBrickStatus(volume.Name)
			if err != nil {
				logger.Error("volume", volume.Name, "err", err, "msg", "Error getting bricks status")
				continue
			}
		}

		for _, entry := range brickStatus {
			labels := prometheus.Labels{
				"cluster_id": conf.ClusterID(),
				"volume":     volume.Name,
				"hostname":   entry.Hostname,
				"brick_path": entry.Path,
				"peer_id":    entry.PeerID,
				"pid":        strconv.Itoa(entry.PID),
			}
			collectGaugeVec(c.brickUpVec, labels, float64(entry.Status), ch)
		}
	}

	return
}
