package xendit_repository

import (
	mongo_model "app/domain/model/mongo"
	xendit_model "app/domain/model/xendit"
	"app/helpers"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func (r *xenditRepo) GenereteSnapLink(ctx context.Context, purchase mongo_model.Purchase) (result helpers.Response, err error) {
	itemName := ""

	if purchase.IsCheckoutPackage {
		itemName = fmt.Sprintf("%s (%s)", purchase.Series.Name, purchase.Season.Name)
	} else {
		itemName = purchase.Tickets[0].Name
	}

	generateSnapUrlRequest := struct {
		ExternalId      string                   `json:"external_id"`
		Amount          int64                    `json:"amount"`
		PayerEmail      string                   `json:"payer_email"`
		Currency        string                   `json:"currency"`
		Locale          string                   `json:"locale"`
		Description     string                   `json:"description"`
		InvoiceDuration int64                    `json:"invoice_duration"`
		Customer        map[string]interface{}   `json:"customer"`
		Items           []map[string]interface{} `json:"items"`
		Metadata        map[string]interface{}   `json:"metadata"`
	}{
		ExternalId:      purchase.Invoice.InvoiceExternalID,
		Amount:          int64(purchase.GrandTotal),
		PayerEmail:      purchase.Member.Email,
		Currency:        "IDR",
		Locale:          "id",
		Description:     r.metadataIssuer,
		InvoiceDuration: r.xenditInvoiceDuration,
		Customer: map[string]interface{}{
			"given_names": purchase.Member.Name,
			"email":       purchase.Member.Email,
			"phone":       purchase.Member.Phone,
		},
		Items: []map[string]interface{}{
			{
				"name":     itemName,
				"quantity": purchase.Amount,
				"price":    purchase.Price,
			},
		},
		Metadata: map[string]interface{}{
			"issuer": r.metadataIssuer,
		},
	}

	// marshal json
	jsonByte, err := json.Marshal(generateSnapUrlRequest)

	// convert to string
	payload := strings.NewReader(string(jsonByte))

	// send request
	req, err := http.NewRequest("POST", r.generateSnapURL, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+r.secretBasicAuth)

	res, err := r.Client.Do(req)
	if err != nil {
		logrus.Error("Generete Snap Link", err)
		return
	}
	defer res.Body.Close()

	// read body
	body, err := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		failed := xendit_model.XenditResponseError{}
		err = json.Unmarshal(body, &failed)
		if err != nil {
			logrus.Error("Generete Snap Link Response Unmarshal", err)
			result.Status = 400
			result.Message = "Generete Snap Link Response Unmarshal"
			return
		}
		result.Status = res.StatusCode
		logrus.Error("Generete Snap Link Response", failed)
		result.Message = failed.Message
		return
	}

	if err != nil {
		logrus.Error("Generete Snap Link Response ReadBody", err)
		result.Status = 400
		result.Message = err.Error()
		return
	}

	// unmarshal to struct
	successData := xendit_model.XenditSnapLinkSuccessResponse{}
	err = json.Unmarshal(body, &successData)
	if err != nil {
		logrus.Error("Generete Snap Link Response Unmarshal", err)
		result.Status = 400
		result.Message = err.Error()
		return
	}

	result.Status = res.StatusCode
	result.Message = "success"
	result.Data = successData

	return
}
