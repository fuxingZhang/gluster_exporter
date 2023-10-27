package gluster

import (
	"encoding/json"
	"fmt"
	"gluster_exporter/pkg/logger"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// LVMStat represents LVM details
type LVMStat struct {
	Device          string
	UUID            string
	Name            string
	DataPercent     float64
	PoolLV          string
	Attr            string
	Size            float64
	Path            string
	MetadataSize    float64
	MetadataPercent float64
	VGName          string
	VGExtentTotal   float64
	VGExtentAlloc   float64
}

// ThinPoolStat represents thin pool LV details
type ThinPoolStat struct {
	ThinPoolName          string
	ThinPoolVGName        string
	ThinPoolDataTotal     float64
	ThinPoolDataUsed      float64
	ThinPoolMetadataTotal float64
	ThinPoolMetadataUsed  float64
}

// VGReport represents VG details
type VGReport struct {
	Report []VGs `json:"report"`
}

// VGs represents list VG Details
type VGs struct {
	Vgs []VGDetails `json:"vg"`
}

// VGDetails represents a single VG detail
type VGDetails struct {
	LVUUID          string `json:"lv_uuid"`
	LVName          string `json:"lv_name"`
	DataPercent     string `json:"data_percent"`
	PoolLV          string `json:"pool_lv"`
	LVAttr          string `json:"lv_attr"`
	LVSize          string `json:"lv_size"`
	LVPath          string `json:"lv_path"`
	LVMetadataSize  string `json:"lv_metadata_size"`
	MetadataPercent string `json:"metadata_percent"`
	VGName          string `json:"vg_name"`
	VGExtentTotal   string `json:"vg_extent_count"`
	VGExtentFree    string `json:"vg_free_count"`
}

func getLVS() ([]LVMStat, []ThinPoolStat, error) {
	cmd := "lvm vgs --unquoted --reportformat=json --noheading --nosuffix --units m -o lv_uuid,lv_name,data_percent,pool_lv,lv_attr,lv_size,lv_path,lv_metadata_size,metadata_percent,vg_name,vg_extent_count,vg_free_count"

	out, err := exec.Command("sh", "-c", cmd).Output() // #nosec
	lvmDet := []LVMStat{}
	thinPool := []ThinPoolStat{}
	var vgExtentFreeTemp float64
	if err != nil {
		logger.Error(
			"msg", "Error getting lvm usage details",
			"err", err,
		)
		return lvmDet, thinPool, err
	}
	var vgReport VGReport
	if err := json.Unmarshal(out, &vgReport); err != nil {
		logger.Debug("msg", "Error parsing lvm usage details", "err", err)
		return lvmDet, thinPool, err
	}

	for _, vg := range vgReport.Report[0].Vgs {
		var obj LVMStat
		obj.UUID = vg.LVUUID
		obj.Name = vg.LVName
		if vg.DataPercent == "" {
			obj.DataPercent = 0.0
		} else {
			if obj.DataPercent, err = strconv.ParseFloat(vg.DataPercent, 64); err != nil {
				logger.Error(
					"msg", "Error parsing DataPercent value of lvm usage",
					"err", err,
				)
				return lvmDet, thinPool, err
			}
		}
		obj.PoolLV = vg.PoolLV
		obj.Attr = vg.LVAttr
		if vg.LVSize == "" {
			obj.Size = 0.0
		} else {
			if obj.Size, err = strconv.ParseFloat(vg.LVSize, 64); err != nil {
				logger.Error(
					"msg", "Error parsing LVSize value of lvm usage",
					"err", err,
				)
				return lvmDet, thinPool, err
			}
		}
		obj.Path = vg.LVPath
		if vg.LVMetadataSize == "" {
			obj.MetadataSize = 0.0
		} else {
			if obj.MetadataSize, err = strconv.ParseFloat(vg.LVMetadataSize, 64); err != nil {
				logger.Error(
					"msg", "Error parsing LVMetadataSize value of lvm usage",
					"err", err,
				)
				return lvmDet, thinPool, err
			}
		}
		if vg.MetadataPercent == "" {
			obj.MetadataPercent = 0.0
		} else {
			obj.MetadataPercent, err = strconv.ParseFloat(vg.MetadataPercent, 64)
			if err != nil {
				logger.Error(
					"msg", "Error parsing MetadataPercent value of lvm usage",
					"err", err,
				)
				return lvmDet, thinPool, err
			}
		}
		if vg.VGExtentTotal == "" {
			obj.VGExtentTotal = 0.0
		} else {
			obj.VGExtentTotal, err = strconv.ParseFloat(vg.VGExtentTotal, 64)
			if err != nil {
				logger.Error(
					"msg", "Error parsing VGExtenTotal value of lvm usage",
					"err", err,
				)
				return lvmDet, thinPool, err
			}
		}
		if vg.VGExtentFree == "" {
			vgExtentFreeTemp = 0.0
		} else {
			vgExtentFreeTemp, err = strconv.ParseFloat(vg.VGExtentFree, 64)
			if err != nil {
				logger.Error(
					"msg", "Error parsing VGExtentAlloc value of lvm usage",
					"err", err,
				)
				return lvmDet, thinPool, err
			}
		}
		obj.VGExtentAlloc = obj.VGExtentTotal - vgExtentFreeTemp
		obj.VGName = vg.VGName
		if obj.Attr[0] == 't' {
			obj.Device = fmt.Sprintf("%s/%s", obj.VGName, obj.Name)
			var TPUsage ThinPoolStat
			TPUsage.ThinPoolName = obj.Name
			TPUsage.ThinPoolVGName = obj.VGName
			TPUsage.ThinPoolDataTotal = obj.Size
			TPUsage.ThinPoolDataUsed = (obj.Size * obj.DataPercent) / 100
			TPUsage.ThinPoolMetadataTotal = obj.MetadataSize
			TPUsage.ThinPoolMetadataUsed = (obj.MetadataSize * obj.MetadataPercent) / 100
			thinPool = append(thinPool, TPUsage)
		} else {
			obj.Device, err = filepath.EvalSymlinks(obj.Path)
			if err != nil {
				logger.Error(
					"msg", "Error evaluating realpath",
					"path", obj.Path,
					"err", err,
				)
				return lvmDet, thinPool, err
			}
		}
		lvmDet = append(lvmDet, obj)
	}
	return lvmDet, thinPool, nil
}

// ProcMounts represents list of items from /proc/mounts
type ProcMounts struct {
	Name         string
	Device       string
	FSType       string
	MountOptions string
}

func parseProcMounts() ([]ProcMounts, error) {
	procMounts := []ProcMounts{}
	b, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return procMounts, err
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "/") {
			tokens := strings.Fields(line)
			procMounts = append(procMounts,
				ProcMounts{Name: tokens[1], Device: tokens[0], FSType: tokens[2], MountOptions: tokens[3]})
		}
	}
	return procMounts, nil
}

