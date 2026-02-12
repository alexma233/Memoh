package telegram

import (
	"context"
	"strings"
	"testing"

	"github.com/memohai/memoh/internal/channel"
)

func TestTelegramOutboundStream_CloseNil(t *testing.T) {
	t.Parallel()

	var s *telegramOutboundStream
	ctx := context.Background()
	if err := s.Close(ctx); err != nil {
		t.Fatalf("Close on nil stream should return nil: %v", err)
	}
}

func TestTelegramOutboundStream_PushClosed(t *testing.T) {
	t.Parallel()

	adapter := NewTelegramAdapter(nil)
	s := &telegramOutboundStream{adapter: adapter}
	s.closed.Store(true)

	ctx := context.Background()
	err := s.Push(ctx, channel.StreamEvent{Type: channel.StreamEventDelta, Delta: "x"})
	if err == nil {
		t.Fatal("Push on closed stream should return error")
	}
	if !strings.Contains(err.Error(), "closed") {
		t.Fatalf("expected closed error: %v", err)
	}
}

func TestTelegramOutboundStream_PushStatusNoOp(t *testing.T) {
	t.Parallel()

	adapter := NewTelegramAdapter(nil)
	s := &telegramOutboundStream{adapter: adapter}

	ctx := context.Background()
	err := s.Push(ctx, channel.StreamEvent{Type: channel.StreamEventStatus})
	if err != nil {
		t.Fatalf("StreamEventStatus should be no-op: %v", err)
	}
}

func TestTelegramOutboundStream_PushNilAdapter(t *testing.T) {
	t.Parallel()

	s := &telegramOutboundStream{adapter: nil}
	ctx := context.Background()
	err := s.Push(ctx, channel.StreamEvent{Type: channel.StreamEventDelta, Delta: "x"})
	if err == nil {
		t.Fatal("Push with nil adapter should return error")
	}
	if !strings.Contains(err.Error(), "not configured") {
		t.Fatalf("expected not configured error: %v", err)
	}
}

func TestTelegramOutboundStream_PushUnsupportedEventType(t *testing.T) {
	t.Parallel()

	adapter := NewTelegramAdapter(nil)
	s := &telegramOutboundStream{adapter: adapter}
	ctx := context.Background()

	err := s.Push(ctx, channel.StreamEvent{Type: channel.StreamEventType("unknown")})
	if err == nil {
		t.Fatal("Push with unknown event type should return error")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("expected unsupported error: %v", err)
	}
}

func TestTelegramOutboundStream_PushEmptyDeltaNoOp(t *testing.T) {
	t.Parallel()

	adapter := NewTelegramAdapter(nil)
	s := &telegramOutboundStream{adapter: adapter}
	ctx := context.Background()

	err := s.Push(ctx, channel.StreamEvent{Type: channel.StreamEventDelta, Delta: ""})
	if err != nil {
		t.Fatalf("empty delta should be no-op: %v", err)
	}
}

func TestTelegramOutboundStream_PushErrorEventEmptyNoOp(t *testing.T) {
	t.Parallel()

	adapter := NewTelegramAdapter(nil)
	s := &telegramOutboundStream{adapter: adapter}
	ctx := context.Background()

	err := s.Push(ctx, channel.StreamEvent{Type: channel.StreamEventError, Error: ""})
	if err != nil {
		t.Fatalf("empty error event should be no-op: %v", err)
	}
}

func TestTelegramOutboundStream_CloseContextCanceled(t *testing.T) {
	t.Parallel()

	adapter := NewTelegramAdapter(nil)
	s := &telegramOutboundStream{adapter: adapter}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := s.Close(ctx)
	if err != context.Canceled {
		t.Fatalf("Close with canceled context should return context.Canceled: %v", err)
	}
}
