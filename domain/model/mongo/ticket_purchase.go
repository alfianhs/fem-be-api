package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketPurchase struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Member     MemberPurchaseFK   `bson:"member" json:"member"`
	Ticket     TicketFK           `bson:"ticket" json:"ticket"`
	Venue      VenueFK            `bson:"venue" json:"venue"`
	PurchaseID string             `bson:"purchaseId" json:"purchaseId"`
	Code       string             `bson:"code" json:"code"`
	IsUsed     bool               `bson:"isUsed" json:"isUsed"`
	UsedAt     *time.Time         `bson:"usedAt" json:"usedAt"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time         `bson:"deletedAt" json:"deletedAt"`
}
