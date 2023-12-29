package integration

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	apiv1 "github.com/pirosiki197/emoine/pkg/infra/proto/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreate_GetComments(t *testing.T) {
	s := NewAPIServer(t)
	t.Cleanup(s.Close)
	client := s.Client()
	ctx := context.Background()

	event := &apiv1.Event{
		Title:   "test event 1",
		StartAt: timestamppb.New(time.Now()),
	}
	t.Run("create event", func(t *testing.T) {
		req := connect.NewRequest(&apiv1.CreateEventRequest{
			Title:   event.Title,
			StartAt: event.StartAt,
		})
		res, err := client.CreateEvent(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		event.Id = res.Msg.Id
	})

	comments := []*apiv1.Comment{
		{
			UserId:  "test user 1",
			EventId: event.Id,
			Text:    "test comment 1",
		},
		{
			UserId:  "test user 2",
			EventId: event.Id,
			Text:    "test comment 2",
		},
	}

	t.Run("create comments", func(t *testing.T) {
		for _, c := range comments {
			req := connect.NewRequest(&apiv1.SendCommentRequest{
				UserId:  c.UserId,
				EventId: c.EventId,
				Text:    c.Text,
			})
			res, err := client.SendComment(ctx, req)
			if err != nil {
				t.Fatal(err)
			}
			c.Id = res.Msg.Id
			c.CreatedAt = timestamppb.Now()
		}
	})

	t.Run("get comments", func(t *testing.T) {
		req := connect.NewRequest(&apiv1.GetCommentsRequest{
			EventId: event.Id,
		})
		res, err := client.GetComments(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Msg.Comments) != len(comments) {
			t.Fatalf("want %d, got %d", len(comments), len(res.Msg.Comments))
		}
		gotComments := res.Msg.Comments
		slices.SortFunc(gotComments, func(a, b *apiv1.Comment) int { return strings.Compare(a.Id, b.Id) })
		slices.SortFunc(comments, func(a, b *apiv1.Comment) int { return strings.Compare(a.Id, b.Id) })
		for i, c := range comments {
			if !CommentEqual(c, gotComments[i]) {
				t.Fatalf("want %v, got %v", c, gotComments[i])
			}
		}
	})
}
