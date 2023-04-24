package gopool

import "errors"

// Options for WithConnection() function
type Options struct {
	Delete    func(p *Pool, conn *Connection) error
	Enable    func(p *Pool, conn *Connection) error
	Disable   func(p *Pool, conn *Connection) error
	IsEnabled func(conn *Connection) bool
}

// Disable a connection in the connection pool
func disable(p *Pool, conn *Connection) error {
	// Check if the provided connection already exists
	if !p.exists(conn) {
		return errors.New("connection does not exist")
	}

	// Check if the connection has already been disabled
	if !conn.enabled {
		return errors.New("connection already disabled")
	}

	// Disable the connection
	conn.enabled = false
	return nil
}

// Enable a connection in the connection pool
func enable(p *Pool, conn *Connection) error {
	// Check if the provided connection already exists
	if !p.exists(conn) {
		return errors.New("connection does not exist")
	}

	// Check if the connection has already been enabled
	if conn.enabled {
		return errors.New("connection already enabled")
	}

	// Enable the connection
	conn.enabled = true
	return nil
}

// Delete a connection from the connection pool
func delete(p *Pool, conn *Connection) error {
	if !p.exists(conn) {
		return errors.New("connection does not exist")
	}
	for i := 0; i < len(p.connections); i++ {
		if p.connections[i] == conn {
			p.connections = append(p.connections[:i], p.connections[i+1:]...)
			break
		}
	}
	return nil
}

// Check if the connection is enabled
func isEnabled(conn *Connection) bool {
	return conn.enabled
}
