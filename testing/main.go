package main

import (
	gp "github.com/realTristan/gopool"
)

// Main function for testing/examples
func main() {
	// Initialize a pool
	var pool *gp.Pool = gp.InitPool()

	// Initalize a connection
	var client *gp.Client[any] = nil // whatever your client is

	// Add the connection to the pool
	// Expire in 10 seconds, -1 for no expiration
	pool.New(client, 10)

	// Access a connection from the pool
	var _, _ = pool.WithConnection(func(conn *gp.Connection, opts *gp.Options) any {
		// Example to disable connection
		if err := opts.Disable(pool, conn); err != nil {
			return err
		}

		// Use the connection client
		conn.WithClient(func(client *gp.Client[any]) any {
			// await client. (... whatever you're trying to do with your database client)
			return nil
		})
		return nil
	})

	// Access a connection from the pool with a connection timeout
	var timeout int64 = 1000 // 1000 milliseconds till timeout (1 second)
	var _, _ = pool.WithConnectionTimeout(timeout, func(conn *gp.Connection, opts *gp.Options) any {
		// Example to disable connection
		if err := opts.Disable(pool, conn); err != nil {
			return err
		}

		// Use the connection client
		conn.WithClient(func(client *gp.Client[any]) any {
			// await client. (... whatever you're trying to do with your database client)
			return nil
		})
		return nil
	})
}
