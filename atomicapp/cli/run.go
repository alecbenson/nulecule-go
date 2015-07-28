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

//RunCommand returns an initialized CLI stop command
func RunCommand() *Command {
	return &Command{
		Name:    "run",
		Func:    runFunction(),
		FlagSet: runFlagSet(),
	}
}

func runFlagSet() *flag.FlagSet {
	runFlagSet := flag.NewFlagSet("run", flag.PanicOnError)
	runFlagSet.Bool("ask", false, "Ask for params even if the defaul value is provided")
	runFlagSet.String("write", "", "A file which will contain anwsers provided in interactive mode")
	return runFlagSet
}

func runFunction() func(cmd *Command, args []string) {
	return func(cmd *Command, args []string) {
		//Set up parameters
		var (
			app    string
			target string
		)
		if len(args) > 0 {
			target = args[0]
		}
		flags := cmd.FlagSet
		answersFile := getVal(flags, "write").(string)
		//Start run sequence
		defer cleanWorkDirectory(target)
		base := nulecule.New(target, app)
		if err := base.ReadMainFile(); err != nil {
			return
		}
		base.CheckSpecVersion()
		base.LoadAnswersFromPath(answersFile)
		base.CheckAllArtifacts()
		deployGraph(base)
	}
}

//deployGraph starts the deploy process for all components
func deployGraph(b *nulecule.Base) {
	graph := b.MainfileData.Graph
	targetPath := b.Target()

	if len(graph) == 0 {
		logrus.Errorf("Graph not specified in %v file\n", constants.MAIN_FILE)
	}
	for _, component := range graph {
		processComponent(component, targetPath)
	}
}

//processComponent iterates through the artifacts in the component and deploy them
func processComponent(c nulecule.Component, targetPath string) {
	//Iterate through the component and deploy all of its providers
	var prov provider.Provider
	for providerName, artifactEntries := range c.Artifacts {
		processArtifacts(c, providerName, targetPath)

		logrus.Infof("Deploying provider: %s...", providerName)
		prov = provider.New(providerName, targetPath)
		prov.SetArtifacts(artifactEntries)
		prov.Init()
		prov.Deploy()
	}
}

//processArtifacts iterates through each artifact entry
//It then substitutes the Nulecule parameters in and saves them in the workdir
func processArtifacts(c nulecule.Component, provider, targetPath string) {
	for _, artifactEntry := range c.Artifacts[provider] {
		//Process inherited artifacts as well
		if len(artifactEntry.Repo.Inherit) > 0 {
			for _, inheritedProvider := range artifactEntry.Repo.Inherit {
				processArtifacts(c, inheritedProvider, targetPath)
			}
		}
		//sanitize the prefix from the file path
		santitizedPath := utils.SanitizePath(artifactEntry.Path)
		//Form the absolute path of the artifact
		fullPath := filepath.Join(targetPath, santitizedPath)
		nulecule.ApplyTemplate(fullPath, targetPath, c.Params)
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
