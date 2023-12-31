package handler

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/pirosiki197/emoine/pkg/domain"
	apiv1 "github.com/pirosiki197/emoine/pkg/infra/proto/api/v1"
	"github.com/pirosiki197/emoine/pkg/infra/proto/pbconv"
)

type streamManager struct {
	cps     domain.PubSub[*domain.Comment]
	clients []*client

	mu sync.Mutex
}

func (sm *streamManager) run(ctx context.Context) {
	// 今はコメントのみだけど増えるかも
	sub, stop := sm.cps.Subscribe(ctx)
	defer stop()
	for {
		select {
		case msg := <-sub:
			if msg.Err != nil {
				slog.Error("failed to receive comment", slog.String("err", msg.Err.Error()))
				continue
			}

			// broadcast
			sm.mu.Lock() // you can't use defer here
			for _, c := range sm.clients {
				if c.eventID == msg.Msg.EventID {
					c.ch <- message{
						comment: pbconv.FromCommentModel(msg.Msg),
					}
				}
			}
			sm.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (sm *streamManager) sendComment(comment *domain.Comment) error {
	return sm.cps.Publish(context.Background(), comment)
}

func (sm *streamManager) connectToStream(eventID uuid.UUID) *client {
	ch := make(chan message)
	c := &client{
		eventID: eventID,
		ch:      ch,
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.clients = append(sm.clients, c)

	return c
}

func (sm *streamManager) disconnectFromStream(c *client) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.clients = slices.DeleteFunc(sm.clients, func(d *client) bool { return d == c })
}

type client struct {
	eventID uuid.UUID
	ch      chan message
}

func (c *client) receive() <-chan message {
	return c.ch
}

type message struct {
	comment *apiv1.Comment
}
