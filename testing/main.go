package main

import gp "github.com/realTristan/gopool"

// Main function for testing/examples
func main() {
	// Initialize a pool
	// Max pool size of 4
	var pool *gp.Pool = gp.InitPool(4)

	// Initalize a client
	var client *gp.Client = nil // whatever your client is

	// Add the connection to the pool
	// Expire in 10 seconds, -1 for no expiration
	pool.Add(client, 10)

	// Access a connection from the pool
	pool.WithConnection(func(conn *gp.Connection) any {
		// Use the connection client
		conn.WithClient(func(client gp.Client) any {
			// await client. (... whatever you're trying to do with your database client)
			return nil
		})
		return nil
	})

	// Access a connection from the pool with a connection timeout
	var timeout int64 = 1000 // 1000 milliseconds till timeout (1 second)
	pool.WithConnectionTimeout(timeout, func(conn gp.Connection) any {
		// Use the connection client
		conn.WithClient(func(client gp.Client) any {
			// await client. (... whatever you're trying to do with your database client)
			return nil
		})
		return nil
	})
}
