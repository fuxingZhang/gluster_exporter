package gluster

import (
	"gluster_exporter/pkg/conf"
)

var gluster GInterface

func Init(config *conf.Config) {
	gluster = MakeGluster(config)
}

func Gluster() GInterface {
	return gluster
}
