package message

import (
	"context"
	"errors"
	"testing"

	"github.com/memohai/memoh/internal/channel"
	mcpgw "github.com/memohai/memoh/internal/mcp"
)

type fakeSender struct {
	err error
}

func (f *fakeSender) Send(ctx context.Context, botID string, channelType channel.ChannelType, req channel.SendRequest) error {
	return f.err
}

type fakeResolver struct {
	ct  channel.ChannelType
	err error
}

func (f *fakeResolver) ParseChannelType(raw string) (channel.ChannelType, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.ct, nil
}

func TestExecutor_ListTools_NilDeps(t *testing.T) {
	exec := NewExecutor(nil, nil, nil)
	tools, err := exec.ListTools(context.Background(), mcpgw.ToolSessionContext{})
	if err != nil {
		t.Fatal(err)
	}
	if len(tools) != 0 {
		t.Errorf("expected 0 tools when deps nil, got %d", len(tools))
	}
}

func TestExecutor_ListTools(t *testing.T) {
	sender := &fakeSender{}
	resolver := &fakeResolver{ct: channel.ChannelType("feishu")}
	exec := NewExecutor(nil, sender, resolver)
	tools, err := exec.ListTools(context.Background(), mcpgw.ToolSessionContext{})
	if err != nil {
		t.Fatal(err)
	}
	if len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools))
	}
	if tools[0].Name != toolSendMessage {
		t.Errorf("tool name = %q, want %q", tools[0].Name, toolSendMessage)
	}
}

func TestExecutor_CallTool_NotFound(t *testing.T) {
	sender := &fakeSender{}
	resolver := &fakeResolver{ct: channel.ChannelType("feishu")}
	exec := NewExecutor(nil, sender, resolver)
	_, err := exec.CallTool(context.Background(), mcpgw.ToolSessionContext{BotID: "bot1"}, "other_tool", nil)
	if err != mcpgw.ErrToolNotFound {
		t.Errorf("expected ErrToolNotFound, got %v", err)
	}
}

