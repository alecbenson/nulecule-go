package cli

import "flag"

//A subcommand is used to represent a specific action in atomic: run, stop, install, etc
//name - the name of the subcommand
//execute - the method to trigget when this flag is Set
//flagset - the set of options that the command can be configured with
type Command struct {
	Name    string
	Func    func(cmd *Command)
	FlagSet *flag.FlagSet
}

//Iterates throuh the commands and runs the ones that are supplied in the command line
func ParseCommands(commands []*Command) {
	args := flag.Args()

	for _, cmd := range commands {
		if cmd.Name == args[0] || cmd.Name == "general" && cmd.Func != nil {
			cmd.FlagSet.Parse(args[1:])
			args = cmd.FlagSet.Args()
			cmd.Func(cmd)
			return
		}
	}
}

//Retrieves the argument from the provided flagset with the give name
func getVal(set *flag.FlagSet, arg string) interface{} {
	requestedFlag := set.Lookup(arg)
	if requestedFlag == nil {
		return nil
	}
	return requestedFlag.Value.(flag.Getter).Get()
}
