package db

import (
	"database/sql"
	"encoding/json"
)

type NullString sql.NullString
type NullInt32 sql.NullInt32
type NullInt64 sql.NullInt64

type Year uint16
type NullYear struct {
	Year Year
	Valid  bool
}

func (ny *NullYear) MarshalJSON() ([]byte, error) {
	if !ny.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ny.Year)
}

func (ny *NullYear) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v == nil {
		ny.Valid = false
	} else {
		ny.Valid = true
		ny.Year = Year(v.(float64))
	}

	return nil
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v == nil {
		ns.Valid = false
	} else {
		ns.Valid = true
		ns.String = v.(string)
	}

	return nil
}