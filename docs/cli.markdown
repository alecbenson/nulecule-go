# cli
--
    import "github.com/alecbenson/nulecule-go/atomicapp/cli"


## Usage

#### func  InitGeneralFlags

```go
func InitGeneralFlags()
```
A general command can be run with any subcommand after it

#### func  ParseCommands

```go
func ParseCommands(commands []*Command)
```
Iterates throuh the commands and runs the ones that are supplied in the command
line

#### type Command

```go
type Command struct {
	Name    string
	Func    func(cmd *Command)
	FlagSet *flag.FlagSet
}
```

A subcommand is used to represent a specific action in atomic: run, stop,
install, etc name - the name of the subcommand execute - the method to trigget
when this flag is Set flagset - the set of options that the command can be
configured with

#### func  InstallCommand

```go
func InstallCommand() *Command
```
InstallCommand returns an initialized CLI stop command

#### func  RunCommand

```go
func RunCommand() *Command
```
RunCommand returns an initialized CLI stop command

#### func  StopCommand

```go
func StopCommand() *Command
```
StopCommand returns an initialized CLI stop command
