package domain

import "context"

type PubSub interface {
	PublishComment(context.Context, *Comment) error
	SubscribeComment(context.Context) (sub <-chan Message[Comment], stop func())
}

type Message[T any] struct {
	Msg *T
	Err error
}

func (m Message[T]) SetMsg(msg *T) Message[T] {
	m.Msg = msg
	return m
}

func (m Message[T]) SetErr(err error) Message[T] {
	m.Err = err
	return m
}
