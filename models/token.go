package models

import "time"

type FCMToken struct {
	UserID    string    `bson:"userId" json:"userId"`
	Token     string    `bson:"token" json:"token"`
	Platform  string    `bson:"platform,omitempty" json:"platform,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
