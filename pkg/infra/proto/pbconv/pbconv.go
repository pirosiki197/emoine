package pbconv

import (
	"time"

	"github.com/google/uuid"
	"github.com/pirosiki197/emoine/pkg/domain"
	apiv1 "github.com/pirosiki197/emoine/pkg/infra/proto/api/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromEventModel(e *domain.Event) *apiv1.Event {
	return &apiv1.Event{
		Id:      e.ID.String(),
		Title:   e.Title,
		StartAt: timestamppb.New(e.StartAt),
		EndAt:   lo.Ternary(e.EndAt.IsZero(), nil, timestamppb.New(e.EndAt)),
	}
}

func ToEventModel(e *apiv1.Event) *domain.Event {
	return &domain.Event{
		ID:      uuid.MustParse(e.Id),
		Title:   e.Title,
		StartAt: ToTime(e.StartAt),
		EndAt:   lo.Ternary(e.EndAt == nil, time.Time{}, ToTime(e.EndAt)),
	}
}

func FromCommentModel(c *domain.Comment) *apiv1.Comment {
	return &apiv1.Comment{
		Id:        c.ID.String(),
		UserId:    c.UserID,
		EventId:   c.EventID.String(),
		Text:      c.Text,
		CreatedAt: timestamppb.New(c.CreatedAt),
	}
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func ToTime(t *timestamppb.Timestamp) time.Time {
	return t.AsTime().In(jst)
}
