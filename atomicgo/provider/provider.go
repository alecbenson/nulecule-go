package provider

import (
	"github.com/alecbenson/nulecule-go/atomicgo/constants"

	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
)

//Provider defines functions that a provider plugin must include
type Provider interface {
	Init() error
	Deploy() error
	Undeploy() error

	CLIPath() []string
	Artifacts() []string
	SetArtifacts(artifacts []string)
	DryRun() bool
	addCLIPaths(paths ...string)
}

//Config contains general configuration parameters that are used by
//each supported provider
type Config struct {
	//A list of artifacts for the provider to deploy
	artifacts []string
	//If true, run in Dry run mode.
	dryRun bool
	//A list of paths to check when trying to run the provider program
	cliPath []string
	//True if the provider is being called from within a container
	InContainer bool
	//Name of the namespace to run the provider in
	Namespace string
	//TargetPath is the path that all provider resources live in
	targetPath string
}

//addCLIPaths adds filepath(s) to check for the provider program in
func (c *Config) addCLIPaths(paths ...string) {
	c.cliPath = append(paths, c.cliPath...)
}

//CLIPath gets a list of paths to search for the provider executable
func (c *Config) CLIPath() []string {
	return c.cliPath
}

//Artifacts gets the list of artifacts belonging to the provider
func (c *Config) Artifacts() []string {
	return c.artifacts
}

//SetArtifacts sets the list of artifacts belonging to the provider
func (c *Config) SetArtifacts(artifacts []string) {
	c.artifacts = artifacts
}

//DryRun the dry run value. In a dry run, no commands are actually run.
func (c *Config) DryRun() bool {
	return c.dryRun
}

//WorkDirectory returns the absolute path to the .workdir directory
func (c *Config) WorkDirectory() string {
	return filepath.Join(c.targetPath, constants.WORKDIR)
}

//New instantiates the provider with the give name
func New(provider string, targetPath string, dryRun bool) Provider {
	sanitizedProvider := strings.ToLower(provider)
	switch sanitizedProvider {
	case "kubernetes":
		return NewKubernetes(targetPath, dryRun)
	case "docker":
		return NewDocker(targetPath, dryRun)
	case "openshift":
		return NewOpenshift(targetPath, dryRun)
	default:
		logrus.Errorf("Unrecognized provider: %s. Defaulting to %s", sanitizedProvider, constants.DEFAULT_PROVIDER)
		return NewKubernetes(targetPath, dryRun)
	}
}
