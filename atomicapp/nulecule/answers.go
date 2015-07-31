package nulecule

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/constants"
	"github.com/alecbenson/nulecule-go/atomicapp/parser"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"
	"gopkg.in/yaml.v1"
)

//LoadAnswersFromPath reads the answers from the provided path
func (b *Base) LoadAnswersFromPath(path string) error {
	if !utils.PathExists(path) {
		return errors.New("Bad answers filepath")
	}
	p := parser.NewParser(path)

	//If a file was provided...
	if utils.PathIsFile(path) {
		err := p.Unmarshal(&b.AnswersData)
		if err != nil {
			return errors.New("Failed to parse file")
		}
	}

	//if a directory was provided...
	if utils.PathIsDirectory(path) {
		//Construct the full path by combining with the ANSWERS_FILE constant
		path = filepath.Join(path, constants.ANSWERS_FILE)
		if !utils.PathExists(path) {
			return errors.New("Failed to read answers from path")
		}
		//..try to parse the file
		err := p.Unmarshal(&b.AnswersData)
		if err != nil {
			return errors.New("Failed to parse file")
		}
	}
	return nil
}

//WriteAnswersSample creates a new answers.conf.sample file to update to
func (b *Base) WriteAnswersSample() {
	sampleAnswersPath := filepath.Join(b.Target(), constants.ANSWERS_FILE_SAMPLE)
	b.writeAnswers(sampleAnswersPath)
}

//Returns an *os.File pointing to the answers file. If one does not exist, it is created and returned.
func (b *Base) answersFile(sampleAnswersPath string) (*os.File, error) {
	if !utils.PathExists(sampleAnswersPath) {
		sampleAnswersFile, err := os.Create(sampleAnswersPath)
		if err != nil {
			logrus.Fatalf("Unable to create sample answers file: %s", err)
			return nil, err
		}
		return sampleAnswersFile, nil
	}
	sampleAnswersFile, err := os.Open(sampleAnswersPath)
	if err != nil {
		logrus.Errorf("Could not open answers file: %s", err)
		return nil, err
	}
	return sampleAnswersFile, nil
}

//writeAnswers will marshal the answers map and write it to the sample file
func (b *Base) writeAnswers(sampleAnswersPath string) error {
	sampleAnswersFile, err := b.answersFile(sampleAnswersPath)
	if err != nil {
		logrus.Errorf("Could not create/open answers file: %s", err)
		return err
	}
	data, err := yaml.Marshal(&b.AnswersData)
	if err != nil {
		logrus.Errorf("Failed to marshal answers to yaml: %s", err)
		return err
	}
	utils.WriteToFile(data, sampleAnswersFile)
	return nil
}

//UpdateAnswers updates the base answers data with a component
func (b *Base) UpdateAnswers(c Component) {
	if _, ok := b.AnswersData[c.Name]; ok {
		logrus.Infof("Component already exists in answers, updating...")
	}
	b.AnswersData[c.Name] = c.Params
}
