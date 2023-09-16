package handler

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"connectrpc.com/connect"
	apiv1 "github.com/pirosiki197/emoine/pkg/proto/api/v1"
	"github.com/pirosiki197/emoine/pkg/proto/api/v1/apiv1connect"
	"github.com/pirosiki197/emoine/pkg/repository"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func setUpServer() *httptest.Server {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	repo := repository.NewRepository(db)
	h := NewHandlre(repo)
	mux := http.NewServeMux()
	mux.Handle(apiv1connect.NewAPIServiceHandler(h))
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	return server
}

func TestCreateEvent(t *testing.T) {
	server := setUpServer()
	t.Cleanup(server.Close)

	client := apiv1connect.NewAPIServiceClient(server.Client(), server.URL)
	ctx := context.Background()

	cases := []struct {
		name string
		req  *connect.Request[apiv1.CreateEventRequest]
		want error
	}{
		{
			name: "success",
			req: connect.NewRequest(&apiv1.CreateEventRequest{
				Title:   "test",
				StartAt: timestamppb.New(time.Now()),
				EndAt:   timestamppb.New(time.Now()),
			}),
			want: nil,
		},
		{
			name: "success with empty endAt",
			req: connect.NewRequest(&apiv1.CreateEventRequest{
				Title:   "test",
				StartAt: timestamppb.New(time.Now()),
				EndAt:   nil,
			}),
			want: nil,
		},
		{
			name: "invalid argument with empty title",
			req: connect.NewRequest(&apiv1.CreateEventRequest{
				Title: "",
			}),
			want: &connect.Error{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := client.CreateEvent(ctx, c.req)
			if reflect.TypeOf(err) != reflect.TypeOf(c.want) {
				t.Errorf("got %v, want %v", reflect.TypeOf(err), reflect.TypeOf(c.want))
			}
		})
	}
}

func TestGetEvents(t *testing.T) {
	server := setUpServer()
	t.Cleanup(server.Close)

	client := apiv1connect.NewAPIServiceClient(server.Client(), server.URL)
	ctx := context.Background()

	events := []*apiv1.Event{
		{
			Title:   "test 1",
			StartAt: timestamppb.New(time.Now()),
			EndAt:   timestamppb.New(time.Now()),
		},
		{
			Title:   "test 2",
			StartAt: timestamppb.New(time.Now().Add(-time.Hour)),
			EndAt:   nil,
		},
		{
			Title:   "test 3",
			StartAt: timestamppb.New(time.Now().Add(-time.Hour * 2)),
			EndAt:   nil,
		},
	}
	for _, e := range events {
		_, err := client.CreateEvent(ctx, connect.NewRequest(&apiv1.CreateEventRequest{
			Title:   e.Title,
			StartAt: e.StartAt,
			EndAt:   e.EndAt,
		}))
		if err != nil {
			t.Fatal(err)
		}
	}

	res, err := client.GetEvents(ctx, connect.NewRequest(&apiv1.GetEventsRequest{}))
	if err != nil {
		t.Fatal(err)
	}
	for i, e := range events {
		if e.Title != res.Msg.Events[i].Title {
			t.Errorf("got %v, want %v", res.Msg.Events[i].Title, e.Title)
		}
	}
}
