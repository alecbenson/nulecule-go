# nulecule
--
    import "github.com/alecbenson/nulecule-go/atomicapp/nulecule"


## Usage

#### func  ApplyTemplate

```go
func ApplyTemplate(artifactPath string, targetPath string, params []Param) ([]byte, error)
```
ApplyTemplate reads the file located at artifactPath and replaces all parameters
With those specified in the nulecule parameters objects artifactPath is the path
provided by the artifact TargetPath is the base directory to install the workdir
directory into

#### func  GetSourceImage

```go
func GetSourceImage(component Component) (string, error)
```
GetSourceImage fetches the sanitized source path of the image

#### func  IsExternal

```go
func IsExternal(component Component) bool
```
IsExternal returns true if the component is an external resource, and false if
the component is A local resource

#### func  SaveArtifact

```go
func SaveArtifact(data []byte, targetPath, name string) error
```
Saves a templated artifact to the .workdir directory. If .workdir does not
exist, it is created. data - a []byte of the templated file name - the name of
the file to write to

#### type Answers

```go
type Answers struct {
}
```

Needs to be built out

#### type ArtifactEntry

```go
type ArtifactEntry struct {
	Path string
	Repo SrcControlRepo
}
```

ArtifactEntry is a source control repository struct used to specify an artifact

#### func (*ArtifactEntry) SetYAML

```go
func (a *ArtifactEntry) SetYAML(tag string, value interface{}) bool
```
SetYAML is implemented by v1 of the go-yaml package. This method is invoked when
go-yaml attempts to parse an ArtifactEntry via Unmarshal() Because of the way
that Nulecule specifies an artifact, we have to define our own rules for parsing
Artifact structs. An artifact contains either a Source Control Repository object
or a URL. The former is represented as a structured set of parameters, while the
latter is represneted as an unlabeled string. Therefore, because an artifact
entry can be either a struct or an unlabeled string (and the parser is not smart
enough to know how to deal with this), we must do it ourselves.

#### type Base

```go
type Base struct {
	AnswersData  string
	ParamsData   string
	MainfileData *Mainfile

	Nodeps  bool
	AppID   string
	AppPath string

	Ask                bool
	WriteSampleAnswers bool
}
```

Base contains a set of nulecule config properties It is set by the atomicapp
subcommands

#### func (*Base) CheckAllArtifacts

```go
func (b *Base) CheckAllArtifacts()
```
CheckAllArtifacts will iterate through each entry in graph and check for
validity

#### func (*Base) CheckComponentArtifacts

```go
func (b *Base) CheckComponentArtifacts(c Component) []string
```
CheckComponentArtifacts will verify that valid artifacts exist for each provider
in the component

#### func (*Base) CheckSpecVersion

```go
func (b *Base) CheckSpecVersion() error
```
CheckSpecVersion verifies that a proper spec version has been provided

#### func (*Base) LoadAnswersFromPath

```go
func (b *Base) LoadAnswersFromPath(path string) (Answers, error)
```
LoadAnswers reads the answers from the provided path

#### func (*Base) New

```go
func (b *Base) New(targetPath string, ask bool) error
```
New creates a new base Nulecule object and initializes the fields

#### func (*Base) ReadMainFile

```go
func (b *Base) ReadMainFile(filepath string) error
```
ReadMainFile will read the Nulecule file and fill the MainfileData field

#### func (*Base) TargetPath

```go
func (b *Base) TargetPath() string
```

#### type Component

```go
type Component struct {
	Name      string
	Source    string
	Params    []Param
	Artifacts map[string][]ArtifactEntry
}
```

Component represents a graph item of the Nulecule file

#### type Constraint

```go
type Constraint struct {
	AllowedPattern string
	Description    string
}
```

Constraint is a struct representing a constaint for a parameter object

#### type Mainfile

```go
type Mainfile struct {
	Specversion string
	ID          string
	Graph       []Component
}
```

Mainfile is a struct representation of the Nulecule specification file

#### type Param

```go
type Param struct {
	Name        string
	Description string
	Constraints []Constraint
	Default     string
	Hidden      bool
}
```

Param represents the Component parameters

#### type SrcControlRepo

```go
type SrcControlRepo struct {
	Inherit []string
	Source  string
	Path    string
	Type    string
	Branch  string
	Tag     string
}
```
