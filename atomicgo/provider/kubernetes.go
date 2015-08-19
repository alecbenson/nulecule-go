package provider

import (
	"errors"
	"fmt"
	"strings"

	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicgo/constants"
	"github.com/alecbenson/nulecule-go/atomicgo/parser"
	"github.com/alecbenson/nulecule-go/atomicgo/utils"
)

//Kubernetes is a provider for kubernetes orchestration
type Kubernetes struct {
	*Config
	OrderedArtifacts []string
	KubeCtl          string
}

type kubeConfig struct {
	Kind     string
	Metadata metadata
}

type metadata struct {
	Name string
}

//NewKubernetes instantiates a new Kubernetes provider
func NewKubernetes(targetPath string, dryRun bool) *Kubernetes {
	provider := new(Kubernetes)
	provider.Config = new(Config)
	provider.targetPath = targetPath
	provider.Config.dryRun = dryRun
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
		//Form the absolute path of the artifact
		base := filepath.Base(artifact)
		fullPath := filepath.Join(p.WorkDirectory(), base)
		p.kubectlCmd(fullPath)
	}
	return nil
}

//Undeploy the kubernetes provider and reset replication controllers
func (p *Kubernetes) Undeploy() error {
	for _, artifact := range p.Artifacts() {
		base := filepath.Base(artifact)
		fullPath := filepath.Join(p.WorkDirectory(), base)

		markupParser := parser.NewParser(fullPath)
		result := &kubeConfig{}
		markupParser.Unmarshal(result)
		name := result.Metadata.Name
		kind := result.Kind

		if isReplicationController(kind) {
			p.resetReplica(name)
		}
		p.deletePod(fullPath)
	}
	return nil
}

//Returns true if the passed in kind is a replication controller
func isReplicationController(kind string) bool {
	if kind == "" {
		return false
	}

	validRC := []string{"replicationcontroller", "rc"}
	for _, rcName := range validRC {
		if strings.ToLower(kind) == rcName {
			return true
		}
	}
	return false
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
			logrus.Debugf("Found kubectl at %s", path)
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

//Resets the replication controller with the given name
func (p *Kubernetes) resetReplica(name string) error {
	if name == "" {
		logrus.Errorf("No pod name provided in call to reset replication controller.")
		return errors.New("No pod name provided")
	}
	namespaceFlag := fmt.Sprintf("--namespace=%s", p.Namespace)
	resizeCmd := exec.Command(p.KubeCtl, "resize", "rc", name, "--replicas=0", namespaceFlag)

	//In a dry run, we don't actually execute the commands, so we log and return here
	if p.DryRun() {
		logrus.Infof("DRY RUN: %s\n", resizeCmd.Args)
		return nil
	}

	if _, err := utils.CheckCommandOutput(resizeCmd, false); err != nil {
		logrus.Errorf("Error reseting kubernetes replication controller")
		return err
	}
	return nil
}

//deletePod removes the kubernetes pod
func (p *Kubernetes) deletePod(artifactPath string) error {
	if !utils.PathExists(artifactPath) {
		logrus.Errorf("No valid artifact could be found at %s", artifactPath)
		return errors.New("Path does not exist")
	}
	namespaceFlag := fmt.Sprintf("--namespace=%s", p.Namespace)
	deleteCmd := exec.Command(p.KubeCtl, "delete", "-f", artifactPath, namespaceFlag)

	//In a dry run, we don't actually execute the commands, so we log and return here
	if p.DryRun() {
		logrus.Infof("DRY RUN: %s\n", deleteCmd.Args)
		return nil
	}

	if _, err := utils.CheckCommandOutput(deleteCmd, false); err != nil {
		return err
	}
	return nil
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

	if _, err := utils.CheckCommandOutput(kubeCmd, false); err != nil {
		return err
	}
	return nil
}
