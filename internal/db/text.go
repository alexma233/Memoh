package db

import "github.com/jackc/pgx/v5/pgtype"

// TextToString returns the string value of pgtype.Text, or "" when invalid.
func TextToString(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
