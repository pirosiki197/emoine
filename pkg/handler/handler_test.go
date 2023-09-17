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

func TestStream(t *testing.T) {
	server := setUpServer()
	t.Cleanup(server.Close)

	client := apiv1connect.NewAPIServiceClient(server.Client(), server.URL)
	ctx := context.Background()

	// コメントを受け取るイベントを作成
	res, err := client.CreateEvent(ctx, connect.NewRequest(&apiv1.CreateEventRequest{
		Title: "test",
	}))
	if err != nil {
		t.Fatal(err)
	}
	id := res.Msg.Id
	// コメントを受け取らないイベントを作成
	res, err = client.CreateEvent(ctx, connect.NewRequest(&apiv1.CreateEventRequest{
		Title: "this event will not received",
	}))
	if err != nil {
		t.Fatal(err)
	}
	anotherID := res.Msg.Id

	t.Run("receive", func(t *testing.T) {
		t.Parallel()
		defer server.CloseClientConnections()
		stream, err := client.ConnectToStream(ctx, connect.NewRequest(&apiv1.ConnectToStreamRequest{
			EventId: id,
		}))
		if err != nil {
			t.Logf("error: %+v", err)
			t.Fatal(err)
		}

		go func() {
			time.Sleep(3 * time.Second)
			server.CloseClientConnections()
		}()

		for stream.Receive() {
			if c := stream.Msg().GetComment(); c != nil {
				if c.Text != "test" {
					t.Errorf("got %s, want %s", c.Text, "test")
				}
			}
		}
	})

	t.Run("send", func(t *testing.T) {
		t.Parallel()
		// wait for stream
		time.Sleep(1 * time.Second)
		// Textがtestのコメントのみが受け取られる
		comments := []*apiv1.SendCommentRequest{
			{
				EventId: id,
				UserId:  "test",
				Text:    "test",
			},
			{
				EventId: anotherID,
				UserId:  "test",
				Text:    "will not received",
			},
		}

		for _, c := range comments {
			_, err := client.SendComment(ctx, connect.NewRequest(c))
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}
