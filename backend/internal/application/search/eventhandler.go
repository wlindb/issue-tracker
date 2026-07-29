package search

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
	"github.com/wlindb/issue-tracker/internal/pkg/tracker/model"
)

type EventOption func(*SearchEventHandler) error

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
	issueDocumentRepository issuedocument.IssueDocumentRepository
}

func NewSearchEventHandler(issueDocumentRepository issuedocument.IssueDocumentRepository, opts ...EventOption) (SearchEventHandler, error) {
	var zero SearchEventHandler
	handler := SearchEventHandler{issueDocumentRepository: issueDocumentRepository}
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

	document := issuedocument.NewIssueDocument(event.Payload.ID, event.Payload.Title, description, event.OccurredAt)
	if _, err := h.issueDocumentRepository.Create(ctx, document); err != nil {
		slog.Error("create issue document", "error", err.Error())
		return fmt.Errorf("create issue document: %w", err)
	}

	return nil
}

func (h SearchEventHandler) HandleIssueTitleUpdated(_ context.Context, event model.IssueTitleUpdatedEvent) error {
	slog.Info("issue title updated",
		"issue_id", event.Payload.ID,
		"title", event.Payload.Title,
		"occurred_at", event.OccurredAt,
	)

	return nil
}

func (h SearchEventHandler) HandleIssueDescriptionUpdated(_ context.Context, event model.IssueDescriptionUpdatedEvent) error {
	slog.Info("issue description updated",
		"issue_id", event.Payload.ID,
		"occurred_at", event.OccurredAt,
	)

	return nil
}
