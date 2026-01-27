package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/memohai/memoh/internal/db/sqlc"
)

// SelectMemoryModel selects a chat model for memory operations.
func SelectMemoryModel(ctx context.Context, modelsService *Service, queries *sqlc.Queries) (GetResponse, sqlc.LlmProvider, error) {
	// First try to get the memory-enabled model.
	memoryModel, err := modelsService.GetByEnableAs(ctx, EnableAsMemory)
	if err == nil {
		provider, err := FetchProviderByID(ctx, queries, memoryModel.LlmProviderID)
		if err != nil {
			return GetResponse{}, sqlc.LlmProvider{}, err
		}
		return memoryModel, provider, nil
	}

	// Fallback to chat model.
	chatModel, err := modelsService.GetByEnableAs(ctx, EnableAsChat)
	if err == nil {
		provider, err := FetchProviderByID(ctx, queries, chatModel.LlmProviderID)
		if err != nil {
			return GetResponse{}, sqlc.LlmProvider{}, err
		}
		return chatModel, provider, nil
	}

	// If no enabled models, try to find any chat model.
	candidates, err := modelsService.ListByType(ctx, ModelTypeChat)
	if err != nil || len(candidates) == 0 {
		return GetResponse{}, sqlc.LlmProvider{}, fmt.Errorf("no chat models available for memory operations")
	}

	selected := candidates[0]
	provider, err := FetchProviderByID(ctx, queries, selected.LlmProviderID)
	if err != nil {
		return GetResponse{}, sqlc.LlmProvider{}, err
	}
	return selected, provider, nil
}

// FetchProviderByID fetches a provider by ID.
func FetchProviderByID(ctx context.Context, queries *sqlc.Queries, providerID string) (sqlc.LlmProvider, error) {
	if strings.TrimSpace(providerID) == "" {
		return sqlc.LlmProvider{}, fmt.Errorf("provider id missing")
	}
	parsed, err := parseUUID(providerID)
	if err != nil {
		return sqlc.LlmProvider{}, err
	}
	return queries.GetLlmProviderByID(ctx, parsed)
}
