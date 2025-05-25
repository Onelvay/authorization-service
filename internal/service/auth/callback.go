package auth

import (
	"account-service/internal/domain/billing"
	"account-service/internal/provider/epay"
	"account-service/pkg/log"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"strings"
)

func (s *Service) Callback(ctx context.Context, id string, body []byte) (err error) {
	logger := log.LoggerFromContext(ctx).Named("Callback")

	data, err := s.userRepository.GetBillingByID(ctx, id)
	if err != nil {
		logger.Error("failed to get billing by id", zap.Error(err))
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
			if strings.Contains(billingData.Description, "Premium") ||
				strings.Contains(billingData.Description, "Standard") ||
				strings.Contains(billingData.Description, "Basic") {
				_, err = s.userRepository.CreateSub(ctx, billingData.AccountID, billingData.Description)
				if err != nil {
					logger.Error("failed to save ", zap.Error(err))
					return
				}
			}

			var cards []billing.CardEntity
			cards, err = s.userRepository.GetCards(ctx, callback.AccountID)
			if err != nil {
				logger.Error("failed to get cards", zap.Error(err), zap.Any("account id", callback.AccountID))
				return
			}

			found := false
			for _, card := range cards {
				if card.Mask == callback.CardMask {
					found = true
				}
			}

			if !found {
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
	}
	return
}
