package xendit_repository

import (
	"app/domain"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type xenditRepo struct {
	Client                *http.Client
	baseURL               *url.URL
	generateSnapURL       string
	secret                string
	secretBasicAuth       string
	metadataIssuer        string
	xenditInvoiceDuration int64
}

func NewXenditRepo() domain.XenditRepo {
	client := &http.Client{}
	baseURL, _ := url.Parse(os.Getenv("XENDIT_URL"))
	secret := os.Getenv("XENDIT_SECRET_KEY")
	metadataIssuer := os.Getenv("XENDIT_METADATA_ISSUER")
	xenditInvoiceDuration, _ := strconv.Atoi(os.Getenv("XENDIT_INVOICE_DURATION"))
	if xenditInvoiceDuration == 0 {
		xenditInvoiceDuration = 1800
	}

	return &xenditRepo{
		Client:                client,
		baseURL:               baseURL,
		generateSnapURL:       baseURL.ResolveReference(&url.URL{Path: "/v2/invoices"}).String(),
		secret:                secret,
		secretBasicAuth:       base64.StdEncoding.EncodeToString([]byte(secret + ":")),
		metadataIssuer:        metadataIssuer,
		xenditInvoiceDuration: int64(xenditInvoiceDuration),
	}
}
