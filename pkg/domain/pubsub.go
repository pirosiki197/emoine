package domain

import "context"

type PubSub[T StreamObject] interface {
	Publish(context.Context, T) error
	Subscribe(context.Context) (sub <-chan Message[T], stop func())
}

type Message[T any] struct {
	Msg T
	Err error
}

func (m Message[T]) SetMsg(msg T) Message[T] {
	m.Msg = msg
	return m
}

func (m Message[T]) SetErr(err error) Message[T] {
	m.Err = err
	return m
}
