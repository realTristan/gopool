package gopool

import (
	"errors"
	"sync"
)

// Item Queue Struct
type Queue struct {
	mutex *sync.RWMutex
	items []*Connection
}

// Get an item and remove it from the queue
func (q *Queue) next() *Connection {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// If there's no items in the queue
	if len(q.items) == 0 {
		return nil
	}

	// Get the item from the queue
	var item *Connection = q.items[0]

	// If the queue size is just 1
	if len(q.items) == 1 {
		q.items = []*Connection{}
		return item
	}

	// Else, set the queue items to all those beyond the first index
	q.items = q.items[1:]

	// Return the item
	return item
}

// Check if a client in the pool already exists
func (q *Queue) exists(c *Client) bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	// Iterate over the queue
	for i := 0; i < q.size(); i++ {
		if q.items[i].client == c {
			return true
		}
	}
	return false
}

// Add an item to the queue
func (q *Queue) add(item *Connection) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Append the new item to the end of the queue items
	q.items = append(q.items, item)
}

// Get the size of the queue
func (q *Queue) size() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	// Return the length of the item queue
	return len(q.items)
}

// Delete an item from the queue
func (q *Queue) delete(item *Connection) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Iterate over the queue items
	for i := 0; i < q.size(); i++ {
		if q.items[i] == item {
			q.items = append(q.items[i:], q.items[i+1:]...)
			return nil
		}
	}

	// Return error
	return errors.New("unable to delete item. it does not exist")
}
