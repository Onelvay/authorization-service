package auth

import (
	"account-service/internal/domain/billing"
	"account-service/pkg/log"
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"net/http"
)

func (s *Service) CreatePayment(ctx context.Context, req billing.Entity) (dest string, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreatePayment")

	invoiceID, err := generateInvoiceID(15)
	if err != nil {
		logger.Error("failed to generate invoice id", zap.Error(err))
		return
	}
	req.InvoiceID = &invoiceID
	req.PostLink = "https://authorization-service-4b7m.onrender.com/auth/callback" + "?id="
	req.Source = "epay"
	req.Currency = "KZT"
	req.TerminalID = "67e34d63-102f-4bd1-898e-370781d0074d"
	req.Language = "ru"
	req.CardSave = true
	req.BackLink = "http://localhost:5173/my-appointments"

	if dest, err = s.userRepository.CreateBilling(ctx, req); err != nil {
		logger.Error("failed to create billing", zap.Error(err))
		return
	}
	return
}

func (s *Service) Pay(ctx context.Context, w http.ResponseWriter, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("Pay")

	billingData, err := s.userRepository.GetBillingByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("Платеж не существует")
		}
		logger.Error("failed to get", zap.Error(err))
		return
	}
	req := billing.ParseToEpayRequest(billingData)

	paymentRes, err := s.epayClient.CheckStatus(*billingData.InvoiceID, billingData.TerminalID)
	if err != nil {
		logger.Error("failed to check invoice status", zap.Error(err))
		return
	}
	req.Status = paymentRes

	if err = s.epayClient.PayByTemplate(w, req); err != nil {
		logger.Error("failed to pay by template", zap.Error(err))
		return
	}
	return
}
