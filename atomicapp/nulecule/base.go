package nulecule

import (
	"errors"
	"path/filepath"

	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/constants"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"
	"gopkg.in/yaml.v1"
)

//Base contains a set of nulecule config properties
//It is set by the atomicapp subcommands
type Base struct {
	AnswersData        string
	ParamsData         string
	MainfileData       *Mainfile
	targetPath         string
	Nodeps             bool
	AppID              string
	AppPath            string
	provider           string
	app                string
	Ask                bool
	WriteSampleAnswers bool
}

//Mainfile is a struct representation of the Nulecule specification file
type Mainfile struct {
	Specversion string
	ID          string
	Graph       []Component
}

//Component represents a graph item of the Nulecule file
type Component struct {
	Name      string
	Source    string
	Params    []Param
	Artifacts map[string][]ArtifactEntry
}

//ArtifactEntry is a source control repository struct used to specify an artifact
type ArtifactEntry struct {
	Path string
	Repo SrcControlRepo
}

type SrcControlRepo struct {
	Inherit []string
	Source  string
	Path    string
	Type    string
	Branch  string
	Tag     string
}

//Param represents the Component parameters
type Param struct {
	Name        string
	Description string
	Constraints []Constraint
	Default     string
	Hidden      bool
}

//Constraint is a struct representing a constaint for a parameter object
type Constraint struct {
	AllowedPattern string
	Description    string
}

//New creates a new base Nulecule object and initializes the fields
func (b *Base) New(targetPath string, ask bool) error {
	b.setTargetPath(targetPath)
	b.Ask = ask
	b.MainfileData = &Mainfile{}
	return nil
}

//ReadMainFile will read the Nulecule file and fill the MainfileData field
func (b *Base) ReadMainFile(filepath string) error {
	//Check for valid path
	if !utils.PathExists(filepath) {
		return errors.New("File does not exist")
	}

	//Read the file
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	//Attempt to parse
	err = yaml.Unmarshal(file, b.MainfileData)
	if err != nil {
		logrus.Errorf("Error parsing Nulecule file: %v", err)
		return err
	}
	logrus.Debugf("Read following Nulecule file:\n\n %+v\n", *b.MainfileData)
	return nil
}

//CheckSpecVersion verifies that a proper spec version has been provided
func (b *Base) CheckSpecVersion() error {
	//Check for specversion property
	if b.MainfileData.Specversion == "" {
		logrus.Errorf("data corrupted: couldn't find specversion in main file")
		return errors.New("spec version check failed")
	}

	//Check for valid spec version
	spec := b.MainfileData.Specversion
	if spec == constants.NULECULESPECVERSION {
		logrus.Debug("version check successful: specversion == %s", constants.NULECULESPECVERSION)
	} else {
		logrus.Errorf("your version in %s file (%s) does not match supported version (%s)",
			constants.MAIN_FILE, spec, constants.NULECULESPECVERSION)
		return errors.New("spec version check failed")
	}
	return nil
}

//CheckAllArtifacts will iterate through each entry in graph and check for validity
func (b *Base) CheckAllArtifacts() {
	for _, c := range b.MainfileData.Graph {
		b.CheckComponentArtifacts(c)
		logrus.Infof("Checked artifacts for component: %s", c.Name)
	}
}

//CheckComponentArtifacts will verify that valid artifacts exist for each provider in the component
func (b *Base) CheckComponentArtifacts(c Component) []string {
	checkedProviders := make([]string, 0, 100)
	providerMap := c.Artifacts
	if len(providerMap) == 0 {
		logrus.Debugf("No artifacts for %s", c.Name)
		return checkedProviders
	}

	for provider := range providerMap {
		b.checkProviderArtifact(c, provider, &checkedProviders)
	}
	logrus.Debugf("Checked providers: %v\n", checkedProviders)
	return checkedProviders
}

