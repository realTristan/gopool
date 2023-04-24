package gopool

// Connection struct
type Connection struct {
	enabled bool
	active  bool
	client  *Client[any]
	expire  int64
}

// Set a connection activity. Define in a function to enable defering
func (conn *Connection) setConnectionActivity(activity bool) {
	conn.active = activity
}

// Get the connection client
func (conn *Connection) WithClient(fn func(c *Client[any]) any) any {
	return fn(conn.client)
}
