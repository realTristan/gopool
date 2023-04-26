package gopool

type Options[T any] struct {
	DeferDelete    bool
	DeferSetExpire int64
	ExpiresAt      func(conn Connection[T]) int64
}
