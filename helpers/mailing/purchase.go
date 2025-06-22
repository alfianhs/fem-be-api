package mailing_helpers

import (
	mongo_model "app/domain/model/mongo"
	"app/helpers"
	"fmt"

	"github.com/sirupsen/logrus"
)

func SendTicketPurchase(ticketPurchases []*mongo_model.TicketPurchase) {
	ticketQrByEmail := make(map[string]map[string][]byte)

	for count, ticketPurchase := range ticketPurchases {
		// generate qr png
		qrCodePng, err := helpers.GenerateQRCodePNG(ticketPurchase.Code)
		if err != nil {
			logrus.Error("Failed to generate QR code:", err)
			continue
		}

		// format date
		date := ticketPurchase.Ticket.Date.Format("02-01-2006")

		// filename
		filename := fmt.Sprintf("QR_%s_%s_%d.png", helpers.SanitizeString(ticketPurchase.Ticket.Name), helpers.SanitizeString(date), count)

		// make map ticket qr by email
		if _, ok := ticketQrByEmail[ticketPurchase.Member.Email]; !ok {
			ticketQrByEmail[ticketPurchase.Member.Email] = make(map[string][]byte)
		}

		ticketQrByEmail[ticketPurchase.Member.Email][filename] = qrCodePng
	}

	// get email template
	subject, body := helpers.GetEmailTicketPurchaseQRTemplate()

	// get fe url
	baseFeUrl := helpers.GetFEUrl()
	ticketPurchaseUrl := fmt.Sprintf("%s/member/ticket-purchases", baseFeUrl)

	for email, qrFileNames := range ticketQrByEmail {

		// replace string template
		dataReplace := map[string]string{
			"ticket_purchase_url": ticketPurchaseUrl,
		}

		finalBody := helpers.StringReplacer(body, dataReplace)

		// setup mail content
		mailer := helpers.NewSMTPMailer()
		mailer.To([]string{email})
		mailer.Subject(subject)
		mailer.Body(finalBody)

		for filename, fileContent := range qrFileNames {
			mailer.Attachment(fileContent, filename, "image/png")
		}

		// send
		if err := mailer.Send(); err != nil {
			logrus.Errorf("Send Email to %s error %v", email, err)
		}
	}
}
