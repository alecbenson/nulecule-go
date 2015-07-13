package provider

import (
	"errors"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"

	"os/exec"
	"strconv"
	"strings"
)

//Docker is a provider for Kubernetes
type Docker struct {
	*Config
}

//Instantiates a new Kubernetes provider
func NewDocker() *Docker {
	provider := new(Docker)
	provider.Config = new(Config)
	return provider
}

//Init the Docker provider
func (p *Docker) Init() error {
	p.checkVersion()
	return nil
}

//Deploy the Docker Provider
func (p *Docker) Deploy() error {
	//Iterate through artifact entries for the docker provider
	for _, artifactEntry := range p.Artifacts() {
		if !utils.PathExists(artifactEntry.Path) {
			logrus.Errorf("No such docker artifact: %v\n", artifactEntry.Path)
			return errors.New("No such docker artifact path")
		}

		file, err := ioutil.ReadFile(artifactEntry.Path)
		if err != nil {
			logrus.Errorf("Error reading artifact file: %v\n", artifactEntry.Path)
		}
		cmds := strings.Fields(string(file))
		dockerCmd := exec.Command(cmds[0], cmds[1:]...)

		//In a dry run, we don't actually execute the commands, so we log and return here
		if p.DryRun() {
			logrus.Infof("DRY RUN: %s\n", dockerCmd.Args)
			return nil
		}

		out, _ := dockerCmd.CombinedOutput()
		logrus.Infof(string(out))
	}
	return nil
}

//Check to ensure that we have a valid version of docker before deploying
func (p *Docker) checkVersion() error {
	daemonOut, _ := exec.Command("sudo", "docker", "version", "--format", "'{{.Server.ApiVersion}}'").Output()
	clientOut, _ := exec.Command("sudo", "docker", "version", "--format", "'{{.Client.ApiVersion}}'").Output()

	daemonString := strings.Trim(string(daemonOut), " '\n")
	clientString := strings.Trim(string(clientOut), " '\n")

	daemonVersion, err := strconv.ParseFloat(string(daemonString), 64)
	if err != nil {
		logrus.Errorf("Could not determine docker daemon version: %v\n", err)
		return errors.New("Unknown docker daemon version\n")
	}
	clientVersion, err := strconv.ParseFloat(string(clientString), 64)
	if err != nil {
		logrus.Errorf("Could not determine docker client version: %v\n", err)
		return errors.New("Unknown docker client version\n")
	}

	logrus.Debugf("Docker client is version %v and Docker daemon is version %v\n", clientVersion, daemonVersion)
	if clientVersion > daemonVersion {
		logrus.Errorf("Docker client is newer than daemon: please upgrade your daemon to a newer version\n")
		return errors.New("Docker daemon out of date")
	}
	return nil
}
