package cli

import (
	"flag"
	"os"

	"github.com/alecbenson/nulecule-go/atomicapp/logging"
)

//InitGeneralFlags provides general configuration.
//Any subcommand can be run after it
func InitGeneralFlags() {
	flag.Bool("version", false, "Displays the software version")
	logLevel := flag.Int("log-level", 0, "0 - Quiet | 1 - Normal | 2 - Verbose logging")
	flag.String("answers", "", "Set the path to the answers file")
	flag.Bool("dry-run", false, "Don't actually call provider. The commands that"+
		"should be run will be sent to stdout but not run.")
	flag.Parse()

	//Verify that arguments have been provided
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(0)
	}

	logging.SetLogLevel(*logLevel)
}
