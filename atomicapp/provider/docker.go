package provider

import (
	"errors"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"

	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

//Docker is a provider for Kubernetes
type Docker struct {
	*Config
}

//NewDocker instantiates a new Kubernetes provider
func NewDocker(targetPath string, dryRun bool) *Docker {
	provider := new(Docker)
	provider.Config = new(Config)
	provider.targetPath = targetPath
	provider.Config.dryRun = dryRun
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
	for _, artifact := range p.Artifacts() {
		//Form the absolute path of the artifact
		base := filepath.Base(artifact)
		fullPath := filepath.Join(p.WorkDirectory(), base)

		if !utils.PathExists(fullPath) {
			logrus.Errorf("No such docker artifact: %v\n", fullPath)
			return errors.New("No such docker artifact path")
		}
		err := p.dockerCmd(fullPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Docker) Undeploy() error {
	return nil
}

//Issues the commands in the given file to docker
func (p *Docker) dockerCmd(artifactFile string) error {
	file, err := ioutil.ReadFile(artifactFile)
	if err != nil {
		logrus.Errorf("Error reading artifact file: %v\n", artifactFile)
		return errors.New("Unable to find docker artifact file")
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
	return nil
}

//Check to ensure that we have a valid version of docker before deploying
func (p *Docker) checkVersion() error {
	daemonCmd := exec.Command("docker", "version", "--format", "'{{.Server.ApiVersion}}'")
	daemonOut, err := utils.CheckCommandOutput(daemonCmd, true)
	if err != nil {
		return err
	}
	clientCmd := exec.Command("docker", "version", "--format", "'{{.Client.ApiVersion}}'")
	clientOut, err := utils.CheckCommandOutput(clientCmd, true)
	if err != nil {
		return err
	}

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
