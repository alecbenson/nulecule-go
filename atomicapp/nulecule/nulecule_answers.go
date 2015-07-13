package nulecule

import (
	"errors"
	"path/filepath"

	"github.com/alecbenson/nulecule-go/atomicapp/constants"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"
)

//Needs to be built out
type Answers struct {
}

//LoadAnswers reads the answers from the provided path
func (b *Base) LoadAnswersFromPath(path string) (Answers, error) {
	result := Answers{}

	//Ensure that the path points to a valid file or directory
	if !utils.PathExists(path) {
		return result, errors.New("Bad answers filepath")
	}

	//If a file was provided...
	if utils.PathIsFile(path) {
		err := utils.ParseYAMLFile(path, &result)
		if err != nil {
			return result, errors.New("Failed to parse file")
		}
	}

	//if a directory was provided...
	if utils.PathIsDirectory(path) {
		//Construct the full path by combining with the ANSWERS_FILE constant
		path = filepath.Join(path, constants.ANSWERS_FILE)
		if !utils.PathExists(path) {
			return result, errors.New("Failed to read answers from path")
		}
		//..try to parse the file
		err := utils.ParseYAMLFile(path, &result)
		if err != nil {
			return result, errors.New("Failed to parse file")
		}
	}
	return result, nil
}
