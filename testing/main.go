package main

import (
	"fmt"

	gp "github.com/realTristan/gopool"
)

// Your client
type MyClient struct{}

//
/*

Replace [MyClient] with your client type

*/
//
func main() {
	// Initialize a pool
	// Max pool size of 4
	var pool *gp.Pool[MyClient] = gp.InitPool[MyClient](4)

	// Initalize a client
	var client *gp.Client[MyClient] = gp.NewClient(&MyClient{})

	// Add the connection to the pool
	// Expire in 10 seconds, -1 for no expiration
	// When the connection expires, close the client
	pool.Add(client, 10, func(client *gp.Client[MyClient]) {
		// client.close()
	})

	// Access a connection from the pool
	pool.WithConnection(func(conn gp.Connection[MyClient], opts *gp.Options[MyClient]) any {
		fmt.Println(conn.ExpiresAt())

		// Add 10 seconds to the expiration
		opts.DeferSetExpire = conn.ExpiresAt() + 10

		// Use the connection client
		conn.WithClient(func(client gp.Client[MyClient]) any {
			// await client. (... whatever you're trying to do with your database client)
			return nil
		})
		return nil
	})

	// Print the pool size
	fmt.Println(pool.Size())

	// Access a connection from the pool with a connection timeout
	var timeout int64 = 10000 // 10000 milliseconds till timeout (10 second)
	pool.WithConnectionTimeout(timeout, func(conn gp.Connection[MyClient], opts *gp.Options[MyClient]) any {
		fmt.Println(conn.ExpiresAt())

		// Delete the connection once finished
		opts.DeferDelete = true

		// Use the connection client
		conn.WithClient(func(client gp.Client[MyClient]) any {
			// await client. (... whatever you're trying to do with your database client)
			return nil
		})
		return nil
	})
	fmt.Println(pool.Size())
}
