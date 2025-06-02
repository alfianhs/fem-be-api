package request

type CreatePurchaseRequest struct {
	ProductId string `json:"productId"`
	Amount    int64  `json:"amount"`
}
