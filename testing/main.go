package main

import (
	"fmt"

	gp "github.com/realTristan/gopool"
)

//
/*

Replace [string] with your client type

*/
//
func main() {
	// Initialize a pool
	// Max pool size of 4
	var pool *gp.Pool[string] = gp.InitPool[string](4)

	// Initalize a client
	var client *gp.Client[string] = gp.NewClient[string]("Hello World!")

	// Add the connection to the pool
	// Expire in 10 seconds, -1 for no expiration
	// When the connection expires, close the client
	pool.Add(client, 10, func(client *gp.Client[string]) {
		// client.close()
	})

	// Access a connection from the pool
	pool.WithConnection(func(conn gp.Connection[string], opts *gp.Options[string]) any {
		// Use the connection client
		conn.WithClient(func(client gp.Client[string]) any {
			// await client. (... whatever you're trying to do with your database client)
			return nil
		})
		return nil
	})

	// Print the pool size
	fmt.Println(pool.Size())

	// Access a connection from the pool with a connection timeout
	var timeout int64 = 10000 // 10000 milliseconds till timeout (10 second)
	pool.WithConnectionTimeout(timeout, func(conn gp.Connection[string], opts *gp.Options[string]) any {
		// Use the connection client
		opts.DeferDelete = true
		conn.WithClient(func(client gp.Client[string]) any {
			// await client. (... whatever you're trying to do with your database client)
			return nil
		})
		return nil
	})
	fmt.Println(pool.Size())
}
