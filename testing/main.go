package main

import gp "github.com/realTristan/gopool"

//
/*

Replace [any] with your client type

*/
//
func main() {
	// Initialize a pool
	// Max pool size of 4
	var pool *gp.Pool[any] = gp.InitPool[any](4)

	// Initalize a client
	var client *gp.Client[any] = nil // whatever your client is

	// Add the connection to the pool
	// Expire in 10 seconds, -1 for no expiration
	// When the connection expires, close the client
	pool.Add(client, 10, func(client *gp.Client[any]) {
		// client.close()
	})

	// Access a connection from the pool
	pool.WithConnection(func(conn gp.Connection[any], opts *gp.Options[any]) any {
		// Use the connection client
		conn.WithClient(func(client gp.Client[any]) any {
			// await client. (... whatever you're trying to do with your database client)
			return nil
		})
		return nil
	})

	// Access a connection from the pool with a connection timeout
	var timeout int64 = 1000 // 1000 milliseconds till timeout (1 second)
	pool.WithConnectionTimeout(timeout, func(conn gp.Connection[any], opts *gp.Options[any]) any {
		// Use the connection client
		conn.WithClient(func(client gp.Client[any]) any {
			// await client. (... whatever you're trying to do with your database client)
			return nil
		})
		return nil
	})
}
