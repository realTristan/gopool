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
	Enabled bool
	Active  bool
	client  *Client[any]
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
		Enabled: true,
		Active:  false,
		client:  client,
	}
}

// Disable a connection in the connection pool
func (p *Pool) Disable(conn *Connection) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.connections[conn]; !ok {
		return errors.New("connection does not exist")
	}
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
		if !v["active"] && v["enabled"] {
			return conn, nil
		}
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

// Get the connection client
func (conn *Connection) WithClient(fn func(c Client[any]) any) any {
	return fn(*conn.client)
}

// Execute a function with a pool connection
func (p *Pool) WithConnection(conn *Connection, fn func(c Connection) any) any {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Execute the function
	return fn(*conn)
}
