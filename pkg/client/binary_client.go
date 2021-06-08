package client

import "github.com/oslokommune/okctl/pkg/config/state"

type BinaryService interface {
	Add(binary state.Binary) error
	Remove(binary state.Binary) error
	List() []state.Binary
}

type BinaryStore interface {
	Add(binary state.Binary) error
	Remove(binary state.Binary) error
}
