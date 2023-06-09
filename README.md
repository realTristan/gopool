# gopool ![Stars](https://img.shields.io/github/stars/realTristan/gopool?color=brightgreen) ![Watchers](https://img.shields.io/github/watchers/realTristan/gopool?label=Watchers)
![image-6](https://user-images.githubusercontent.com/75189508/234116253-eec9af68-66c5-44a4-8c30-cf82dbf936c1.png)

# Example
```go
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
	// Replace nil with your client
	var client *gp.Client[any] = gp.NewClient(nil)

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
}
```

# License
MIT License

Copyright (c) 2023 Tristan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
