package gopool

import (
	"errors"
	"sync"
	"time"
)

// Connection Pool Struct *Connection
type Pool struct {
	maxSize     int
	currentSize int
	mutex       *sync.RWMutex
	connections *Queue
}

// Initialize the Connection Pool
func InitPool(maxSize int) *Pool {
	return &Pool{
		maxSize:     maxSize,
		currentSize: 0,
		mutex:       &sync.RWMutex{},
		connections: &Queue{
			mutex: &sync.RWMutex{},
			items: []*Connection{},
		},
	}
}

// Get the current pool size
func (p *Pool) Size() (int, int) {
	var (
		copyCurSize int = p.currentSize
		copyMaxSize int = p.maxSize
	)
	return copyCurSize, copyMaxSize
}

// Add a new connection to the connection pool
func (p *Pool) Add(client *Client, expire int64) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Check if the maximum pool connection has been reached
	if p.currentSize+1 > p.maxSize {
		return errors.New("maximum pool capacity reached")
	}

	// If the connection already exists in the pool
	if p.connections.exists(client) {
		return errors.New("client already exists in connection pool")
	}

	// If the provided expiration is greater than 0
	// aka the user wants to set an expiration...
	if expire > 0 {
		expire = time.Now().Unix() + expire
	}

	// Append the connection to the existing pool connections
	p.connections.add(&Connection{
		client: client,
		expire: expire,
	})

	// Increase the current pool size
	p.currentSize++

	// Return no errors
	return nil
}

// Get a connection from the connection pool
func (p *Pool) get() (*Connection, error) {
	var conn *Connection = p.connections.next()
	if conn == nil {
		return nil, errors.New("no connections are available")
	}

	// If the connection is not active and has expired
	if conn.expire > 0 && conn.expire-time.Now().UnixMilli() <= 0 {
		if err := p.connections.delete(conn); err == nil {
			p.currentSize--
		}
		return p.get()
	}

	// Return the connection and no error
	return conn, nil
}

// Continuously fetch for a connection until the timeout time has been reached
func (p *Pool) getTimeout(timeout int64) (*Connection, error) {
	var (
		start int64       = time.Now().UnixMilli()
		conn  *Connection = nil
		err   error       = nil
	)
	for conn == nil || err != nil {
		// Check if the provided timeout has been reached
		if timeout > 0 && start+timeout > time.Now().UnixMilli() {
			return nil, errors.New("timeout reached. no available connections could be found")
		}

		// Get the connection
		conn, err = p.get()
	}
	return conn, nil
}

// Execute a function with a pool connection
func (p *Pool) WithConnection(fn func(c *Connection) any) (any, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Get a connection from the connections pool
	if conn, err := p.get(); err != nil {
		return nil, err
	} else {
		// Add the connection back to the pool once function returns
		defer p.connections.add(conn)

		// Execute the function
		return fn(conn), nil
	}
}

// Execute a function with a pool connection
func (p *Pool) WithConnectionTimeout(timeout int64, fn func(c Connection) any) (any, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Get a connection from the connections pool
	if conn, err := p.getTimeout(timeout); err != nil {
		return nil, err
	} else {
		// Add the connection back to the pool once function returns
		defer p.connections.add(conn)

		// Execute the function
		return fn(*conn), nil
	}
}
