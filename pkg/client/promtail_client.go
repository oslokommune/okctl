package client

// PromtailState describes an interface to handle state for a Loki instance
type PromtailState interface {
	HasPromtail() (bool, error)
}
