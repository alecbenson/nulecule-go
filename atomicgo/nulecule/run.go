package nulecule

import (
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicgo/constants"
	"github.com/alecbenson/nulecule-go/atomicgo/provider"
	"github.com/alecbenson/nulecule-go/atomicgo/utils"
)

//Run starts the nulecule run process to deploy a clustered application
func (b *Base) Run(ask bool, answersFile string) {
	if err := b.ReadMainFile(); err != nil {
		return
	}
	b.setAnswersDir(answersFile)
	b.CheckSpecVersion()
	b.LoadAnswers()
	b.CheckAllArtifacts()
	b.deployGraph(ask, answersFile)
}

//processComponent iterates through the artifacts in the component and deploy them
func (b *Base) processComponent(c *Component, ask bool) {
	//Iterate through the component and deploy all of its providers
	var prov provider.Provider
	for providerName, artifactEntries := range c.Artifacts {
		artifacts := arrangeArtifacts(artifactEntries)
		b.processArtifacts(c, providerName, ask)

		logrus.Infof("Deploying provider: %s...", providerName)
		prov = provider.New(providerName, b.Target(), b.DryRun)
		prov.SetArtifacts(artifacts)
		prov.Init()
		prov.Deploy()
	}
}

//processArtifacts iterates through each artifact entry
//It then substitutes the Nulecule parameters in and saves them in the workdir
func (b *Base) processArtifacts(c *Component, provider string, ask bool) {
	for _, artifactEntry := range c.Artifacts[provider] {
		//Process inherited artifacts as well
		if len(artifactEntry.Repo.Inherit) > 0 {
			for _, inheritedProvider := range artifactEntry.Repo.Inherit {
				b.processArtifacts(c, inheritedProvider, ask)
			}
		}
		//sanitize the prefix from the file path
		santitizedPath := utils.SanitizePath(artifactEntry.Path)
		//Form the absolute path of the artifact
		artifactPath := filepath.Join(b.Target(), santitizedPath)
		b.applyTemplate(artifactPath, c, ask)
	}
}

//deployGraph starts the deploy process for all components
func (b *Base) deployGraph(ask bool, answersFile string) {
	graph := b.MainfileData.Graph

	if len(graph) == 0 {
		logrus.Errorf("Graph not specified in %v file\n", constants.MAIN_FILE)
	}
	for _, c := range graph {
		if IsExternal(c) {
			externalDir := b.GetExternallAppDirectory(c)
			externalBase := New(externalDir, "", b.DryRun)
			externalBase.Run(ask, b.AnswersDir())
		}
		b.processComponent(&c, ask)
	}
}

func arrangeArtifacts(artifacts []ArtifactEntry) []string {
	result := make([]string, 0, len(artifacts))
	for _, artifact := range artifacts {
		if artifact.Path == "" {
			continue
		}
		sanitized := utils.SanitizePath(artifact.Path)
		result = append(result, sanitized)
	}
	return result
}
