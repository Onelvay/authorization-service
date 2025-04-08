package auth

import (
	"account-service/internal/domain/billing"
	"context"
	"fmt"
)

func (s *Service) GetCards(ctx context.Context, accId string) (cards []billing.CardEntity, err error) {
	cards, err = s.userRepository.GetCards(ctx, accId)
	if err != nil {
		fmt.Println(err.Error)
	}
	return
}

func (s *Service) DeleteCard(ctx context.Context, cardId string) (err error) {
	err = s.userRepository.DeleteCardByID(ctx, cardId)
	if err != nil {
		fmt.Println(err.Error)
		return
	}
	return
}
