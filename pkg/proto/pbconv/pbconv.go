package pbconv

import (
	"github.com/pirosiki197/emoine/pkg/model"
	apiv1 "github.com/pirosiki197/emoine/pkg/proto/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromEventModel(e *model.Event) *apiv1.Event {
	return &apiv1.Event{
		Id:      e.ID.String(),
		Title:   e.Title,
		StartAt: timestamppb.New(e.StartAt),
		EndAt:   timestamppb.New(e.EndAt),
	}
}

func FromCommentModel(c *model.Comment) *apiv1.Comment {
	return &apiv1.Comment{
		Id:        c.ID.String(),
		UserId:    c.UserID,
		EventId:   c.EventID.String(),
		Text:      c.Text,
		CreatedAt: timestamppb.New(c.CreatedAt),
	}
}
