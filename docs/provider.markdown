# provider
--
    import "github.com/alecbenson/nulecule-go/atomicapp/provider"


## Usage

#### type Config

```go
type Config struct {

	//True if the provider is being called from within a container
	InContainer bool
	//Name of the namespace to run the provider in
	Namespace string
}
```

Config contains general configuration parameters that are used by each supported
provider

#### func (*Config) Artifacts

```go
func (c *Config) Artifacts() []nulecule.ArtifactEntry
```
Gets the list of artifacts belonging to the provider

#### func (*Config) CLIPath

```go
func (c *Config) CLIPath() []string
```
Gets a list of paths to search for the provider executable

#### func (*Config) DryRun

```go
func (c *Config) DryRun() bool
```
Gets the dry run value. In a dry run, no commands are actually run.

#### func (*Config) SetArtifacts

```go
func (c *Config) SetArtifacts(artifacts []nulecule.ArtifactEntry)
```
Sets the list of artifacts belonging to the provider

#### func (*Config) WorkDirectory

```go
func (c *Config) WorkDirectory() string
```

#### type Docker

```go
type Docker struct {
	*Config
}
```

Docker is a provider for Kubernetes

#### func  NewDocker

```go
func NewDocker(targetPath string) *Docker
```
Instantiates a new Kubernetes provider

#### func (*Docker) Deploy

```go
func (p *Docker) Deploy() error
```
Deploy the Docker Provider

#### func (*Docker) Init

```go
func (p *Docker) Init() error
```
Init the Docker provider

#### type Kubernetes

```go
type Kubernetes struct {
	*Config
	OrderedArtifacts []string
	KubeCtl          string
}
```

Kubernetes is a provider for kubernetes orchestration

#### func  NewKubernetes

```go
func NewKubernetes(targetPath string) *Kubernetes
```
Instantiates a new Kubernetes provider

#### func (*Kubernetes) Deploy

```go
func (p *Kubernetes) Deploy() error
```
Deploy the Kubernetes Provider

#### func (*Kubernetes) Init

```go
func (p *Kubernetes) Init() error
```
Init the Docker Provider

#### type Provider

```go
type Provider interface {
	Init() error
	Deploy() error

	CLIPath() []string
	Artifacts() []nulecule.ArtifactEntry
	SetArtifacts(artifacts []nulecule.ArtifactEntry)
	DryRun() bool
	// contains filtered or unexported methods
}
```

Provider defines functions that a provider plugin must include

#### func  New

```go
func New(provider string, targetPath string) Provider
```
New instantiates the provider with the give name
