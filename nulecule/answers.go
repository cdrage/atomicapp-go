package nulecule

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/cdrage/atomicapp-go/constants"
	"github.com/cdrage/atomicapp-go/parser"
	"github.com/cdrage/atomicapp-go/utils"
	"gopkg.in/yaml.v1"
)

//A map of configuration parameters and values
type Answers map[string]string

//Unmarshals the answers from the answers.conf file into the base method
func (b *Base) LoadAnswers() error {
	//if a directory was provided...
	fp := b.AnswersDir()
	if utils.PathIsDirectory(fp) {
		//Construct the full path by combining with the ANSWERS_FILE constant
		fp = filepath.Join(fp, constants.ANSWERS_FILE)
		if !utils.PathExists(fp) {
			return errors.New("Failed to read answers from path")
		}
	}
	//..try to parse the file
	if !utils.PathExists(fp) {
		return errors.New("Bad answers filepath")
	}
	p := parser.NewParser(fp)
	err := p.Unmarshal(&b.AnswersData)
	if err != nil {
		return errors.New("Failed to parse file")
	}
	return nil
}

//Creates a new answers.conf.sample file
func (b *Base) WriteAnswersSample() error {
	sampleAnswersPath := filepath.Join(b.AnswersDir(), constants.ANSWERS_FILE_SAMPLE)
	if err := b.writeAnswers(sampleAnswersPath); err != nil {
		logrus.Errorf("Failed to write answers: %s", err)
		return err
	}
	return nil
}

//Marshal the answers map and write it to the sample file
func (b *Base) writeAnswers(sampleAnswersPath string) error {
	sampleAnswersFile, err := b.createAnswersFile(sampleAnswersPath)
	if err != nil {
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

//Returns an *os.File pointing to the answers file. If one does not exist, it is created and returned.
func (b *Base) createAnswersFile(sampleAnswersPath string) (*os.File, error) {
	if utils.PathExists(sampleAnswersPath) {
		return os.OpenFile(sampleAnswersPath, os.O_APPEND|os.O_WRONLY, 0600)
	}
	sampleAnswersFile, err := os.Create(sampleAnswersPath)
	if err != nil {
		logrus.Fatalf("Unable to create sample answers file: %s", err)
		return nil, err
	}
	return sampleAnswersFile, nil
}

//Update the base answers data with a component
func (b *Base) updateComponentAnswers(c Component) error {
	componentAnswers := b.AnswersData[c.Name]
	if componentAnswers == nil {
		componentAnswers = make(Answers)
	}
	for _, p := range c.Params {
		paramName := p.Name
		value, err := b.getParamValue(p)
		if err != nil {
			return err
		}
		componentAnswers[paramName] = value
	}
	//Set the value to the resulting map before returning
	b.AnswersData[c.Name] = componentAnswers
	return nil
}
