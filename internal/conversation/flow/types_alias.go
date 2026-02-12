package flow

import "github.com/memohai/memoh/internal/conversation"

type ModelMessage = conversation.ModelMessage
type ContentPart = conversation.ContentPart
type ToolCall = conversation.ToolCall
type ToolCallFunction = conversation.ToolCallFunction
type AssistantOutput = conversation.AssistantOutput
type ChatRequest = conversation.ChatRequest
type ChatResponse = conversation.ChatResponse
type StreamChunk = conversation.StreamChunk
type Settings = conversation.Settings

var NewTextContent = conversation.NewTextContent
