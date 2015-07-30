package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"unicode"
)

//Parser is responsible for reading a markup language and parsing
//it into a struct
type Parser struct {
	markupFile string
}

//NewParser instantiates a new parser for reading markup files
func NewParser(markupFile string) *Parser {
	return &Parser{markupFile: markupFile}
}

//Unmarshal arses a markup file into the given interface object
func (p *Parser) Unmarshal(result interface{}) error {
	if !utils.PathExists(p.markupFile) {
		return errors.New("File does not exist")
	}

	f, err := ioutil.ReadFile(p.markupFile)
	if err != nil {
		return err
	}

	if isJSON(f) {
		return json.Unmarshal(f, result)
	}
	return yaml.Unmarshal(f, result)
}

//isJSON checks the file for an opening brace
//to determine if it is a json file or not
func isJSON(data []byte) bool {
	spacesRemoved := bytes.TrimLeftFunc(data, unicode.IsSpace)
	prefix := []byte("{")
	if bytes.HasPrefix(spacesRemoved, prefix) {
		return true
	}
	return false
}
