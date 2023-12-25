package pubsub_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pirosiki197/emoine/pkg/domain"
	"github.com/pirosiki197/emoine/pkg/pubsub"
)

func TestPubSub_Comment(t *testing.T) {
	p := pubsub.NewPubsub()
	err := p.PublishComment(context.Background(), &domain.Comment{
		ID:        uuid.New(),
		UserID:    "user-id",
		EventID:   uuid.New(),
		Text:      "text",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
}

func TestPubSub_Comment_Concurrency(t *testing.T) {
	p := pubsub.NewPubsub()

	t.Run("publish", func(t *testing.T) {
		t.Parallel()
		// wait for subscribe
		time.Sleep(1 * time.Second)
		for i := 0; i < 1000; i++ {
			err := p.PublishComment(context.Background(), &domain.Comment{
				ID:        uuid.New(),
				UserID:    "user-id",
				EventID:   uuid.New(),
				Text:      "text",
				CreatedAt: time.Now(),
			})
			if err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("subscribe", func(t *testing.T) {
		t.Parallel()
		ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()
		ch, _ := p.SubscribeComment(ctx)
		var count int
		for m := range ch {
			if m.Err != nil {
				t.Error(m.Err)
			}
			t.Log(m.Msg)
			count++
		}
		if count != 1000 {
			t.Errorf("want %d, got %d", 1000, count)
		}
	})
}
