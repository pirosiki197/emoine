package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	apiv1 "github.com/pirosiki197/emoine/pkg/gen/api/v1"
)

type handler struct {
	repo Repository
}

func NewHandlre(repo Repository) *handler {
	return &handler{
		repo: repo,
	}
}

func (s *handler) CreateEvent(
	ctx context.Context,
	req *connect.Request[apiv1.CreateEventRequest],
) (*connect.Response[apiv1.CreateEventResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (s *handler) GetEvents(
	ctx context.Context,
	req *connect.Request[apiv1.GetEventsRequest],
) (*connect.Response[apiv1.GetEventsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (s *handler) SendComment(
	ctx context.Context,
	req *connect.Request[apiv1.SendCommentRequest],
) (*connect.Response[apiv1.SendCommentResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (s *handler) GetComments(
	ctx context.Context,
	req *connect.Request[apiv1.GetCommentsRequest],
) (*connect.Response[apiv1.GetCommentsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (s *handler) ConnectToStream(
	ctx context.Context,
	req *connect.Request[apiv1.ConnectToStreamRequest],
	stream *connect.ServerStream[apiv1.ConnectToStreamResponse],
) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}
