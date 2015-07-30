package cli

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/constants"
	"github.com/alecbenson/nulecule-go/atomicapp/nulecule"
	"github.com/alecbenson/nulecule-go/atomicapp/provider"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"
)

func stopFlagSet() *flag.FlagSet {
	stopFlagSet := flag.NewFlagSet("stop", flag.PanicOnError)
	return stopFlagSet
}

func stopFunction() func(cmd *Command, args []string) {
	return func(cmd *Command, args []string) {
		//Set up parameters
		var (
			app    string
			target string
		)
		if len(args) > 0 {
			target = args[0]
		}
		defer cleanWorkDirectory(target)
		base := nulecule.New(target, app)
		if err := base.ReadMainFile(); err != nil {
			return
		}
		undeployAllProviders(base, base.Target())
	}
}

//StopCommand returns an initialized CLI stop command
func StopCommand() *Command {
	return &Command{
		Name:    "stop",
		Func:    stopFunction(),
		FlagSet: stopFlagSet(),
	}
}

func undeployAllProviders(base *nulecule.Base, targetPath string) {
	//Iterate through the component and undeploy all of its providers
	graph := base.MainfileData.Graph
	var prov provider.Provider
	for _, c := range graph {
		for providerName, artifactEntries := range c.Artifacts {
			logrus.Infof("Undeploying provider: %s...", providerName)
			prov = provider.New(providerName, targetPath)
			prov.SetArtifacts(artifactEntries)
			prov.Init()
			prov.Undeploy()
		}
	}
}

//cleanWorkDirectory removes the .workdir directory once the graph has been deployed.
func cleanWorkDirectory(targetPath string) {
	workDirectory := filepath.Join(targetPath, constants.WORKDIR)
	if utils.PathExists(workDirectory) {
		logrus.Debugf("Cleaning up work directory at %s\n", workDirectory)
		os.RemoveAll(workDirectory)
	}
}
