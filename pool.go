package gopool

import (
	"errors"
	"sync"
	"time"
)

// Connection Pool Struct *Connection
type Pool[T any] struct {
	maxConnections     int
	currentConnections int
	mutex              *sync.RWMutex
	connections        *ConnectionQueue[T]
}

// Initialize the Connection Pool
func InitPool[T any](maxSize int) *Pool[T] {
	return &Pool[T]{
		maxConnections:     maxSize,
		currentConnections: 0,
		mutex:              &sync.RWMutex{},
		connections: &ConnectionQueue[T]{
			mutex:       &sync.RWMutex{},
			connections: []*Connection[T]{},
		},
	}
}

// Get the current pool size
func (p *Pool[T]) Size() int {
	var copy int = p.currentConnections
	return copy
}

// Add a new connection to the connection pool
func (p *Pool[T]) Add(client *Client[T], expire int64, onExpire func(client *Client[T])) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Check if the maximum pool connection has been reached
	if p.currentConnections+1 > p.maxConnections {
		return errors.New("maximum pool capacity reached")
	}

	// If the connection already exists in the pool
	if p.connections.clientExists(client) {
		return errors.New("client already exists in connection pool")
	}

	// If the provided expiration is greater than 0
	// aka the user wants to set an expiration...
	if expire > 0 {
		expire = time.Now().Unix() + expire
	}

	// Append the connection to the existing pool connections
	p.connections.add(&Connection[T]{
		client:   client,
		expire:   expire,
		onExpire: onExpire,
	})

	// Increase the current pool size
	p.currentConnections++

	// Return no errors
	return nil
}

// Get a connection from the connection pool
func (p *Pool[T]) get() (*Connection[T], error) {
	var conn *Connection[T] = p.connections.next()
	if conn == nil {
		return nil, errors.New("no connections are available")
	}

	// If the connection is not active and has expired
	if conn.expire > 0 && conn.expire-time.Now().Unix() <= 0 {
		if err := p.connections.delete(conn); err == nil {
			p.mutex.Lock()
			p.currentConnections--
			p.mutex.Unlock()
			conn.onExpire(conn.client)
		}
		return p.get()
	}

	// Return the connection and no error
	return conn, nil
}

// Continuously fetch for a connection until the timeout time has been reached
func (p *Pool[T]) getTimeout(timeout int64) (*Connection[T], error) {
	var (
		start int64          = time.Now().UnixMilli()
		conn  *Connection[T] = nil
		err   error          = nil
	)
	for conn == nil || err != nil {
		// Check if the provided timeout has been reached
		if timeout > 0 && start+timeout > time.Now().UnixMilli() {
			return nil, errors.New("timeout reached. no available connections could be found")
		}

		// Get a new connection
		conn, err = p.get()
	}
	return conn, nil
}

// Execute a function with a pool connection
func (p *Pool[T]) WithConnection(fn func(c Connection[T], opts *Options[T]) any) (any, error) {
	// Get a connection from the connections pool
	if conn, err := p.get(); err != nil {
		return nil, err
	} else {
		// Initialize the options
		var opts *Options[T] = &Options[T]{
			DeferDelete:    false,
			DeferSetExpire: -2,
		}

		// Add the connection back to the pool once function returns
		defer func() {
			// If the user wants to delete the connection on defer
			// then they can set the DeferDelete option to true
			if !opts.DeferDelete {
				p.connections.add(conn)

				// If the user wants to update the connection expiration
				// on defer, they can set the DeferSetExpire option to
				// a number greater than -2 (i.e -1 for no expiration or > 0 for an expiration)
				if opts.DeferSetExpire > -2 {
					conn.expire = opts.DeferSetExpire
				}
			}
		}()

		// Execute the function
		return fn(*conn, opts), nil
	}
}

// Execute a function with a pool connection
func (p *Pool[T]) WithConnectionTimeout(timeout int64, fn func(c Connection[T], opts *Options[T]) any) (any, error) {
	// Get a connection from the connections pool
	if conn, err := p.getTimeout(timeout); err != nil {
		return nil, err
	} else {
		// Initialize the options
		var opts *Options[T] = &Options[T]{
			DeferDelete:    false,
			DeferSetExpire: -2,
		}

		// Add the connection back to the pool once function returns
		defer func() {
			// If the user wants to delete the connection on defer,
			// they can set the DeferDelete option to true
			if !opts.DeferDelete {
				p.connections.add(conn)

				// If the user wants to update the connection expiration
				// on defer, they can set the DeferSetExpire option to
				// a number greater than -2 (i.e -1 for no expiration or > 0 for an expiration)
				if opts.DeferSetExpire > -2 {
					conn.expire = opts.DeferSetExpire
				}
			}
		}()

		// Execute the function
		return fn(*conn, opts), nil
	}
}
