package search

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
	"github.com/wlindb/issue-tracker/internal/pkg/tracker/model"
)

type EventOption func(*SearchEventHandler) error

// IssueDocumentService is what the handler needs from the domain.
type IssueDocumentService interface {
	Create(ctx context.Context, id uuid.UUID, title string, description string, updatedAt time.Time) error
	UpdateTitle(ctx context.Context, id uuid.UUID, title string) error
	UpdateDescription(ctx context.Context, id uuid.UUID, description string) error
}

func WithIssueCreated(subscriber event.SubscriberOf[model.IssueCreatedEvent]) EventOption {
	return func(h *SearchEventHandler) error {
		return subscriber.Subscribe(h.HandleIssueCreated)
	}
}

func WithIssueTitleUpdated(subscriber event.SubscriberOf[model.IssueTitleUpdatedEvent]) EventOption {
	return func(h *SearchEventHandler) error {
		return subscriber.Subscribe(h.HandleIssueTitleUpdated)
	}
}

func WithIssueDescriptionUpdated(subscriber event.SubscriberOf[model.IssueDescriptionUpdatedEvent]) EventOption {
	return func(h *SearchEventHandler) error {
		return subscriber.Subscribe(h.HandleIssueDescriptionUpdated)
	}
}

type SearchEventHandler struct {
	issueDocumentService IssueDocumentService
}

func NewSearchEventHandler(issueDocumentService IssueDocumentService, opts ...EventOption) (SearchEventHandler, error) {
	var zero SearchEventHandler
	handler := SearchEventHandler{issueDocumentService: issueDocumentService}
	for _, opt := range opts {
		if err := opt(&handler); err != nil {
			return zero, fmt.Errorf("search event handler option: %w", err)
		}
	}
	return handler, nil
}

func (h SearchEventHandler) HandleIssueCreated(ctx context.Context, event model.IssueCreatedEvent) error {
	var description string
	if event.Payload.Description != nil {
		description = *event.Payload.Description
	}

	if err := h.issueDocumentService.Create(ctx, event.Payload.ID, event.Payload.Title, description, event.OccurredAt); err != nil {
		slog.Error("create issue document", "error", err.Error())
		return err //nolint:wrapcheck // domain service already returns a wrapped, contextual error
	}

	return nil
}

func (h SearchEventHandler) HandleIssueTitleUpdated(ctx context.Context, event model.IssueTitleUpdatedEvent) error {
	if err := h.issueDocumentService.UpdateTitle(ctx, event.Payload.ID, event.Payload.Title); err != nil {
		slog.Error("update issue document title", "error", err.Error())
		return fmt.Errorf("handle issue title updated: %w", err)
	}

	return nil
}

func (h SearchEventHandler) HandleIssueDescriptionUpdated(ctx context.Context, event model.IssueDescriptionUpdatedEvent) error {
	var description string
	if event.Payload.Description != nil {
		description = *event.Payload.Description
	}

	if err := h.issueDocumentService.UpdateDescription(ctx, event.Payload.ID, description); err != nil {
		slog.Error("update issue document description", "error", err.Error())
		return fmt.Errorf("handle issue description updated: %w", err)
	}

	return nil
}
