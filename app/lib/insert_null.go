package lib

import "database/sql"

func NewNullString(str string) sql.NullString {

	if len(str) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: str,
		Valid:  true,
	}

}
