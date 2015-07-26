# utils
--
    import "github.com/alecbenson/nulecule-go/atomicapp/utils"


## Usage

```go
const (
	FILE_PREFIX   = "file://"
	DOCKER_PREFIX = "docker://"
)
```

#### func  ParseYAMLFile

```go
func ParseYAMLFile(filepath string, result interface{}) error
```
ParseYAMLFile will read in a path to a yaml file and parse to a map

#### func  PathExists

```go
func PathExists(path string) bool
```
PathExists verifies that the given path points to a valid file or directory

#### func  PathIsDirectory

```go
func PathIsDirectory(path string) bool
```
PathIsDirectory returns true if the path points to a valid directory

#### func  PathIsFile

```go
func PathIsFile(path string) bool
```
PathIsFile returns true if the path points to a vald file

#### func  SanitizePath

```go
func SanitizePath(url string) string
```
SanitizePath Strips the protocol from the provided path
