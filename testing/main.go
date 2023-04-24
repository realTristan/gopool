package main

import gp "github.com/realTristan/gopool"

// Main function for testing/examples
func main() {
	// Initialize a pool
	var pool *gp.Pool = gp.InitPool()

	// Initalize a connection
	var client *gp.Client[any] = nil
	var conn = gp.InitConn(client)

	// Add the connection to the pool
	pool.Add(conn)
}