func LvmUsage(path string) (stats []LVMStat, thinPoolStats []ThinPoolStat, err error) {
	mountPoints, err := parseProcMounts()
	if err != nil {
		return stats, thinPoolStats, err
	}
	var thinPoolNames []string
	lvs, tpStats, err := getLVS()
	if err != nil {
		return stats, thinPoolStats, err
	}
	for _, lv := range lvs {
		for _, mount := range mountPoints {
			dev, err := filepath.EvalSymlinks(mount.Device)
			if err != nil {
				logger.Error(
					"msg", "Error evaluating realpath",
					"path", mount.Device,
					"err", err,
				)
				continue
			}
			// Check if the logical volume is mounted as a gluster brick
			if lv.Device == dev && strings.HasPrefix(path, mount.Name) && strings.Contains(path, "/"+lv.VGName+"/"+lv.Name) {
				// Check if the LV is a thinly provisioned volume and if yes then get the thin pool LV name
				if lv.Attr[0] == 'V' {
					tpName := lv.PoolLV
					thinPoolNames = append(thinPoolNames, tpName)
				}
				stats = append(stats, lv)
			}
		}
	}
	// Iterate and select only those thin pool LVs whose thinly provisioned volumes are mounted as gluster bricks
	for _, tpName := range thinPoolNames {
		for _, tpStat := range tpStats {
			if tpName == tpStat.ThinPoolName {
				thinPoolStats = append(thinPoolStats, tpStat)
			}
		}
	}

	return stats, thinPoolStats, nil
}
