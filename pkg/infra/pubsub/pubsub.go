package pubsub

import (
	"context"
	"sync"

	"github.com/go-json-experiment/json"
	"github.com/pirosiki197/emoine/pkg/domain"
	"github.com/redis/go-redis/v9"
)

type Pubsub struct {
	rdb *redis.Client
}

func NewPubsub() *Pubsub {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	return &Pubsub{
		rdb: rdb,
	}
}

const (
	CommentChannel = "comments"
)

func (p *Pubsub) PublishComment(ctx context.Context, c *domain.Comment) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return p.rdb.Publish(ctx, CommentChannel, b).Err()
}

// SubscribeComment returns a channel that receives comments.
// The channel is closed when the close function is called.
// When ctx is canceled, the channel is closed, so you don't need to call the close function.
//
// If an error occurs, the error is set to the Err field of the Message and keep receiving.
func (p *Pubsub) SubscribeComment(ctx context.Context) (sub <-chan domain.Message[domain.Comment], stop func()) {
	pubsub := p.rdb.Subscribe(ctx, CommentChannel)
	ch := make(chan domain.Message[domain.Comment])

	stop = sync.OnceFunc(func() {
		close(ch)
		_ = pubsub.Close()
	})

	redisCh := pubsub.Channel()

	// start goroutine to stop receiving comments
	go func() {
		<-ctx.Done()
		stop()
	}()
	// start goroutine to receive comments
	go func() {
		for c := range redisCh {
			var msg domain.Message[domain.Comment]
			comment := &domain.Comment{}

			err := json.Unmarshal([]byte(c.Payload), comment)
			if err != nil {
				ch <- msg.SetErr(err)
				continue
			}

			select {
			case ch <- msg.SetMsg(comment):
			case <-ctx.Done():
				stop()
				return
			}
		}
	}()

	return ch, stop
}