//CheckComponentArtifacts will verify that valid artifacts exist for the specified provider
//The specified provider must be a member of the given component
func (b *Base) checkProviderArtifact(c Component, provider string, checkedProviders *[]string) {
	logrus.Debugf("Provider : %v", provider)

	//If the provider has already been checked, skip it
	if providerAlreadyChecked(checkedProviders, provider) {
		return
	}

	if artifacts, ok := c.Artifacts[provider]; ok {
		//Iterate through each individual artifact entry for each provider
		for _, artifactEntry := range artifacts {
			//If the entry has a path field, check it for validity
			if artifactEntry.Path != "" {
				fullPath := filepath.Join(b.targetPath, utils.SanitizePath(artifactEntry.Path))
				if utils.PathExists(fullPath) && utils.PathIsFile(fullPath) {
					logrus.Infof("Artifact %s: OK.", fullPath)
				} else {
					logrus.Errorf("Artifact %s: MISSING.", fullPath)
				}
			}
			//For this artfiact to be 'fully checked',
			//we need to verify that the inherited providers (if any) are valid as well
			b.checkInheritence(c, provider, artifactEntry.Repo.Inherit, checkedProviders)
		}
		*checkedProviders = append(*checkedProviders, provider)
	}
}

//Iterates through all providers in the inherit list and checks them for validity
func (b *Base) checkInheritence(c Component, provider string, inheritList []string, checkedProviders *[]string) {
	if len(inheritList) == 0 {
		return
	}

	for _, inheritProvider := range inheritList {
		//Verify that the provider does not inherit itself
		if provider == inheritProvider {
			logrus.Errorf("Provider %v cannot inherit itself. This entry will be ignored.\n", provider)
			return
		}
		//If the inherited provider has not already been checked, check it
		if providerAlreadyChecked(checkedProviders, inheritProvider) {
			continue
		}
		b.checkProviderArtifact(c, inheritProvider, checkedProviders)
	}
}

//Returns true if the provider is in the checked provider list
func providerAlreadyChecked(checkedProviders *[]string, provider string) bool {
	for _, checkedProvider := range *checkedProviders {
		if checkedProvider == provider {
			return true
		}
	}
	return false
}

func (b *Base) setTargetPath(filepath string) error {
	var err error
	if filepath == "" {
		logrus.Debugf("No target path provided, using current working directory")
		filepath, err = os.Getwd()
		if err != nil {
			logrus.Fatalf("Failed to get working directory")
			return errors.New("Failed to set target path")
		}
	}
	if utils.PathExists(filepath) {
		b.targetPath = filepath
		return nil
	}
	logrus.Fatalf("user provided an invalid target path: %s", filepath)
	return errors.New("Failed to set target path")
}

func (b *Base) TargetPath() string {
	return b.targetPath
}

//SetYAML is implemented by v1 of the go-yaml package. This method is invoked when go-yaml attempts to parse an ArtifactEntry via Unmarshal()
//Because of the way that Nulecule specifies an artifact, we have to define our own rules for parsing Artifact structs.
//An artifact contains either a Source Control Repository object or a URL. The former is represented as a structured set of parameters, while
//the latter is represneted as an unlabeled string. Therefore, because an artifact entry can be either a struct or an unlabeled string (and the parser is not
//smart enough to know how to deal with this), we must do it ourselves.
func (a *ArtifactEntry) SetYAML(tag string, value interface{}) bool {
	switch typedEntry := value.(type) {
	case string:
		a.Path = typedEntry
	case map[interface{}]interface{}:
		//In this case, we just want go-yaml to continue implicitly unmarshaling this struct.
		//Unfortunately, we have to remarshal the value to get the []byte type we need to call unmarshal again.
		b, err := yaml.Marshal(value)
		if err != nil {
			logrus.Errorf("could not set yaml for Artifact entry struct")
		}
		yaml.Unmarshal(b, &a.Repo)
	default:
		return false
	}
	return true
}
