package gluster

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gluster_exporter/pkg/conf"
	"gluster_exporter/pkg/consts"
)

var (
	peerIDPattern = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}")
)

// IsLeader returns true or false based on whether the node is the leader of the cluster or not
func (g *GD1) IsLeader() (bool, error) {
	setDefaultConfig(g.config)
	peerList, err := g.Peers()
	if err != nil {
		return false, err
	}
	peerID, err := g.LocalPeerID()
	if err != nil {
		return false, err
	}
	var maxPeerID string
	//This for loop iterates among all the peers and finds the peer with the maximum UUID (lexicographically)
	for i, pr := range peerList {
		if pr.Online {
			if peerList[i].ID > maxPeerID {
				maxPeerID = peerList[i].ID
			}
		}
	}
	//Checks and returns true if maximum peerID is equal to the local peerID
	if maxPeerID == peerID {
		return true, nil
	}
	return false, nil
}

// IsLeader returns true or false based on whether the node is the leader of the cluster or not
func (g *GD2) IsLeader() (bool, error) {
	peerList, err := g.Peers()
	if err != nil {
		return false, err
	}
	peerID, err := g.LocalPeerID()
	if err != nil {
		return false, err
	}
	for _, pr := range peerList {
		if pr.Online {
			if peerID == pr.ID {
				return true, nil
			}
			return false, nil
		}
	}
	// This would imply none of the peers are online and return false
	return false, nil
}

// MakeGluster returns respective gluster obj based on configuration
func MakeGluster(gConfig *conf.Config) (gi GInterface) {
	if gConfig == nil {
		return nil
	}
	setDefaultConfig(gConfig)
	gi = &GD2{config: gConfig}
	if gConfig.GlusterMgmt == "" || gConfig.GlusterMgmt == consts.MgmtGlusterd {
		gi = &GD1{config: gConfig}
	}
	return
}

func readPeerID(fileStream io.ReadCloser, keywordID string) (string, error) {
	defer func() {
		err := fileStream.Close()
		if err != nil {
			// TODO: Log here
			return
		}
	}()

	scanner := bufio.NewScanner(fileStream)
	for scanner.Scan() {
		lines := strings.Split(scanner.Text(), "\n")
		for _, line := range lines {
			if strings.Contains(line, keywordID) {
				parts := strings.Split(line, "=")
				unformattedPeerID := parts[1]
				peerID := peerIDPattern.FindString(unformattedPeerID)
				if peerID == "" {
					return "", errors.New("unable to find peer address")
				}
				return peerID, nil
			}
		}
	}
	return "", errors.New("unable to find peer address")
}

// LocalPeerID returns local peer ID of glusterd
func (g *GD1) LocalPeerID() (string, error) {
	keywordID := "UUID"
	peeridFile := g.config.GlusterdWorkdir + "/glusterd.info"
	fileStream, err := os.Open(filepath.Clean(peeridFile))
	if err != nil {
		return "", err
	}
	return readPeerID(fileStream, keywordID)
}

// LocalPeerID returns local peer ID of glusterd2
func (g *GD2) LocalPeerID() (string, error) {
	keywordID := "peer-id"
	peeridFile := g.config.GlusterdWorkdir + "/uuid.toml"
	fileStream, err := os.Open(filepath.Clean(peeridFile))
	if err != nil {
		return "", err
	}
	return readPeerID(fileStream, keywordID)
}

// GetClusterID returns local clusterd ID
func GetClusterID() (clusterID string) {
	if clusterID = os.Getenv(consts.EnvGlusterClusterID); clusterID == "" {
		clusterID = consts.DefaultGlusterClusterID
	}
	return
}

// GConfig method returns the configuration
// this implements the 'conf.GConfigInterface'
func (g *GD1) Config() *conf.Config {
	return g.config
}

// GConfig method returns the configuration
// this implements the 'conf.GConfigInterface'
func (g *GD2) Config() *conf.Config {
	return g.config
}

// GetGlusterVersion gets the glusterfs version
func GetGlusterVersion() (string, error) {
	cmd := "glusterfs --version | head -1"
	bytes, err := ExecuteCmd(cmd)
	if err != nil {
		return "", err
	}
	stdout := string(bytes[:])
	fields := strings.Fields(stdout)
	return fields[1], err
}

// ExecuteCmd enables to execute system cmds and returns stdout, err
func ExecuteCmd(cmd string) ([]byte, error) {
	cmdfields := strings.Fields(cmd)
	cmdstr := cmdfields[0]
	if fullcmd, err := exec.LookPath(cmdfields[0]); err == nil {
		cmdstr = fullcmd
	}
	args := cmdfields[1:]
	out, err := exec.Command(cmdstr, args...).Output() // #nosec
	return out, err
}
