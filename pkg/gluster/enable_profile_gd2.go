package gluster

import (
	"gluster_exporter/pkg/consts"
	"gluster_exporter/pkg/logger"

	"github.com/gluster/glusterd2/pkg/api"
)

// EnableVolumeProfiling enables profiling for a volume
func (g *GD2) EnableVolumeProfiling(volume Volume) error {
	client, err := initRESTClient(g.config)
	if err != nil {
		return err
	}

	value, exists := volume.Options[consts.CountFOPHitsGD2]
	if !exists {
		// Enable profiling for the volumes as its not set
		err := client.VolumeSet(
			volume.Name,
			api.VolOptionReq{
				Options: map[string]string{
					consts.CountFOPHitsGD2:       "on",
					consts.LatencyMeasurementGD2: "on",
				},
				VolOptionFlags: api.VolOptionFlags{
					AllowAdvanced: true,
				},
			},
		)
		if err != nil {
			return err
		}
	} else {
		if value == "off" {
			logger.Debug(
				"msg", "Volume profiling is explicitly disabled. No profile metrics would be exposed.",
				"volume", volume.Name,
			)
		}
	}
	return nil
}
