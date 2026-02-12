package db

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestTextToString(t *testing.T) {
	tests := []struct {
		name  string
		value pgtype.Text
		want  string
	}{
		{"valid", pgtype.Text{String: "hello", Valid: true}, "hello"},
		{"invalid", pgtype.Text{}, ""},
		{"valid empty", pgtype.Text{String: "", Valid: true}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TextToString(tt.value); got != tt.want {
				t.Errorf("TextToString() = %q, want %q", got, tt.want)
			}
		})
	}
}
