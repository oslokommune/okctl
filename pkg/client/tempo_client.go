package client

// TempoState describes an interface to handle state for a Tempo instance
type TempoState interface {
	HasTempo() (bool, error)
}
