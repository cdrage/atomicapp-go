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

type Base struct {
	AnswersData        map[string]Answers
	AnswersDirectory   string
	MainfileData       *Mainfile
	TargetPath         string
	Nodeps             bool
	DryRun             bool
	AppPath            string
	Application        string
	WriteSampleAnswers bool
}

//Struct representation of the Nulecule specification file
type Mainfile struct {
	Specversion string
	ID          string
	Graph       []Component
}

//Representing the graph file of Nulecule
type Component struct {
	Name      string
	Source    string
	Params    []Param
	Artifacts map[string][]ArtifactEntry
}

//New creates a new base Nulecule object and initializes the fields
func New(TargetPath, app string, dryRun bool) *Base {
	b := &Base{Application: app}
	b.setTargetPath(TargetPath)
	b.setAnswersDir(b.Target())
	b.MainfileData = &Mainfile{}
	b.AnswersData = make(map[string]Answers)
	b.DryRun = dryRun
	return b
}

//Read the Nulecule file and fill the MainfileData field
func (b *Base) ReadMainFile() error {
	//Check for valid path
	targetFile := filepath.Join(b.Target(), constants.MAIN_FILE)
	if !utils.PathExists(targetFile) {
		logrus.Fatalf("Could not find %s file in %s", constants.MAIN_FILE, b.Target())
		return errors.New("File does not exist")
	}

	//Attempt to parse
	p := parser.NewParser(targetFile)
	err := p.Unmarshal(b.MainfileData)
	if err != nil {
		logrus.Errorf("Error parsing Nulecule file: %v", err)
		return err
	}
	return nil
}

//Verifies that a proper spec version has been provided
func (b *Base) CheckSpecVersion() error {
	//Check for specversion property
	if b.MainfileData.Specversion == "" {
		logrus.Errorf("data corrupted: couldn't find specversion in main file")
		return errors.New("spec version check failed")
	}

	//Check for valid spec version
	spec := b.MainfileData.Specversion
	if spec == constants.NULECULESPECVERSION {
		logrus.Debugf("version check successful: specversion == %s", constants.NULECULESPECVERSION)
	} else {
		logrus.Errorf("your version in %s file (%s) does not match supported version (%s)",
			constants.MAIN_FILE, spec, constants.NULECULESPECVERSION)
		return errors.New("spec version check failed")
	}
	return nil
}

//Iterate through each entry in graph and check for validity
func (b *Base) CheckAllArtifacts() {
	for _, c := range b.MainfileData.Graph {
		b.CheckComponentArtifacts(c)
		logrus.Infof("Checked artifacts for component: %s", c.Name)
	}
}

//Verify that valid artifacts exist for each provider in the component
func (b *Base) CheckComponentArtifacts(c Component) []string {
	checkedProviders := make([]string, 0, 100)
	providerMap := c.Artifacts
	if len(providerMap) == 0 {
		logrus.Debugf("No artifacts for %s", c.Name)
		return checkedProviders
	}

	for provider := range providerMap {
		b.checkProviderArtifact(c, provider, &checkedProviders)
	}
	logrus.Debugf("Checked providers: %v\n", checkedProviders)
	return checkedProviders
}

//Verify that valid artifacts exist for the specified provider
//The specified provider must be a member of the given component
func (b *Base) checkProviderArtifact(c Component, provider string, checkedProviders *[]string) {
	logrus.Debugf("Provider : %v", provider)

	//If the provider has already been checked, skip it
	if providerAlreadyChecked(checkedProviders, provider) {
		return
	}

	if artifacts, ok := c.Artifacts[provider]; ok {
		//Iterate through each individual artifact entry for each provider
		for _, artifactEntry := range artifacts {
			//If the entry has a path field, check it for validity
			if artifactEntry.Path != "" {
				fullPath := filepath.Join(b.Target(), utils.SanitizePath(artifactEntry.Path))
				if utils.PathExists(fullPath) && utils.PathIsFile(fullPath) {
					logrus.Infof("Artifact %s: OK.", fullPath)
				} else {
					logrus.Errorf("Artifact %s: MISSING.", fullPath)
				}
			}
			//For this artifact to be 'fully checked',
			//we need to verify that the inherited providers (if any) are valid as well
			b.checkInheritence(c, provider, artifactEntry.Repo.Inherit, checkedProviders)
		}
		*checkedProviders = append(*checkedProviders, provider)
	}
}

