package gopool

import (
	"sync"
)

// Item Queue Struct
type ConnectionQueue[T any] struct {
	mutex       *sync.RWMutex
	connections []*Connection[T]
}

// Get an item and remove it from the queue
func (q *ConnectionQueue[T]) next() *Connection[T] {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// If there's no items in the queue
	if len(q.connections) == 0 {
		return nil
	}

	// Get the item from the queue
	var conn *Connection[T] = q.connections[0]

	// If the queue size is just 1
	if len(q.connections) == 1 {
		q.connections = []*Connection[T]{}
		return conn
	}

	// Else, set the queue items to all those beyond the first index
	q.connections = q.connections[1:]
	return conn
}

// Check if a client already exists in the pool
func (q *ConnectionQueue[T]) clientExists(c *Client[T]) bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	// Iterate over the queue
	for i := 0; i < len(q.connections); i++ {
		if q.connections[i].client == c {
			return true
		}
	}
	return false
}

// Add an item to the queue
func (q *ConnectionQueue[T]) add(item *Connection[T]) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Append the new item to the end of the queue items
	q.connections = append(q.connections, item)
}

/* Delete an item from the queue
func (q *ConnectionQueue[T]) delete(item *Connection[T]) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Iterate over the queue items
	for i := 0; i < len(q.connections); i++ {
		if q.connections[i] == item {
			q.connections = append(q.connections[:i], q.connections[i+1:]...)
			return nil
		}
	}

	// Return error
	return errors.New("unable to delete item. it does not exist")
}
*/
