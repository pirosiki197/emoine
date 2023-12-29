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

func TestCreate_GetEvents(t *testing.T) {
	s := NewAPIServer(t)
	t.Cleanup(s.Close)
	client := s.Client()
	ctx := context.Background()

	events := []*apiv1.Event{
		{
			Title:   "test event 1",
			StartAt: timestamppb.New(time.Now()),
		},
		{
			Title:   "test event 2",
			StartAt: timestamppb.New(time.Now().Add(-1 * time.Hour)),
			EndAt:   timestamppb.New(time.Now()),
		},
	}

	t.Run("create events", func(t *testing.T) {
		for _, e := range events {
			req := connect.NewRequest(&apiv1.CreateEventRequest{
				Title:   e.Title,
				StartAt: e.StartAt,
				EndAt:   e.EndAt,
			})
			res, err := client.CreateEvent(ctx, req)
			if err != nil {
				t.Fatal(err)
			}
			e.Id = res.Msg.Id
		}
	})

	t.Run("get events", func(t *testing.T) {
		req := connect.NewRequest(&apiv1.GetEventsRequest{})
		res, err := client.GetEvents(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Msg.Events) != len(events) {
			t.Fatalf("want %d, got %d", len(events), len(res.Msg.Events))
		}
		gotEvents := res.Msg.Events
		slices.SortFunc(gotEvents, func(a, b *apiv1.Event) int { return strings.Compare(a.Id, b.Id) })
		slices.SortFunc(events, func(a, b *apiv1.Event) int { return strings.Compare(a.Id, b.Id) })
		for i, e := range events {
			if !EventEqual(e, gotEvents[i]) {
				t.Errorf("want %v, got %v", e, gotEvents[i])
			}
		}
	})
}

func TestCreate_GetEvent(t *testing.T) {
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
			EndAt:   event.EndAt,
		})
		res, err := client.CreateEvent(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		event.Id = res.Msg.Id
	})

	t.Run("get event", func(t *testing.T) {
		req := connect.NewRequest(&apiv1.GetEventRequest{Id: event.Id})
		res, err := client.GetEvent(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if !EventEqual(event, res.Msg.Event) {
			t.Errorf("want %v, got %v", event, res.Msg.Event)
		}
	})
}
