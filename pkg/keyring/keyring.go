package keyring

import (
	"fmt"
	krPkg "github.com/99designs/keyring"
	"github.com/oslokommune/okctl/pkg/config"
)

type KeyType string

const (
	KeyTypeUserPassword =  "userPassword"
	KeyTypeGithubToken  =  "githubToken"
)

type Keyringer interface {
	Store(key KeyType, val string) error
	Fetch(key KeyType) (string, error)
}

type Keyring struct {
	ring krPkg.Keyring
}

var _ Keyringer = &Keyring{}

func New(keyring krPkg.Keyring) (*Keyring, error) {
	return &Keyring{
		ring: keyring,
	}, nil
}

func DefaultKeyring() (krPkg.Keyring, error) {
	ring, err := krPkg.Open(krPkg.Config{
		ServiceName: config.DefaultKeyringServiceName,
	})
	return ring, err
}

func (k *Keyring) Store(key KeyType, val string) error {
	if len(val) == 0 {
		return fmt.Errorf("key of type %s cannot store empty value", key)
	}
	return k.ring.Set(krPkg.Item{
		Key:  string(key),
		Data: []byte(val),
	})
}

func (k *Keyring) Fetch(key KeyType) (string, error) {
	get, err := k.ring.Get(string(key))
	if err != nil {
		return "", err
	}
	return string(get.Data), err
}

type InMemoryKeyring struct {
	krPkg.Keyring
	item krPkg.Item
}

var _ krPkg.Keyring = &InMemoryKeyring{}

func NewInMemoryKeyring() *InMemoryKeyring {
	return &InMemoryKeyring{}
}

func (i *InMemoryKeyring) Get(key string) (krPkg.Item, error) {
	return i.item, nil
}

func (i *InMemoryKeyring) Set(item krPkg.Item) error {
	i.item = item

	return nil
}