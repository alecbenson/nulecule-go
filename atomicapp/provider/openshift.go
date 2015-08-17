package provider

import ()

//Openshift is a provider for Kubernetes
type Openshift struct {
	*Config
}

//NewOpenshift instantiates a new Kubernetes provider
func NewOpenshift(targetPath string, dryRun bool) *Openshift {
	provider := new(Openshift)
	provider.Config = new(Config)
	provider.targetPath = targetPath
	provider.Config.dryRun = dryRun
	return provider
}

//Init the Openshift provider
func (p *Openshift) Init() error {
	return nil
}

//Deploy the Openshift Provider
func (p *Openshift) Deploy() error {
	//look for template 'kind'
	//processTemplate()
	//call creaete on each artifact with --config
	return nil
}

//ProcessTemplate will invoke openshift's process command
//when given a path to a template artifact,
func (p *Openshift) ProcessTemplate() error {
	//call --config=%s, process -f %s
	//write to output path
	return nil
}

//Undeploy removes the openshift provider and its configuration from the system
func (p *Openshift) Undeploy() error {
	return nil
}
