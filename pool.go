package gopool

import (
	"errors"
	"sync"
	"time"
)

// Connection Pool Struct
type Pool struct {
	connections []*Connection
	mutex       *sync.RWMutex
}

// Initialize the Connection Pool
func InitPool() *Pool {
	return &Pool{
		connections: []*Connection{},
		mutex:       &sync.RWMutex{},
	}
}

// Add a new connection to the connection pool
func (p *Pool) New(client *Client[any]) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Create a new connection
	var conn *Connection = &Connection{
		enabled: true,
		active:  false,
		client:  client,
	}

	// Append the connection to the existing pool connections
	p.connections = append(p.connections, conn)
	return nil
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

// Check if a connection exists in the connection pool
func (p *Pool) exists(conn *Connection) bool {
	for i := 0; i < len(p.connections); i++ {
		if p.connections[i] == conn {
			return true
		}
	}
	return false
}

// Get a connection from the connection pool
func (p *Pool) get() (*Connection, error) {
	for i := 0; i < len(p.connections); i++ {
		var conn *Connection = p.connections[i]
		if !conn.active && conn.enabled {
			return conn, nil
		}
	}
	return nil, errors.New("no connections are available")
}

// Continuously fetch for a connection until the timeout time has been reached
func (p *Pool) getTimeout(timeout int64) (*Connection, error) {
	var (
		start int64 = time.Now().UnixMilli()
		conn  *Connection
		err   error
	)
	for conn == nil || err != nil {
		conn, err = p.get()
		if start+timeout > time.Now().UnixMilli() {
			return nil, errors.New("timeout reached. no available connection could be fetched")
		}
	}
	return conn, nil
}
