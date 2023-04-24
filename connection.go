package gopool

// Connection struct
type Connection struct {
	enabled bool
	active  bool
	client  *Client[any]
}
