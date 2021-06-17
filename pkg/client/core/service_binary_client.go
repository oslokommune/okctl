package core

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type binaryService struct {
	config *config.Config
}

func (b *binaryService) Add(binary state.Binary) error {
	b.config.UserState.Binaries = append(b.config.UserState.Binaries, binary)

	err := b.config.WriteCurrentUserData()
	if err != nil {
		return fmt.Errorf("writing current user data: %w", err)
	}

	return nil
}

func (b *binaryService) Remove(binary state.Binary) error {
	b2, err := b.removeOne(b.config.UserState.Binaries, binary)
	if err != nil {
		return fmt.Errorf("removing from user state binaries: %w", err)
	}

	b.config.UserState.Binaries = b2

	err = b.config.WriteCurrentUserData()
	if err != nil {
		return fmt.Errorf("writing current user data: %w", err)
	}

	return nil
}

func (b *binaryService) removeOne(slice []state.Binary, binaryToRemove state.Binary) ([]state.Binary, error) {
	for i, binary := range slice {
		if binary.Id() == binaryToRemove.Id() {
			return append(slice[:i], slice[i+1:]...), nil
		}
	}

	return nil, fmt.Errorf(fmt.Sprintf("binary %s not found", binaryToRemove.Id()))
}

func (b *binaryService) List() []state.Binary {
	return b.config.UserState.Binaries
}

func NewBinaryService(config *config.Config) client.BinaryService {
	return &binaryService{
		config: config,
	}
}
