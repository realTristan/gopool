package gopool

type Options struct {
	ExpiresIn func(conn *Connection) int64
}
