package gopool

type Options[T any] struct {
	DeferDelete    bool
	DeferSetExpire int64
}
