package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type CustomTime struct {
	time.Time
}

// --- JSON ---
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		ct.Time = time.Time{}
		return nil
	}
	layout := "15:04" // hh:mm 24-hour
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + ct.Time.Format("15:04") + `"`), nil
}

// --- BSON ---
func (ct CustomTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if ct.Time.IsZero() {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(ct.Time)
}

func (ct *CustomTime) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var tm time.Time
	if err := bson.UnmarshalValue(t, data, &tm); err != nil {
		return err
	}
	ct.Time = tm
	return nil
}
