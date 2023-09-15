package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	apiv1 "github.com/pirosiki197/emoine/pkg/gen/api/v1"
	"github.com/pirosiki197/emoine/pkg/model"
	"github.com/pirosiki197/emoine/pkg/pbconv"
	"github.com/samber/lo"
)

type handler struct {
	repo      Repository
	validator *protovalidate.Validator
}

func NewHandlre(repo Repository) *handler {
	v, err := protovalidate.New()
	if err != nil {
		panic(err)
	}

	return &handler{
		repo:      repo,
		validator: v,
	}
}

func (s *handler) CreateEvent(
	ctx context.Context,
	req *connect.Request[apiv1.CreateEventRequest],
) (*connect.Response[apiv1.CreateEventResponse], error) {
	if err := s.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	e := &model.Event{
		ID:      uuid.New(),
		Title:   req.Msg.Title,
		StartAt: req.Msg.StartAt.AsTime(),
		EndAt:   req.Msg.EndAt.AsTime(),
	}

	if err := s.repo.CreateEvent(ctx, e); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&apiv1.CreateEventResponse{
		Id: e.ID.String(),
	})

	return res, nil
}

func (s *handler) GetEvents(
	ctx context.Context,
	req *connect.Request[apiv1.GetEventsRequest],
) (*connect.Response[apiv1.GetEventsResponse], error) {
	events, err := s.repo.GetEvents(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&apiv1.GetEventsResponse{
		Events: lo.Map(events, func(e model.Event, _ int) *apiv1.Event {
			return pbconv.FromEventModel(&e)
		}),
	})

	return res, nil
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
