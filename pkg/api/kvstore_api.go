package api

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var storeNameRe = regexp.MustCompile(`[a-zA-Z0-9_.-]+`)

const (
	// minStoreNameLength defines the minimum length of a key / value store name
	minStoreNameLength = 3
	// maxStoreNameLength defines the maximum length of a key / value store name
	maxStoreNameLength = 255
)

// StoreName defines a key-value store name
type StoreName string

// Validate ensures the name is legal
func (n StoreName) Validate() error {
	return validation.Validate(string(n),
		validation.Length(minStoreNameLength, maxStoreNameLength),
		validation.Match(storeNameRe),
	)
}

// CreateStoreOpts defines required data for creating a table
type CreateStoreOpts struct {
	// ClusterID identifies which cluster the store should be created in context of
	ClusterID ID
	// Name defines the name of the key / value store to create
	Name StoreName
	// Keys defines the provisioned keys the store should accept
	Keys []string
}

// Validate ensures opts contains the required and correct data
func (o CreateStoreOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterID, validation.Required),
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.Keys, validation.Required, validation.Length(1, 0)),
	)
}

// DeleteStoreOpts defines required data for deleting a table
type DeleteStoreOpts struct {
	// ClusterID identifies which cluster the store should be deleted in context of
	ClusterID ID
	// Name defines the name of the key / value store to delete
	Name StoreName
}

// Validate ensures opts contains the required and correct data
func (o DeleteStoreOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterID, validation.Required),
		validation.Field(&o.Name, validation.Required),
	)
}

// InsertItemOpts defines required data for adding a table
type InsertItemOpts struct {
	// ClusterID identifies which cluster the store exists in context of
	ClusterID ID
	// TableName defines the name of the key / value store to insert an item in
	TableName StoreName
	// Item defines the item to insert
	Item KeyValueStoreItem
}

// Validate ensures opts contains the required and correct data
func (o InsertItemOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterID, validation.Required),
		validation.Field(&o.TableName, validation.Required),
		validation.Field(&o.Item, validation.Required),
	)
}

// DeleteItemOpts defines required data for deleting an item
type DeleteItemOpts struct {
	// ClusterID identifies which cluster the item should be deleted in context of
	ClusterID ID
	// TableName defines the name of the key / value store to delete an item from
	TableName StoreName
	// Field defines the field of where the key identifying the item should be queried
	Field string
	// Key defines the key identifying the item to delete
	Key string
}

// Validate ensures opts contains the required and correct data
func (o DeleteItemOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterID, validation.Required),
		validation.Field(&o.TableName, validation.Required),
		validation.Field(&o.Field, validation.Required),
		validation.Field(&o.Key, validation.Required),
	)
}

// ItemSelector defines required data for selecting an item
type ItemSelector struct {
	// Key defines a key to match with a value when selecting an item
	Key string
	// Value defines the value a certain key must have which can identify an item
	Value string
}

// GetStringOpts defines required data for retrieving a string
type GetStringOpts struct {
	// ClusterID identifies which cluster to do the operation in context of
	ClusterID ID
	// TableName defines the name of the key / value store to do the operation in context of
	TableName StoreName
	// Selector defines the necessary information to identify an item in the store
	Selector ItemSelector
	// Field defines which attribute of an item to return
	Field string
}

// Validate ensures opts contains the required and correct data
func (o GetStringOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterID, validation.Required),
		validation.Field(&o.TableName, validation.Required),
		validation.Field(&o.Selector, validation.Required),
		validation.Field(&o.Field, validation.Required),
	)
}

// KeyValueStoreItem represents one item in a key value store
type KeyValueStoreItem struct {
	Fields map[string]interface{}
}

// Storer defines operations done with stores
type Storer interface {
	// CreateStore defines a function that creates a key / value store
	CreateStore(CreateStoreOpts) error
	// DeleteStore defines a function that deletes a key / value store
	DeleteStore(DeleteStoreOpts) error
	// ListStores defines a function that lists all available stores
	ListStores() ([]string, error)
}

// Stringer defines operations done with string items
type Stringer interface {
	// GetString defines a function that retrieves a string from a key / value store
	GetString(opts GetStringOpts) (string, error)
}

// Itemer defines operations done with items of tables
type Itemer interface {
	Stringer
	// InsertItem defines a function that inserts a new item in a key / value store
	InsertItem(opts InsertItemOpts) error
	// RemoveItem defines a function that removes an item from a key / value store
	RemoveItem(opts DeleteItemOpts) error
}

// KeyValueStoreService defines operations on a key-value database service
type KeyValueStoreService interface {
	Storer
	Itemer
}

// KeyValueStoreCloudProvider defines operations on a key-value database provider
type KeyValueStoreCloudProvider interface {
	Storer
	Itemer
}
