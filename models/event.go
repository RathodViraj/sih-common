package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrganizationUID string             `bson:"organization_uid" json:"organization_uid"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Location        string             `bson:"location" json:"location"`
	Date            CustomDate         `bson:"date" json:"date"`
	Time            CustomTime         `bson:"time" json:"time"`
	Image1          string             `bson:"image1,omitempty" json:"image1,omitempty"`
	Image2          string             `bson:"image2,omitempty" json:"image2,omitempty"`
	Link            string             `bson:"link,omitempty" json:"link,omitempty"`
	SenderName      string             `bson:"sender_name" json:"sender_name"`
	Pincodes        []string           `bson:"pincodes,omitempty" json:"pincodes,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	Status          string             `bson:"status" json:"status"`
	VerifiedBy      string             `bson:"verified_by,omitempty" json:"verified_by,omitempty"`
	VerifiedAt      *time.Time         `bson:"verified_at,omitempty" json:"verified_at,omitempty"`
	Remarks         string             `bson:"verification_remarks,omitempty" json:"verification_remarks,omitempty"`
	IDHex           string             `bson:"-" json:"-"`
}
