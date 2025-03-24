package auth

import (
	"account-service/internal/provider/epay"
	"account-service/pkg/log"
	"context"
	"go.uber.org/zap"
	"net/http"
)

func (s *Service) Pay(ctx context.Context, w http.ResponseWriter) (err error) {
	logger := log.LoggerFromContext(ctx).Named("Pay")

	invoiceId, err := generateInvoiceID(15)
	if err != nil {
		logger.Error("failed to generate invoice id", zap.Error(err))
		return
	}
	postLink := "https://authorization-service-4b7m.onrender.com/auth" + "/callback"

	req := epay.Request{
		ID:            "some id",
		IIN:           "iiiinnnn",
		CorrelationID: "asdfasdfasdfasdfsa",
		InvoiceID:     invoiceId,
		Amount:        "5002",
		Currency:      "KZT",
		TerminalID:    "67e34d63-102f-4bd1-898e-370781d0074d",
		Description:   "asdfaserqwerqwerasdf",
		AccountID:     "qwerqwerqwerqwerqwerw",
		Name:          "Testing Test",
		Email:         "testest@gmail.com",
		Phone:         "+77475684141",
		Language:      "ru",
		PostLink:      postLink,
	}

	if err = s.epayClient.PayByTemplate(w, req); err != nil {
		logger.Error("failed to pay by template", zap.Error(err))
		return
	}
	return
}
