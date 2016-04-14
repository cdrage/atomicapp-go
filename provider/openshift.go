package provider

import ()

type Openshift struct {
	*Config
}

func NewOpenshift(targetPath string, dryRun bool) *Openshift {
	provider := new(Openshift)
	provider.Config = new(Config)
	provider.targetPath = targetPath
	provider.Config.dryRun = dryRun
	return provider
}

func (p *Openshift) Init() error {
	return nil
}

func (p *Openshift) Deploy() error {
	//look for template 'kind'
	//processTemplate()
	//call creaete on each artifact with --config
	return nil
}

//Invoke openshift's process command
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
