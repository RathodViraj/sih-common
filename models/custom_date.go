package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type CustomDate struct {
	time.Time
}

// --- JSON ---
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		cd.Time = time.Time{}
		return nil
	}
	layout := "2006-01-02"
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if cd.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + cd.Time.Format("2006-01-02") + `"`), nil
}

// --- BSON ---
func (cd CustomDate) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if cd.Time.IsZero() {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(cd.Time)
}

func (cd *CustomDate) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var tm time.Time
	if err := bson.UnmarshalValue(t, data, &tm); err != nil {
		return err
	}
	cd.Time = tm
	return nil
}
