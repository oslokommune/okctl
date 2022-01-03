package breeze

import (
	"fmt"
	"os"

	"github.com/asdine/storm/v3/index"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/json"
)

// Breeze contains the required state
type Breeze struct {
	path     *string
	addend   []string
	writable *bool
}

// configurer defines functionality for configuring
type configurer interface {
	// SetDatabaseFilePath indicates where Breeze should look for a database file
	SetDatabaseFilePath(string)
	// SetWritable enables mutating operations
	SetWritable(bool)
	// IsWritable returns true if Breeze is allowed to mutate data
	IsWritable() bool
}

// Client extracts the interface we are currently using from storm
// we can simply extend this interface in the future if we want to
// use more of the available functionality
type Client interface {
	configurer
	All(to interface{}, options ...func(*index.Options)) error
	AllByIndex(fieldName string, to interface{}, options ...func(*index.Options)) error
	DeleteStruct(data interface{}) error
	From(addend ...string) Client
	Init(data interface{}) error
	One(fieldName string, value interface{}, to interface{}) error
	Save(data interface{}) error
}

// New returns an initialised client
func New() Breeze {
	path := ""
	writable := false

	return Breeze{path: &path, writable: &writable}
}

// From returns the client with the added addend, we
// use this later on when using the storm From
func (b Breeze) From(addend ...string) Client {
	from := b
	from.addend = append(from.addend, addend...)

	return from
}

// Init the data type
func (b Breeze) Init(data interface{}) error {
	db, node, err := b.open()
	if err != nil {
		return err
	}

	defer func() {
		err = db.Close()
	}()

	return node.Init(data)
}

// AllByIndex fetches all items by using the index
func (b Breeze) AllByIndex(fieldName string, to interface{}, options ...func(*index.Options)) error {
	db, node, err := b.open()
	if err != nil {
		return err
	}

	defer func() {
		err = db.Close()
	}()

	return node.AllByIndex(fieldName, to, options...)
}

// All returns all items
func (b Breeze) All(to interface{}, options ...func(*index.Options)) error {
	db, node, err := b.open()
	if err != nil {
		return err
	}

	defer func() {
		err = db.Close()
	}()

	return node.All(to, options...)
}

// One returns one item
func (b Breeze) One(fieldName string, value interface{}, to interface{}) error {
	db, node, err := b.open()
	if err != nil {
		return err
	}

	defer func() {
		err = db.Close()
	}()

	return node.One(fieldName, value, to)
}

// Save the provided data
func (b Breeze) Save(data interface{}) error {
	db, node, err := b.open()
	if err != nil {
		return err
	}

	defer func() {
		err = db.Close()
	}()

	return node.Save(data)
}

// DeleteStruct removes the data
func (b Breeze) DeleteStruct(data interface{}) error {
	db, node, err := b.open()
	if err != nil {
		return err
	}

	defer func() {
		err = db.Close()
	}()

	return node.DeleteStruct(data)
}

// SetDatabaseFilePath configures the path to where breeze will look for the database file when doing operations
func (b Breeze) SetDatabaseFilePath(path string) {
	*b.path = path
}

// SetWritable sets a flag enabling or disabling database mutation
func (b Breeze) SetWritable(writable bool) {
	*b.writable = writable
}

// IsWritable returns the flag representing if the database can be mutated
func (b Breeze) IsWritable() bool {
	return *b.writable
}

// open the storm database
func (b Breeze) open() (*storm.DB, storm.Node, error) {
	if *b.path == "" {
		return nil, nil, fmt.Errorf("no database path is set")
	}

	var fileMode os.FileMode = 0o400

	if *b.writable {
		fileMode = 0o600
	}

	db, err := storm.Open(*b.path, storm.Codec(json.Codec), storm.BoltOptions(fileMode, nil))
	if err != nil {
		return nil, nil, fmt.Errorf("loading state database: %w", err)
	}

	if b.addend != nil {
		return db, db.From(b.addend...), nil
	}

	return db, db, nil
}
