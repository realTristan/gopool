package gopool

// Get the connection client
func (conn *Connection) WithClient(fn func(c *Client[any]) any) any {
	return fn(conn.client)
}

// Execute a function with a pool connection
func (p *Pool) WithConnection(fn func(c *Connection, opts *Options) any) (any, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Get a connection from the connections pool
	if conn, err := p.get(); err != nil {
		return nil, err
	} else {
		// Set the connection to active
		conn.setConnectionActivity(true)
		defer conn.setConnectionActivity(false)

		// Execute the function
		return fn(conn, &Options{
			Delete:    delete,
			Enable:    enable,
			Disable:   disable,
			IsEnabled: isEnabled,
		}), nil
	}
}

// Execute a function with a pool connection
func (p *Pool) WithConnectionTimeout(timeout int64, fn func(c *Connection, opts *Options) any) (any, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Get a connection from the connections pool
	if conn, err := p.getTimeout(timeout); err != nil {
		return nil, err
	} else {
		// Set the connection to active
		conn.setConnectionActivity(true)
		defer conn.setConnectionActivity(false)

		// Execute the function
		return fn(conn, &Options{
			Delete:    delete,
			Enable:    enable,
			Disable:   disable,
			IsEnabled: isEnabled,
		}), nil
	}
}
