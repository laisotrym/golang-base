package safeweb_lib_helper

import (
    "database/sql"
    "encoding/json"
)

type NullString struct {
    sql.NullString
}

type NullFloat64 struct {
    sql.NullFloat64
}

type NullInt64 struct {
    sql.NullInt64
}

var nullString = []byte("null")

// Marshalling
func (n NullString) MarshalJSON() ([]byte, error) {
    if n.Valid {
        return json.Marshal(n.String)
    }
    return nullString, nil
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
    if n.Valid {
        return json.Marshal(n.Int64)
    }
    return nullString, nil
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
    if n.Valid {
        return json.Marshal(n.Float64)
    }
    return nullString, nil
}

// Unmarshalling
func (n *NullString) UnmarshalJSON(b []byte) error {
    var s interface{}
    if err := json.Unmarshal(b, &s); err != nil {
        return err
    }
    return n.Scan(s)
}

func (n *NullInt64) UnmarshalJSON(b []byte) error {
    var s interface{}
    if err := json.Unmarshal(b, &s); err != nil {
        return err
    }
    return n.Scan(s)
}

func (n *NullFloat64) UnmarshalJSON(b []byte) error {
    var s interface{}
    if err := json.Unmarshal(b, &s); err != nil {
        return err
    }
    return n.Scan(s)
}
