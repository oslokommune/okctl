// Package breeze provides a wrapper around storm
// for opening and closing the database as required
// instead of keeping the connection open over
// and unlimited timespan. This means storage
// operations will take longer, but we won't
// block other actions.
package breeze
