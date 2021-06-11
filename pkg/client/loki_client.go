package client

// LokiState describes an interface to handle state for a Loki instance
type LokiState interface {
	HasLoki() (bool, error)
}
