package collector

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"

	"gluster_exporter/pkg/conf"
	"gluster_exporter/pkg/gluster"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(NewPeerCountsCollector(), true)
}

func NewPeerCountsCollector() Collector {
	const subsystem = ""

	var (
		// general metric labels
		gnrlMetricLabels = []string{"cluster_id", "name", "peer_id"}
		// an additional information of 'vg_name' is added
		// this specifies which Volume Group the LV or ThinPool count belongs
		withVgMetricLabels = []string{"cluster_id", "name", "peer_id", "vg_name"}
	)

	return &PeerCountsCollector{
		pvCountVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "pv_count",
			Help:      "No: of Physical Volumes",
		}, gnrlMetricLabels),
		lvCountVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "lv_count",
			Help:      "No: of Logical Volumes in a Volume Group",
		}, withVgMetricLabels),
		vgCountVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "vg_count",
			Help:      "No: of Volume Groups",
		}, gnrlMetricLabels),
		tpCountVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "thinpool_count",
			Help:      "No: of thinpools in a Volume Group",
		}, withVgMetricLabels),
	}
}

type PeerCountsCollector struct {
	pvCountVec *prometheus.GaugeVec
	lvCountVec *prometheus.GaugeVec
	vgCountVec *prometheus.GaugeVec
	tpCountVec *prometheus.GaugeVec
}

func (c *PeerCountsCollector) Name() string {
	return "peer_counts"
}

func (c *PeerCountsCollector) Collect(ch chan<- prometheus.Metric) (err error) {
	peerID, err := gluster.Gluster().LocalPeerID()
	if err != nil {
		return
	}

	pMetrics, err := NewPeerMetrics()
	if err != nil {
		return
	}

	genrlLbls := prometheus.Labels{
		"cluster_id": conf.ClusterID(),
		"name":       "Physical_Volumes",
		"peer_id":    peerID,
	}
	collectGaugeVec(c.pvCountVec, genrlLbls, float64(pMetrics.PVCount), ch)

	genrlLbls = prometheus.Labels{
		"cluster_id": conf.ClusterID(),
		"name":       "Volume_Groups",
		"peer_id":    peerID,
	}
	collectGaugeVec(c.vgCountVec, genrlLbls, float64(pMetrics.VGCount), ch)

	// logical volume counts are added specific to each VG
	for vgName, lvCount := range pMetrics.LVCountMap {
		genrlLbls = prometheus.Labels{
			"cluster_id": conf.ClusterID(),
			"name":       "Logical_Volumes",
			"peer_id":    peerID,
			"vg_name":    vgName,
		}
		collectGaugeVec(c.lvCountVec, genrlLbls, float64(lvCount), ch)
	}
	// similarly thinpool counts are also added per VG
	for vgName, tpCount := range pMetrics.ThinPoolCountMap {
		genrlLbls = prometheus.Labels{
			"cluster_id": conf.ClusterID(),
			"name":       "ThinPool_Count",
			"peer_id":    peerID,
			"vg_name":    vgName,
		}
		collectGaugeVec(c.tpCountVec, genrlLbls, float64(tpCount), ch)
	}
	return
}

// PeerMetrics : exposes PV, LV, VG counts
type PeerMetrics struct {
	PVCount          int            // total Physical Volume counts
	LVCountMap       map[string]int // collects lv count for each Volume Group
	ThinPoolCountMap map[string]int // collects thinpool count for each Volume Group
	VGCount          int            // no: of Volume Groups
}

type myVGDetails struct {
	LVUUID     string `json:"lv_uuid"`
	LVName     string `json:"lv_name"`
	PoolLV     string `json:"pool_lv"`
	VGName     string `json:"vg_name"`
	LVPath     string `json:"lv_path"`
	LVCount    string `json:"lv_count"`
	PVCount    string `json:"pv_count"`
	PoolLVUUID string `json:"pool_lv_uuid"`
	LVAttr     string `json:"lv_attr"`
}

// NewPeerMetrics : provides a way to get the consolidated metrics (such PV, LV, VG counts)
func NewPeerMetrics() (*PeerMetrics, error) {
	cmdStr := "lvm vgs --noheading --reportformat=json -o lv_uuid,lv_name,pool_lv,vg_name,lv_path,lv_count,pv_count,pool_lv_uuid,lv_attr"
	outBs, err := gluster.ExecuteCmd(cmdStr)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(outBs))
	var vgs []myVGDetails
	// collect the details from the JSON output
	for {
		t, decErr := dec.Token()
		if decErr == io.EOF {
			break
		}
		if decErr != nil {
			err = errors.New("Unable to parse JSON output: " + decErr.Error())
			return nil, err
		}
		// if the token is 'vg', collect/decode the details into VGDetails array
		if t == "vg" {
			decErr = dec.Decode(&vgs)
		}
		if decErr != nil {
			err = errors.New("JSON output changed, parse failed: " + decErr.Error())
			return nil, err
		}
	}
	pMetrics := &PeerMetrics{
		PVCount:          0,
		VGCount:          0,
		LVCountMap:       make(map[string]int),
		ThinPoolCountMap: make(map[string]int),
	}
	var vgMap = make(map[string]myVGDetails)
	for _, vg := range vgs {
		// collect the unique vgs into the map
		if _, ok := vgMap[vg.VGName]; !ok {
			vgMap[vg.VGName] = vg
			// increment the VG counter, for each new Volume Group
			pMetrics.VGCount++
			// if there is no error while integer conversion, add that number
			if count, convErr := strconv.Atoi(strings.TrimSpace(vg.PVCount)); convErr == nil {
				pMetrics.PVCount += count
			}
			// by default set the LV count to Zero
			pMetrics.LVCountMap[vg.VGName] = 0
			if count, convErr := strconv.Atoi(strings.TrimSpace(vg.LVCount)); convErr == nil {
				// if there are no errors, update the LV count
				pMetrics.LVCountMap[vg.VGName] = count
			}
		}
		// before adding into 'thinPoolMap', check the attribute
		// if attribute string starts with 't', consider it as thinpool
		//
		// converting to a rune array because of 'utf-8' enconded
		// string handling in 'Go'
		if len(vg.LVAttr) > 0 && []rune(vg.LVAttr)[0] == 't' {
			// increment the thin pool count for that particular VG
			pMetrics.ThinPoolCountMap[vg.VGName]++
		}
	}
	return pMetrics, nil
}
