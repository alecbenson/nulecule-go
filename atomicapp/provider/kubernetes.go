package provider

import (
	"errors"
	"fmt"

	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"
	"github.com/alecbenson/nulecule-go/atomicapp/constants"
)

//Kubernetes is a provider for kubernetes orchestration
type Kubernetes struct {
	*Config
	OrderedArtifacts []string
	KubeCtl          string
}

//Instantiates a new Kubernetes provider
func NewKubernetes(targetPath string) *Kubernetes {
	provider := new(Kubernetes)
	provider.Config = new(Config)
	provider.targetPath = targetPath
	return provider
}

//Init the Docker Provider
func (p *Kubernetes) Init() error {
	var err error
	p.KubeCtl, err = p.findKubeCtlPath()
	if err != nil {
		return err
	}
	return nil
}

//Deploy the Kubernetes Provider
func (p *Kubernetes) Deploy() error {
	for _, artifact := range p.Artifacts() {
		//sanitize the prefix from the file path
		santitizedPath := utils.SanitizePath(artifact.Path)
		//Form the absolute path of the artifact
		fullPath := filepath.Join(p.targetPath, santitizedPath)
		p.kubectlCmd(fullPath)
	}
	return nil
}

//Determines the path of kubectl on the host
func (p *Kubernetes) findKubeCtlPath() (string, error) {
	if p.DryRun() {
		return "/usr/bin/kubectl", nil
	}

	//Insert some additional test paths to try
	p.addCLIPaths(
		"/usr/bin/kubectl",
		"/usr/local/bin/kubectl",
	)

	for _, path := range p.CLIPath() {
		//Add /host as a prefix if InContainer set to true
		if p.InContainer {
			path = filepath.Join("/host", path)
		}
		if utils.PathExists(path) {
			logrus.Infof("Found kubectl at %s", path)
			return path, nil
		}
	}

	logrus.Errorf("Unable to find valid kubectl installation. Tried the following locations: %v", p.CLIPath())
	return "", errors.New("No kubectl installation found")
}

//Ensures that all kubernetes commands are issued in order
//Reads the artifact file, gets the 'kind' field, and returns an ordered list
func (p *Kubernetes) prepareOrder() ([]string, error) {
	for artifact := range p.Artifacts() {
		logrus.Infof("%v\n", artifact)
	}
	return []string{}, nil
}

//Issues a command to kubectl.
//Path is the path of the kubernets file to pass to kubectl
func (p *Kubernetes) kubectlCmd(path string) error {
	if p.Namespace == "" {
		p.Namespace = constants.DEFAULT_NAMESPACE
	}
	path = utils.SanitizePath(path)
	namespaceFlag := fmt.Sprintf("--namespace=%s", p.Namespace)
	kubeCmd := exec.Command(p.KubeCtl, "create", "-f", path, namespaceFlag)

	//In a dry run, we don't actually execute the commands, so we log and return here
	if p.DryRun() {
		logrus.Infof("DRY RUN: %s\n", kubeCmd.Args)
		return nil
	}

	//Run the command
	out, _ := kubeCmd.CombinedOutput()
	logrus.Infof(string(out))
	return nil
}
