package gopool

import (
	"errors"
	"sync"
	"time"
)

// Connection Client Struct
type Client[T any] *any

// Connection Pool Struct
type Pool struct {
	connections map[*Connection]map[string]bool
	mutex       *sync.RWMutex
}

// Connection struct
type Connection struct {
	enabled bool
	active  bool
	client  *Client[any]
	mutex   *sync.RWMutex
}

// Initialize the Connection Pool
func InitPool() *Pool {
	return &Pool{
		connections: make(map[*Connection]map[string]bool),
		mutex:       &sync.RWMutex{},
	}
}

// Initialize a Connection
func InitConn(client *Client[any]) *Connection {
	return &Connection{
		enabled: true,
		active:  false,
		client:  client,
		mutex:   &sync.RWMutex{},
	}
}

// Get the client from the connection
func (c *Connection) Client() *Client[any] {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.client
}

// Get whether the current connection is active
func (c *Connection) IsActive() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.active
}

// Get whether the current connection is enabled
func (c *Connection) IsEnabled() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.enabled
}

// Disable a connection in the connection pool
func (p *Pool) Disable(conn *Connection) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.connections[conn]; !ok {
		return errors.New("connection does not exist")
	}
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	p.connections[conn]["enabled"] = false
	return nil
}

// Enable a connection in the connection pool
func (p *Pool) Enable(conn *Connection) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.connections[conn]; !ok {
		return errors.New("connection does not exist")
	}
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	p.connections[conn]["enabled"] = true
	return nil
}

// Add a connection to the connection pool
func (p *Pool) Add(conn *Connection) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Check if the connection already exists
	if _, ok := p.connections[conn]; ok {
		return errors.New("connection already exists")
	}
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	p.connections[conn] = make(map[string]bool)
	return nil
}

// Remove a connection from the connection pool
func (p *Pool) Remove(conn *Connection) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.connections[conn]; !ok {
		return errors.New("connection does not exist")
	}
	delete(p.connections, conn)
	return nil
}

// Get a connection from the connection pool
/*

var (
	conn *Connection = nil
	err error = nil
)

for conn == nil || err != nil {
	conn, err = p.Get()
}

*/
func (p *Pool) Get() (*Connection, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for conn, v := range p.connections {
		conn.mutex.RLock()
		if !v["active"] && v["enabled"] {
			conn.mutex.RUnlock()
			return conn, nil
		}
		conn.mutex.RUnlock()
	}
	return nil, errors.New("no connections are available")
}

// Continuously fetch for a connection until the timeout time has been reached
func (p *Pool) GetTimeout(timeout int64) (*Connection, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var (
		start int64 = time.Now().UnixMilli()
		conn  *Connection
		err   error
	)
	for conn == nil || err != nil {
		conn, err = p.Get()
		if start+timeout > time.Now().UnixMilli() {
			return nil, errors.New("timeout reached. no available connection could be fetched")
		}
	}
	return conn, nil
}
