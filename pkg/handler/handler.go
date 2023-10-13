package handler

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"github.com/pirosiki197/emoine/pkg/domain"
	apiv1 "github.com/pirosiki197/emoine/pkg/proto/api/v1"
	"github.com/pirosiki197/emoine/pkg/proto/pbconv"
	"github.com/pirosiki197/emoine/pkg/pubsub"
	"github.com/samber/lo"
)

type handler struct {
	repo      Repository
	validator *protovalidate.Validator

	sm *streamManager
}

func NewHandlre(repo Repository) *handler {
	v, err := protovalidate.New()
	if err != nil {
		panic(err)
	}

	sm := &streamManager{
		ps: pubsub.NewPubsub(),
	}
	go sm.run(context.Background())

	return &handler{
		repo:      repo,
		validator: v,
		sm:        sm,
	}
}

func (h *handler) CreateEvent(
	ctx context.Context,
	req *connect.Request[apiv1.CreateEventRequest],
) (*connect.Response[apiv1.CreateEventResponse], error) {
	if err := h.validator.Validate(req.Msg); err != nil {
		slog.Error(err.Error(), slog.Any("req", req.Msg))
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	e := &domain.Event{
		ID:      uuid.New(),
		Title:   req.Msg.Title,
		StartAt: req.Msg.StartAt.AsTime(),
		EndAt:   req.Msg.EndAt.AsTime(),
	}

	if err := h.repo.CreateEvent(ctx, e); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&apiv1.CreateEventResponse{
		Id: e.ID.String(),
	})

	return res, nil
}

func (h *handler) GetEvents(
	ctx context.Context,
	req *connect.Request[apiv1.GetEventsRequest],
) (*connect.Response[apiv1.GetEventsResponse], error) {
	events, err := h.repo.GetEvents(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&apiv1.GetEventsResponse{
		Events: lo.Map(events, func(e domain.Event, _ int) *apiv1.Event {
			return pbconv.FromEventModel(&e)
		}),
	})

	return res, nil
}

func (h *handler) SendComment(
	ctx context.Context,
	req *connect.Request[apiv1.SendCommentRequest],
) (*connect.Response[apiv1.SendCommentResponse], error) {
	if err := h.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	c := &domain.Comment{
		ID:        uuid.New(),
		UserID:    req.Msg.UserId,
		EventID:   uuid.MustParse(req.Msg.EventId),
		Text:      req.Msg.Text,
		CreatedAt: time.Now(),
	}

	if err := h.repo.CreateComment(ctx, c); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := h.sm.sendComment(c); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&apiv1.SendCommentResponse{
		Id: c.ID.String(),
	})

	return res, nil
}

func (h *handler) GetComments(
	ctx context.Context,
	req *connect.Request[apiv1.GetCommentsRequest],
) (*connect.Response[apiv1.GetCommentsResponse], error) {
	if err := h.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	comments, err := h.repo.GetEventComments(ctx, req.Msg.EventId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&apiv1.GetCommentsResponse{
		Comments: lo.Map(comments, func(c domain.Comment, _ int) *apiv1.Comment {
			return pbconv.FromCommentModel(&c)
		}),
	})

	return res, nil
}

func (h *handler) ConnectToStream(
	ctx context.Context,
	req *connect.Request[apiv1.ConnectToStreamRequest],
	stream *connect.ServerStream[apiv1.ConnectToStreamResponse],
) error {
	if err := h.validator.Validate(req.Msg); err != nil {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	e, err := h.repo.GetEvent(ctx, req.Msg.EventId)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	if e == nil {
		return connect.NewError(connect.CodeNotFound, nil)
	}

	stream.Send(&apiv1.ConnectToStreamResponse{
		EventOrComment: &apiv1.ConnectToStreamResponse_Event{
			Event: pbconv.FromEventModel(e),
		},
	})

	client := h.sm.connectToStream(e.ID)
	defer h.sm.disconnectFromStream(client)

	for {
		select {
		case msg := <-client.receive():
			stream.Send(&apiv1.ConnectToStreamResponse{
				EventOrComment: &apiv1.ConnectToStreamResponse_Comment{
					Comment: msg.comment,
				},
			})
		case <-ctx.Done():
			return nil
		}
	}
}
