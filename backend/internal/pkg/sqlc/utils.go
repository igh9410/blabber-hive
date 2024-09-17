package sqlc

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUIDToPgtype converts uuid.UUID to pgtype.UUID
func UUIDToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// PgtypeToUUID converts pgtype.UUID to uuid.UUID
func PgtypeToUUID(id pgtype.UUID) uuid.UUID {
	if !id.Valid {
		return uuid.Nil
	}
	return id.Bytes
}

// StringToPgtype converts string to pgtype.Text
func StringToPgtype(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

// PgtypeToString converts pgtype.Text to string
func PgtypeToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// TimeToPgtype converts time.Time to pgtype.Timestamp
func TimeToPgtype(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: !t.IsZero()}
}

// PgtypeToTime converts pgtype.Timestamp to time.Time
func PgtypeToTime(t pgtype.Timestamp) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// Int32ToPgtype converts int32 to pgtype.Int4
func Int32ToPgtype(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

// PgtypeToInt32 converts pgtype.Int4 to int32
func PgtypeToInt32(i pgtype.Int4) int32 {
	if !i.Valid {
		return 0
	}
	return i.Int32
}

// Int64ToPgtype converts int64 to pgtype.Int8
func Int64ToPgtype(i int64) pgtype.Int8 {
	return pgtype.Int8{Int64: i, Valid: true}
}

// PgtypeToInt64 converts pgtype.Int8 to int64
func PgtypeToInt64(i pgtype.Int8) int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

// Float64ToPgtype converts float64 to pgtype.Float8
func Float64ToPgtype(f float64) pgtype.Float8 {
	return pgtype.Float8{Float64: f, Valid: true}
}

// PgtypeToFloat64 converts pgtype.Float8 to float64
func PgtypeToFloat64(f pgtype.Float8) float64 {
	if !f.Valid {
		return 0
	}
	return f.Float64
}

// BoolToPgtype converts bool to pgtype.Bool
func BoolToPgtype(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

// PgtypeToBool converts pgtype.Bool to bool
func PgtypeToBool(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}
