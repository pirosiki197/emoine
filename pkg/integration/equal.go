package integration

import (
	"time"

	apiv1 "github.com/pirosiki197/emoine/pkg/infra/proto/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EventEqual(a, b *apiv1.Event) bool {
	if a.Id != b.Id {
		return false
	}
	if a.Title != b.Title {
		return false
	}
	if !TimestamppbNeallyEqual(a.StartAt, b.StartAt) {
		return false
	}
	if !(a.EndAt == nil && b.EndAt == nil) && !TimestamppbNeallyEqual(a.EndAt, b.EndAt) {
		return false
	}
	return true
}

func CommentEqual(a, b *apiv1.Comment) bool {
	if a.Id != b.Id {
		return false
	}
	if a.UserId != b.UserId {
		return false
	}
	if a.EventId != b.EventId {
		return false
	}
	if a.Text != b.Text {
		return false
	}
	if !TimestamppbNeallyEqual(a.CreatedAt, b.CreatedAt) {
		return false
	}
	return true
}

func TimestamppbNeallyEqual(a, b *timestamppb.Timestamp) bool {
	return a.AsTime().Sub(b.AsTime()).Abs() < 1*time.Second
}