//Iterates through all providers in the inherit list and checks them for validity
func (b *Base) checkInheritence(c Component, provider string, inheritList []string, checkedProviders *[]string) {
	if len(inheritList) == 0 {
		return
	}

	for _, inheritProvider := range inheritList {
		//Verify that the provider does not inherit itself
		if provider == inheritProvider {
			logrus.Errorf("Provider %v cannot inherit itself. This entry will be ignored.\n", provider)
			return
		}
		//If the inherited provider has not already been checked, check it
		if providerAlreadyChecked(checkedProviders, inheritProvider) {
			continue
		}
		b.checkProviderArtifact(c, inheritProvider, checkedProviders)
	}
}

//Returns true if the provider is in the checked provider list
func providerAlreadyChecked(checkedProviders *[]string, provider string) bool {
	for _, checkedProvider := range *checkedProviders {
		if checkedProvider == provider {
			return true
		}
	}
	return false
}

func (b *Base) setTargetPath(target string) error {
	//If no target is specified or if the user specifies a '.',
	//then use the current working directory
	if target == "" || !utils.PathExists(target) {
		cwd, err := os.Getwd()
		if err != nil {
			logrus.Fatalf("Failed to get working directory")
			return errors.New("Failed to set target path")
		}
		b.TargetPath = cwd
		return nil
	}
	b.TargetPath = target
	return nil
}

//Target is a getter for nulecule base's target field
func (b *Base) Target() string {
	if b.TargetPath == "" {
		b.setTargetPath("")
	}
	return b.TargetPath
}

//AnswersDir returns the base directory in which the answers file lives
func (b *Base) AnswersDir() string {
	return b.AnswersDirectory
}

func (b *Base) setAnswersDir(AnswersDirectory string) error {
	if !utils.PathExists(AnswersDirectory) {
		if AnswersDirectory != "" {
			logrus.Warnf("Invalid answers directory provided: '%s'. Using '%s' instead", AnswersDirectory, b.Target())
		}
		b.AnswersDirectory = b.Target()
		return errors.New("Using target path as answers directory")
	}
	b.AnswersDirectory = AnswersDirectory
	return nil
}

//App is a getter for the Nulecule base app field
func (b *Base) App() string {
	return b.Application
}

//SetYAML is implemented by v1 of the go-yaml package. This method is invoked when go-yaml attempts to parse an ArtifactEntry via Unmarshal()
//Because of the way that Nulecule specifies an artifact, we have to define our own rules for parsing Artifact structs.
//An artifact contains either a Source Control Repository object or a URL. The former is represented as a structured set of parameters, while
//the latter is represneted as an unlabeled string. Therefore, because an artifact entry can be either a struct or an unlabeled string (and the parser is not
//smart enough to know how to deal with this), we must do it ourselves.
func (a *ArtifactEntry) SetYAML(tag string, value interface{}) bool {
	switch typedEntry := value.(type) {
	case string:
		a.Path = typedEntry
	case map[interface{}]interface{}:
		//In this case, we just want go-yaml to continue implicitly unmarshaling this struct.
		//Unfortunately, we have to remarshal the value to get the []byte type we need to call unmarshal again.
		b, err := yaml.Marshal(value)
		if err != nil {
			logrus.Errorf("could not set yaml for Artifact entry struct")
		}
		yaml.Unmarshal(b, &a.Repo)
	default:
		return false
	}
	return true
}
