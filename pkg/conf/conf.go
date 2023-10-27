package conf

import (
	"os"
	"strings"

	"gluster_exporter/pkg/consts"
)

type ConfigCmdArgs struct {
	GlusterMgmt         string
	Glusterd2Endpoint   string
	GlusterCmd          string
	GlusterRemoteHost   string
	GlusterGlusterdSock string
	GlusterdWorkdir     string
}

// Config represents Glusterd1/Glusterd2 configurations
type Config struct {
	ConfigCmdArgs
	GlusterClusterID  string
	Glusterd2User     string
	Glusterd2Secret   string
	Glusterd2Cacert   string
	Glusterd2Insecure bool
	Timeout           int64
}

var conf *Config

func Init(args *ConfigCmdArgs) *Config {
	conf = &Config{
		ConfigCmdArgs: *args,
	}
	// by default, use glusterd (that is; GD1)
	if conf.GlusterMgmt == "" {
		conf.GlusterMgmt = consts.MgmtGlusterd
	}

	// Set the Gluster Configurations used in utils
	if conf.GlusterdWorkdir == "" {
		conf.GlusterdWorkdir = getDefaultGlusterdDir(conf.GlusterMgmt)
	}

	// If GD2_ENDPOINTS env variable is set, use that info
	// for making REST API calls
	if endpoint := os.Getenv(consts.EnvGD2Endpoints); endpoint != "" {
		conf.Glusterd2Endpoint = endpoint
	}
	// if there are multiple endpoints, get the first one
	if endpoint := conf.Glusterd2Endpoint; endpoint != "" {
		endpoint = strings.Replace(endpoint, ",", " ", -1)
		endpoint = strings.Fields(endpoint)[0]
		conf.Glusterd2Endpoint = endpoint
	}
	// if GLUSTER_CLUSTER_ID env variable is set, it gets the precedence
	if gClusterID := os.Getenv(consts.EnvGlusterClusterID); gClusterID != "" {
		conf.GlusterClusterID = gClusterID
	}
	// gluster cluster ID is still empty, put the default
	if conf.GlusterClusterID == "" {
		conf.GlusterClusterID = consts.DefaultGlusterClusterID
	}
	return conf
}

func ClusterID() string {
	return conf.GlusterClusterID
}

func Conf() *Config {
	return conf
}

// Below variables are set as flags during build time. The current
// values are just placeholders
var (
	defaultGlusterd1Workdir = "/var/lib/glusterd"
	defaultGlusterd2Workdir = "/var/lib/glusterd2"
)

func getDefaultGlusterdDir(mgmt string) string {
	if mgmt == consts.MgmtGlusterd2 {
		return defaultGlusterd2Workdir
	}
	return defaultGlusterd1Workdir
}
