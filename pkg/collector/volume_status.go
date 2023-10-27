package collector

import (
	"gluster_exporter/pkg/gluster"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(NewVolumeStatusCollector(), true)
}

func NewVolumeStatusCollector() Collector {
	const subsystem = "volume"

	var (
		volStatusBrickCountLabels = []string{"instance", "volume_name"}
		volStatusPerBrickLabels   = []string{"instance", "volume_name", "hostname", "peer_id", "pid", "brick_path"}
	)

	return &VolumeStatusCollector{
		volStatusBrickCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "status_brick_count"),
			"Number of bricks for volume",
			volStatusBrickCountLabels,
			nil,
		),
		volumeBrickStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "brick_status"),
			"Per node brick status for volume",
			volStatusPerBrickLabels,
			nil,
		),
		volumeBrickPort: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "brick_port"),
			"Brick port",
			volStatusPerBrickLabels,
			nil,
		),
		volumeBrickPid: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "brick_pid"),
			"Brick pid",
			volStatusPerBrickLabels,
			nil,
		),
		volumeBrickTotalInodes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "brick_total_inodes"),
			"Brick total inodes",
			volStatusPerBrickLabels,
			nil,
		),
		volumeBrickFreeInodes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "brick_free_inodes"),
			"Brick free inodes",
			volStatusPerBrickLabels,
			nil,
		),
		volumeBrickTotalBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "brick_total_bytes"),
			"Brick total bytes",
			volStatusPerBrickLabels,
			nil,
		),
		volumeBrickFreeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "brick_free_bytes"),
			"Brick free bytes",
			volStatusPerBrickLabels,
			nil,
		),
	}
}

type VolumeStatusCollector struct {
	volStatusBrickCount    *prometheus.Desc
	volumeBrickStatus      *prometheus.Desc
	volumeBrickPort        *prometheus.Desc
	volumeBrickPid         *prometheus.Desc
	volumeBrickTotalInodes *prometheus.Desc
	volumeBrickFreeInodes  *prometheus.Desc
	volumeBrickTotalBytes  *prometheus.Desc
	volumeBrickFreeBytes   *prometheus.Desc
}

func (c *VolumeStatusCollector) Name() string {
	return "volume_status"
}

func (c *VolumeStatusCollector) Collect(ch chan<- prometheus.Metric) (err error) {
	peerID, err := gluster.Gluster().LocalPeerID()
	if err != nil {
		return
	}

	volumes, err := gluster.Gluster().VolumeStatus()
	if err != nil {
		return
	}

	// Get monitored gluster instance FQDN
	peers, err := gluster.Gluster().Peers()
	if err != nil {
		return
	}
	fqdn := "n/a"
	for _, peer := range peers {
		if peer.ID == peerID {
			// TODO: figure out which value of PeerAddresses may
			// be hostname -- or resolve ip ourselves
			fqdn = peer.PeerAddresses[0]
			break
		}
	}

	for _, vol := range volumes {
		ch <- prometheus.MustNewConstMetric(
			c.volStatusBrickCount,
			prometheus.GaugeValue,
			float64(len(vol.Nodes)),
			fqdn,
			vol.Name,
		)

		for _, node := range vol.Nodes {
			brickPid := strconv.Itoa(node.PID)

			perBrickLabels := []string{
				fqdn,
				vol.Name,
				node.Hostname,
				node.PeerID,
				brickPid,
				node.Path,
			}
			ch <- prometheus.MustNewConstMetric(
				c.volumeBrickStatus,
				prometheus.GaugeValue,
				float64(node.Status),
				perBrickLabels...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.volumeBrickPort,
				prometheus.GaugeValue,
				float64(node.Port),
				perBrickLabels...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.volumeBrickPid,
				prometheus.GaugeValue,
				float64(node.PID),
				perBrickLabels...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.volumeBrickTotalInodes,
				prometheus.GaugeValue,
				float64(node.Gd1InodesTotal),
				perBrickLabels...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.volumeBrickFreeInodes,
				prometheus.GaugeValue,
				float64(node.Gd1InodesFree),
				perBrickLabels...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.volumeBrickTotalBytes,
				prometheus.GaugeValue,
				float64(node.Capacity),
				perBrickLabels...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.volumeBrickFreeBytes,
				prometheus.GaugeValue,
				float64(node.Free),
				perBrickLabels...,
			)
		}
	}
	return
}
