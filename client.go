package gopool

// Connection Client Struct
type Client[T any] any

// Initialize a new client
func NewClient[T any](c any) *Client[T] {
	var client *Client[T] = new(Client[T])
	*client = c
	return client
}
