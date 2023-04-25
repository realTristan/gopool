package gopool

// Connection struct
type Connection[T any] struct {
	client   *Client[T]
	expire   int64
	onExpire func(client *Client[T])
}

// Get the connection client
func (conn *Connection[T]) WithClient(fn func(c Client[T]) any) any {
	return fn(*conn.client)
}
