package cli

import (
	"flag"

	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/constants"
	"github.com/alecbenson/nulecule-go/atomicapp/nulecule"
	"github.com/alecbenson/nulecule-go/atomicapp/provider"
)

func runFlagSet() *flag.FlagSet {
	runFlagSet := flag.NewFlagSet("run", flag.PanicOnError)
	runFlagSet.Bool("ask", false, "Ask for params even if the defaul value is provided")
	runFlagSet.String("write", "", "A file which will contain anwsers provided in interactive mode")
	runFlagSet.String("APP", "", "Path to the directory where the image is installed.")
	return runFlagSet
}


//Get the graph item of the main file
//For each component:
	//process
func dispatchGraph(graph *[]nulecule.Component) {
	if len(*graph) == 0 {
		logrus.Errorf("Graph not specified in %v file\n", constants.MAIN_FILE)
	}

	for _, component := range *graph {
		processComponent(component)
	}
}

//Get the provider from the nulecule base
//iniitialize a generic provider
//process all artifacts
//deploy the component
func processComponent(c nulecule.Component) {
	//The provider class to deploy with
	var prov provider.Provider
	//Iterate through the component and deploy all of its providers
	for providerName, artifactEntries := range c.Artifacts {
		processArtifacts(c, providerName)

		logrus.Infof("Deploying provider: %s...", providerName)
		prov = provider.New(providerName)
		prov.SetArtifacts(artifactEntries)
		prov.Init()
		prov.Deploy()
	}
}

//Get all artifacts provided by the component
//for every artifact in the component that belongs to provider,
	//recurse on inherited artifacts
//run load artifact on each
//apply an artifact template
//save the artifact to the destination directory (write the templated artifact to the dest. directory)
func processArtifacts(c nulecule.Component, provider string) {
	for _, artifactEntry := range c.Artifacts[provider] {

		//Process inherited artifacts as well
		if len(artifactEntry.Repo.Inherit) > 0 {
			for _, inherited_provider := range artifactEntry.Repo.Inherit {
				processArtifacts(c, inherited_provider)
			}
		}

	}
}

func runFunction() func(cmd *Command) {
	return func(cmd *Command) {
		//Set up parameters
		flags := cmd.FlagSet
		answersFile := getVal(flags, "write").(string)
		targetPath := getVal(flags, "APP").(string)
		ask := getVal(flags, "ask").(bool)
		targetFile := filepath.Join(targetPath, constants.MAIN_FILE)
		logrus.Debugf("RUN COMMAND: args are %v\n", targetFile)

		//Start run sequence
		base := &nulecule.Base{}
		base.New(targetPath, ask)
		base.ReadMainFile(targetFile)
		base.CheckSpecVersion()
		base.LoadAnswersFromPath(answersFile)
		base.CheckAllArtifacts()
		dispatchGraph(&base.MainfileData.Graph)
	}
}

//RunCommand returns an initialized CLI stop command
func RunCommand() *Command {
	return &Command{
		Name:    "run",
		Func:    runFunction(),
		FlagSet: runFlagSet(),
	}
}
