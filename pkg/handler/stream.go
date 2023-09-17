package handler

import (
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/pirosiki197/emoine/pkg/model"
	apiv1 "github.com/pirosiki197/emoine/pkg/proto/api/v1"
	"github.com/pirosiki197/emoine/pkg/proto/pbconv"
)

type streamManager struct {
	clients []*client

	mu sync.Mutex
}

func (sm *streamManager) sendComment(comment *model.Comment) {
	for _, c := range sm.clients {
		if c.eventID == comment.EventID {
			c.ch <- message{
				comment: pbconv.FromCommentModel(comment),
			}
		}
	}
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

	sm.clients = slices.DeleteFunc(sm.clients, func(d *client) bool {
		return d == c
	})
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
