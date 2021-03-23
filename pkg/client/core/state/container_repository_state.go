package state

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type containerRepositoryState struct {
	state state.ContainerRepositorer
}

func (c *containerRepositoryState) GetContainerRepository(imageName string) (*client.ContainerRepository, error) {
	repository := c.state.GetContainerRepository(imageName)

	return &client.ContainerRepository{ImageName: repository.ImageName}, nil
}

func (c *containerRepositoryState) SaveContainerRepository(updatedRepository *client.ContainerRepository) (*store.Report, error) {
	repository := c.state.GetContainerRepository(updatedRepository.ImageName)

	repository.ImageName = updatedRepository.ImageName
	report, err := c.state.SaveContainerRepository(updatedRepository.ImageName, repository)
	if err != nil {
		return nil, err
	}

	report.Actions = append([]store.Action{
		{
			Name: "ContainerRepository",
			Path: fmt.Sprintf(""),
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

func (c *containerRepositoryState) RemoveContainerRepository(imageName string) (*store.Report, error) {
	report, err := c.state.DeleteContainerRepository(imageName)
	if err != nil {
		return nil, err
	}

	report.Actions = append([]store.Action{
		{
			Name: "ContainerRepository",
			Path: fmt.Sprintf("ImageName=%s", imageName),
			Type: "StateUpdate[removed]",
		},
	}, report.Actions...)

	return report, nil
}

// NewContainerRepositoryState returns an initialised state updater
func NewContainerRepositoryState(state state.ContainerRepositorer) client.ContainerRepositoryState {
	return &containerRepositoryState{
		state: state,
	}
}