func TestExecutor_CallTool_NilDeps(t *testing.T) {
	exec := NewExecutor(nil, nil, nil)
	result, err := exec.CallTool(context.Background(), mcpgw.ToolSessionContext{BotID: "bot1"}, toolSendMessage, map[string]any{
		"platform": "feishu", "target": "t1", "text": "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if isErr, _ := result["isError"].(bool); !isErr {
		t.Error("expected error result when deps nil")
	}
}

func TestExecutor_CallTool_NoBotID(t *testing.T) {
	sender := &fakeSender{}
	resolver := &fakeResolver{ct: channel.ChannelType("feishu")}
	exec := NewExecutor(nil, sender, resolver)
	result, err := exec.CallTool(context.Background(), mcpgw.ToolSessionContext{}, toolSendMessage, map[string]any{
		"platform": "feishu", "target": "t1", "text": "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if isErr, _ := result["isError"].(bool); !isErr {
		t.Error("expected error when bot_id is missing")
	}
}

func TestExecutor_CallTool_BotIDMismatch(t *testing.T) {
	sender := &fakeSender{}
	resolver := &fakeResolver{ct: channel.ChannelType("feishu")}
	exec := NewExecutor(nil, sender, resolver)
	session := mcpgw.ToolSessionContext{BotID: "bot1"}
	result, err := exec.CallTool(context.Background(), session, toolSendMessage, map[string]any{
		"bot_id": "other", "platform": "feishu", "target": "t1", "text": "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if isErr, _ := result["isError"].(bool); !isErr {
		t.Error("expected error when bot_id mismatch")
	}
}

func TestExecutor_CallTool_NoPlatform(t *testing.T) {
	sender := &fakeSender{}
	resolver := &fakeResolver{ct: channel.ChannelType("feishu")}
	exec := NewExecutor(nil, sender, resolver)
	session := mcpgw.ToolSessionContext{BotID: "bot1"}
	result, err := exec.CallTool(context.Background(), session, toolSendMessage, map[string]any{
		"target": "t1", "text": "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if isErr, _ := result["isError"].(bool); !isErr {
		t.Error("expected error when platform is missing")
	}
}

func TestExecutor_CallTool_PlatformParseError(t *testing.T) {
	sender := &fakeSender{}
	resolver := &fakeResolver{err: errors.New("unknown platform")}
	exec := NewExecutor(nil, sender, resolver)
	session := mcpgw.ToolSessionContext{BotID: "bot1", CurrentPlatform: "feishu"}
	result, err := exec.CallTool(context.Background(), session, toolSendMessage, map[string]any{
		"platform": "bad", "target": "t1", "text": "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if isErr, _ := result["isError"].(bool); !isErr {
		t.Error("expected error when platform parse fails")
	}
}

func TestExecutor_CallTool_NoMessage(t *testing.T) {
	sender := &fakeSender{}
	resolver := &fakeResolver{ct: channel.ChannelType("feishu")}
	exec := NewExecutor(nil, sender, resolver)
	session := mcpgw.ToolSessionContext{BotID: "bot1"}
	result, err := exec.CallTool(context.Background(), session, toolSendMessage, map[string]any{
		"platform": "feishu", "target": "t1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if isErr, _ := result["isError"].(bool); !isErr {
		t.Error("expected error when message/text is missing")
	}
}

func TestExecutor_CallTool_NoTarget(t *testing.T) {
	sender := &fakeSender{}
	resolver := &fakeResolver{ct: channel.ChannelType("feishu")}
	exec := NewExecutor(nil, sender, resolver)
	session := mcpgw.ToolSessionContext{BotID: "bot1"}
	result, err := exec.CallTool(context.Background(), session, toolSendMessage, map[string]any{
		"platform": "feishu", "text": "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if isErr, _ := result["isError"].(bool); !isErr {
		t.Error("expected error when target and channel_identity_id are missing")
	}
}

func TestExecutor_CallTool_SendError(t *testing.T) {
	sender := &fakeSender{err: errors.New("send failed")}
	resolver := &fakeResolver{ct: channel.ChannelType("feishu")}
	exec := NewExecutor(nil, sender, resolver)
	session := mcpgw.ToolSessionContext{BotID: "bot1", ReplyTarget: "t1"}
	result, err := exec.CallTool(context.Background(), session, toolSendMessage, map[string]any{
		"platform": "feishu", "text": "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if isErr, _ := result["isError"].(bool); !isErr {
		t.Error("expected error when Send fails")
	}
}

func TestExecutor_CallTool_Success(t *testing.T) {
	sender := &fakeSender{}
	resolver := &fakeResolver{ct: channel.ChannelType("feishu")}
	exec := NewExecutor(nil, sender, resolver)
	session := mcpgw.ToolSessionContext{BotID: "bot1", CurrentPlatform: "feishu", ReplyTarget: "chat1"}
	result, err := exec.CallTool(context.Background(), session, toolSendMessage, map[string]any{
		"text": "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := mcpgw.PayloadError(result); err != nil {
		t.Fatal(err)
	}
	content, _ := result["structuredContent"].(map[string]any)
	if content == nil {
		t.Fatal("no structuredContent")
	}
	if content["ok"] != true {
		t.Errorf("ok = %v", content["ok"])
	}
	if content["platform"] != "feishu" {
		t.Errorf("platform = %v", content["platform"])
	}
}

func TestParseOutboundMessage(t *testing.T) {
	tests := []struct {
		name      string
		args      map[string]any
		fallback  string
		wantEmpty bool
		wantErr   bool
	}{
		{"text fallback", map[string]any{}, "hello", false, false},
		{"message string", map[string]any{"message": "msg"}, "", false, false},
		{"message object", map[string]any{"message": map[string]any{"text": "obj"}}, "", false, false},
		{"empty", map[string]any{}, "", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := parseOutboundMessage(tt.args, tt.fallback)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOutboundMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantEmpty && !msg.IsEmpty() {
				t.Error("expected empty message")
			}
			if !tt.wantEmpty && msg.IsEmpty() {
				t.Error("expected non-empty message")
			}
		})
	}
}
