package pubsub

import (
	"context"
	"sync"

	"github.com/go-json-experiment/json"
	"github.com/pirosiki197/emoine/pkg/domain"
	"github.com/redis/go-redis/v9"
)

// TODO: channelの名前を持つ
type Pubsub[T domain.StreamObject] struct {
	rdb         *redis.Client
	channelName string
}

func NewPubsub[T domain.StreamObject](rdb *redis.Client) *Pubsub[T] {
	var object T
	return &Pubsub[T]{
		rdb:         rdb,
		channelName: object.Type().String(),
	}
}

const (
	CommentChannel = "comments"
)

func (p *Pubsub[T]) Publish(ctx context.Context, c T) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	// switch case
	return p.rdb.Publish(ctx, p.channelName, b).Err()
}

// Subscribe returns a channel that receives comments.
// The channel is closed when the close function is called.
// When ctx is canceled, the channel is closed, so you don't need to call the close function.
//
// If an error occurs, the error is set to the Err field of the Message and keep receiving.
func (p *Pubsub[T]) Subscribe(ctx context.Context) (sub <-chan domain.Message[T], stop func()) {
	// switch channel
	pubsub := p.rdb.Subscribe(ctx, p.channelName)
	ch := make(chan domain.Message[T])

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
			var msg domain.Message[T]
			value := new(T)

			err := json.Unmarshal([]byte(c.Payload), value)
			if err != nil {
				ch <- msg.SetErr(err)
				continue
			}

			select {
			case ch <- msg.SetMsg(*value):
			case <-ctx.Done():
				stop()
				return
			}
		}
	}()

	return ch, stop
}
