package auth

import (
	"account-service/internal/domain/billing"
	"account-service/internal/provider/epay"
	"account-service/pkg/log"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"strings"
)

func (s *Service) Callback(ctx context.Context, id string, body []byte) (err error) {
	logger := log.LoggerFromContext(ctx).Named("Callback")

	data, err := s.userRepository.GetBillingByID(ctx, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error("failed to get billing by id", zap.Error(err))
		}
		return
	}

	err = s.callbackByEpay(ctx, data, body)
	return
}

func (s *Service) callbackByEpay(ctx context.Context, billingData billing.Entity, body []byte) (err error) {
	logger := log.LoggerFromContext(ctx).Named("callbackByEpay")

	billingData, err = s.userRepository.GetBillingByID(ctx, billingData.ID)
	if err != nil {
		logger.Error("failed to get billing", zap.Error(err))
		return
	}

	var callback epay.Response
	if err = json.Unmarshal(body, &callback); err != nil {
		logger.Error("failed to unmarshall epay callback", zap.Error(err), zap.Any("callback struct", string(body)))
		return
	}

	epayCallback, err := s.epayClient.CheckStatus(callback.InvoiceID, "67e34d63-102f-4bd1-898e-370781d0074d")
	if err != nil {
		logger.Error("failed to check invoice status", zap.Error(err))
		return
	}

	//если был передан параметр cardSave = true, тогда к account id мы привязываем card id
	if strings.EqualFold(epayCallback.Transaction.StatusName, "AUTH") || strings.EqualFold(epayCallback.Transaction.StatusName, "CHARGE") {
		if billingData.CardSave {
			_, err = s.userRepository.CreateCard(ctx, billing.CardEntity{
				CardID:     callback.CardID,
				AccountID:  callback.AccountID,
				TerminalID: callback.TerminalID,
				Type:       callback.CardType,
				Mask:       callback.CardMask,
				Issuer:     callback.CardIssuer,
				IsDefault:  false,
			})
			if err != nil {
				logger.Error("failed to save card", zap.Error(err), zap.Any("billing struct", billingData))
				return
			}
		}
	}
	return
}
