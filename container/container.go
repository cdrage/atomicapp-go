package container

import (
	"fmt"
	"github.com/cdrage/atomicapp-go/constants"
	"github.com/cdrage/atomicapp-go/utils"
	"os/exec"
)

const (
	CONTAINTER = "docker"
)

//Removes the container with the specified name
func Remove(containerName string) error {
	dockerRm := exec.Command(CONTAINTER, "rm", containerName)
	if _, err := utils.CheckCommandOutput(dockerRm, true); err != nil {
		return err
	}
	return nil
}

//Creates a container with the given name from the image provided
func Create(containerName, image string) error {
	dockerCreate := exec.Command(CONTAINTER, "create", "--name", containerName, image, "nop")
	if _, err := utils.CheckCommandOutput(dockerCreate, false); err != nil {
		return err
	}
	return nil
}

//Copies the data from the container to another
func Copy(containerName, copyFrom, copyTo string) error {
	cpField := fmt.Sprintf("%s:%s", containerName, constants.APP_ENT_PATH)
	dockerCp := exec.Command(CONTAINTER, "cp", cpField, copyTo)
	if _, err := utils.CheckCommandOutput(dockerCp, false); err != nil {
		return err
	}
	return nil
}

//Pulls the container with the given image name from docker
func Pull(image string) error {
	dockerPull := exec.Command(CONTAINTER, "pull", image)
	if _, err := utils.CheckCommandOutput(dockerPull, false); err != nil {
		return err
	}
	return nil
}
