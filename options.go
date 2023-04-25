package gopool

type Options struct {
	ExpiresAt func(conn *Connection) int64
}
