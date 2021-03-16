// Package keyring handles secrets stored on client machine using okctl
package keyring

import (
	"fmt"
	"runtime"

	"github.com/oslokommune/okctl/pkg/config/constant"

	krPkg "github.com/99designs/keyring"
)

// KeyType : the type of key we are storing or fetching from keyring
type KeyType string

const (
	// KeyTypeUserPassword user password used to login with saml (AD-password)
	KeyTypeUserPassword = "userPassword"
	// KeyTypeGithubToken github token for session
	KeyTypeGithubToken = "githubToken"
)

// Keyringer exposes functions needed from keyring
type Keyringer interface {
	Store(key KeyType, val string) error
	Fetch(key KeyType) (string, error)
}

// Keyring : wrapper for keyring package
type Keyring struct {
	ring krPkg.Keyring
}

var _ Keyringer = &Keyring{}

// New creates a new keyring
func New(keyring krPkg.Keyring, debug bool) (*Keyring, error) {
	krPkg.Debug = debug

	return &Keyring{
		ring: keyring,
	}, nil
}

// DefaultKeyringForOS is the default keyring to store client secrets
func DefaultKeyringForOS() (krPkg.Keyring, error) {
	cfg := krPkg.Config{
		ServiceName: constant.DefaultKeyringServiceName,
	}

	switch runtime.GOOS {
	case "darwin":
		cfg.AllowedBackends = []krPkg.BackendType{krPkg.KeychainBackend}
	case "linux":
	default:
		return nil, fmt.Errorf("no supported keyring backends for your operating system: %s", runtime.GOOS)
	}

	return krPkg.Open(cfg)
}

// Store a value with given keytype and value in keyring
func (k *Keyring) Store(key KeyType, val string) error {
	if len(val) == 0 {
		return fmt.Errorf("key of type %s cannot store empty value", key)
	}

	return k.ring.Set(krPkg.Item{
		Key:  string(key),
		Data: []byte(val),
	})
}

// Fetch a value with given keytype from keyring
func (k *Keyring) Fetch(key KeyType) (string, error) {
	get, err := k.ring.Get(string(key))
	if err != nil {
		return "", err
	}

	return string(get.Data), err
}

// InMemoryKeyring : store and fetch secrets in memory for tests
type InMemoryKeyring struct {
	krPkg.Keyring
	item krPkg.Item
}

var _ krPkg.Keyring = &InMemoryKeyring{}

// NewInMemoryKeyring : Constructor for in memory keyring
func NewInMemoryKeyring() *InMemoryKeyring {
	return &InMemoryKeyring{}
}

// Get : get a value from the in memory keyring
func (i *InMemoryKeyring) Get(key string) (krPkg.Item, error) {
	return i.item, nil
}

// Set : set a value in the in memory keyring
func (i *InMemoryKeyring) Set(item krPkg.Item) error {
	i.item = item

	return nil
}
