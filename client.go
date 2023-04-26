package gopool

// Connection Client Struct
type Client[T any] *T

// Initialize a new client
func NewClient[T any](c *T) *Client[T] {
	var client *Client[T] = new(Client[T])
	*client = c
	return client
}
