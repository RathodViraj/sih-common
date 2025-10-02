package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Advertisement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrganizationUID string             `bson:"organization_uid" json:"organization_uid"`
	ProductName     string             `bson:"product_name" json:"product_name"`
	Manufacturer    string             `bson:"manufacturer" json:"manufacturer"`
	Price           float64            `bson:"price" json:"price"`
	Discount        float64            `bson:"discount,omitempty" json:"discount,omitempty"`
	Image1          string             `bson:"image1,omitempty" json:"image1,omitempty"`
	Image2          string             `bson:"image2,omitempty" json:"image2,omitempty"`
	Description     string             `bson:"description" json:"description"`
	Chemicals       []string           `bson:"chemicals,omitempty" json:"chemicals,omitempty"`
	SellerLinks     []string           `bson:"seller_links,omitempty" json:"seller_links,omitempty"`
	SenderName      string             `bson:"sender_name" json:"sender_name"`
	Pincodes        []string           `bson:"pincodes,omitempty" json:"pincodes,omitempty"`
	ValidUntil      CustomDate         `bson:"valid_until" json:"valid_until"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	Status          string             `bson:"status" json:"status"`
	VerifiedBy      string             `bson:"verified_by,omitempty" json:"verified_by,omitempty"`
	VerifiedAt      *time.Time         `bson:"verified_at,omitempty" json:"verified_at,omitempty"`
	Remarks         string             `bson:"verification_remarks,omitempty" json:"verification_remarks,omitempty"`
	IDHex           string             `bson:"-" json:"-"`
}

func (ad *Advertisement) IDFromHex(id string) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		ad.ID = objID
	}
}
