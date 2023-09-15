package pbconv

import (
	apiv1 "github.com/pirosiki197/emoine/pkg/gen/api/v1"
	"github.com/pirosiki197/emoine/pkg/model"
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
