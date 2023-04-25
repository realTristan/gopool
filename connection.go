package gopool

// Connection struct
type Connection struct {
	client *Client
	expire int64
}

// Get the connection client
func (conn *Connection) WithClient(fn func(c Client) any) any {
	return fn(*conn.client)
}
