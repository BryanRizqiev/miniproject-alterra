package lib

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func NewNullString(str string) sql.NullString {

	if len(str) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: str,
		Valid:  true,
	}

}

func NewNullTime(t time.Time) sql.NullTime {

	if t.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  t,
		Valid: true,
	}

}

func InsertDefaultValue(str string, dflt string) string {

	if str != "" {
		return str
	}

	return dflt

}

func NewUuid() string {

	return uuid.NewString()

}
