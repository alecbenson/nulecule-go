package nulecule

import (
	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicgo/constants"
	"github.com/alecbenson/nulecule-go/atomicgo/provider"
	"github.com/alecbenson/nulecule-go/atomicgo/utils"
	"os"
	"path/filepath"
)

//Stop starts the nulecule run process to undeploy a clustered application
func (b *Base) Stop() {
	if err := b.ReadMainFile(); err != nil {
		return
	}
	b.undeployAllProviders()
}

//undeployAllProviders will iterate through the components
//and undeploy all of its providers
func (b *Base) undeployAllProviders() {
	defer b.cleanWorkDirectory()
	graph := b.MainfileData.Graph
	var prov provider.Provider
	for _, c := range graph {
		if IsExternal(c) {
			externalDir := b.GetExternallAppDirectory(c)
			externalBase := New(externalDir, "", b.DryRun)
			externalBase.Stop()
			continue
		}
		for providerName, artifactEntries := range c.Artifacts {
			logrus.Infof("Undeploying provider: %s...", providerName)
			prov = provider.New(providerName, b.Target(), b.DryRun)
			artifacts := arrangeArtifacts(artifactEntries)
			prov.SetArtifacts(artifacts)
			prov.Init()
			prov.Undeploy()
		}
	}
}

//cleanWorkDirectory removes the .workdir directory once the graph has been deployed.
func (b *Base) cleanWorkDirectory() {
	workDirectory := filepath.Join(b.Target(), constants.WORKDIR)
	if utils.PathExists(workDirectory) {
		logrus.Debugf("Cleaning up work directory at %s\n", workDirectory)
		os.RemoveAll(workDirectory)
	}
}
