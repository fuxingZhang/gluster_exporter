package collector

import (
	"gluster_exporter/pkg/gluster"
	"gluster_exporter/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(NewPeerInfoCollector(), true)
}

func NewPeerInfoCollector() Collector {
	const subsystem = "peer"

	var (
		peerCountMetricLabels = []string{"instance"}
		peerSCMetricLabels    = []string{"instance", "hostname", "uuid"}
	)

	return &PeerInfoCollector{
		peerCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "count"),
			"Number of peers in cluster",
			peerCountMetricLabels,
			nil,
		),
		peerStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "status"),
			"Peer status info",
			peerSCMetricLabels,
			nil,
		),
		peerConnected: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "connected"),
			"Peer connection status",
			peerSCMetricLabels,
			nil,
		),
	}
}

type PeerInfoCollector struct {
	peerCount     *prometheus.Desc
	peerStatus    *prometheus.Desc
	peerConnected *prometheus.Desc
}

func (c *PeerInfoCollector) Name() string {
	return "peer_info"
}

func (c *PeerInfoCollector) Collect(ch chan<- prometheus.Metric) (err error) {
	peerID, err := gluster.Gluster().LocalPeerID()
	if err != nil {
		return
	}

	peers, err := gluster.Gluster().Peers()
	if err != nil {
		logger.Error(
			"msg", "[Gluster Peers] Error:",
			"peer", peerID,
			"err", err,
		)
		return
	}

	fqdn := "n/a"
	for _, peer := range peers {
		if peer.ID == peerID {
			// TODO: figure out which value of PeerAddresses may
			// be hostname -- or resolve ip ourselves
			fqdn = peer.PeerAddresses[0]
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.peerCount,
		prometheus.GaugeValue,
		float64(len(peers)),
		fqdn,
	)

	var connected int
	for _, peer := range peers {
		if peer.Online {
			connected = 1
		} else {
			connected = 0
		}
		// Only update glusterPeerStatus when we retrieved a
		// non-negative peer state, i.e. we're running with the GD1
		// backend.

		if peer.Gd1State > -1 {
			// newMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			// 	Namespace: namespace,
			// 	Subsystem: "peer",
			// 	Name:      "status",
			// 	Help:      "Peer status info",
			// }, []string{"instance", "hostname", "uuid"}).WithLabelValues(fqdn, peer.PeerAddresses[0], peer.ID)
			// newMetric.Set(float64(peer.Gd1State))
			// newMetric.Collect(ch)
			ch <- prometheus.MustNewConstMetric(
				c.peerStatus,
				prometheus.GaugeValue,
				float64(peer.Gd1State),
				fqdn,
				peer.PeerAddresses[0],
				peer.ID,
			)
		}

		ch <- prometheus.MustNewConstMetric(
			c.peerConnected,
			prometheus.GaugeValue,
			float64(connected),
			fqdn,
			peer.PeerAddresses[0],
			peer.ID,
		)
	}
	return
}
