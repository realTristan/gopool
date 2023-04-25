package gopool

type Options[T any] struct {
	ExpiresAt func(conn *Connection[T]) int64
}
