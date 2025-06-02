package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchase struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	Member            MemberPurchaseFK   `bson:"member" json:"member"`
	SeasonId          string             `bson:"seasonId" json:"seasonId"`
	Season            SeasonFK           `bson:"-" json:"season"`
	SeriesId          string             `bson:"seriesId" json:"seriesId"`
	Series            SeriesFK           `bson:"-" json:"series"`
	TicketIds         []string           `bson:"ticketIds" json:"ticketIds"`
	Tickets           []TicketFK         `bson:"-" json:"tickets"`
	Amount            int64              `bson:"amount" json:"amount"`
	Invoice           InvoiceFK          `bson:"invoice" json:"invoice"`
	Price             float64            `bson:"price" json:"price"`
	GrandTotal        float64            `bson:"grandTotal" json:"grandTotal"`
	IsCheckoutPackage bool               `bson:"isCheckoutPackage" json:"isCheckoutPackage"`
	Status            PurchaseStatus     `bson:"status" json:"-"`
	ExpiresAt         time.Time          `bson:"expiresAt" json:"expiresAt"`
	StatusString      string             `bson:"-" json:"status"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt         *time.Time         `bson:"deletedAt" json:"-"`
}

type MemberPurchaseFK struct {
	ID    string `bson:"id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Phone string `bson:"phone" json:"phone"`
}

type InvoiceFK struct {
	InvoiceID          string `bson:"invoiceId" json:"invoiceId"`
	InvoiceExternalID  string `bson:"invoiceExternalId" json:"invoiceExternalId"`
	InvoiceUrl         string `bson:"invoiceUrl" json:"invoiceUrl"`
	PaymentMethod      string `bson:"PaymentMethod" json:"PaymentMethod"`
	MerchantName       string `bson:"merchantName" json:"merchantName"`
	BankCode           string `bson:"bankCode" json:"bankCode"`
	PaymentChannel     string `bson:"paymentChannel" json:"paymentChannel"`
	PaymentDestination string `bson:"paymentDestination" json:"paymentDestination"`
}

func (p *Purchase) Format() *Purchase {
	p.StatusString = PurchaseStatusMap[p.Status].Name
	return p
}
