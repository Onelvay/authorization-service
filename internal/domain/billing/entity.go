package billing

import (
	"account-service/internal/provider/epay"
	"time"
)

type Entity struct {
	ID              string  `db:"id" json:"id"`
	CorrelationID   string  `db:"correlation_id" json:"correlationId"`
	InvoiceID       *string `db:"invoice_id" json:"invoiceId"`
	IIN             string  `db:"iin" json:"iin"`
	Phone           string  `db:"phone" json:"phone"`
	Source          string  `db:"source" json:"source"`
	Amount          string  `db:"amount" json:"amount"`
	Currency        string  `db:"currency" json:"currency"`
	TerminalID      string  `db:"terminal_id" json:"terminalId"`
	Description     string  `db:"description" json:"description"`
	AccountID       string  `db:"account_id" json:"accountId"`
	Name            string  `db:"name" json:"name"`
	Email           string  `db:"email" json:"email"`
	Data            string  `db:"data" json:"data"`
	BackLink        string  `db:"back_link" json:"backLink"`
	FailureBackLink string  `db:"failure_back_link" json:"failureBackLink"`
	PostLink        string  `db:"post_link" json:"post_link"`
	FailurePostLink string  `db:"failure_post_link" json:"failurePostLink"`
	Language        string  `db:"language" json:"language"`
	CardSave        bool    `db:"card_save" json:"cardSave"`
}

type CardEntity struct {
	ID         string    `json:"id" db:"id"`
	CardID     string    `json:"card_id" db:"card_id"`
	AccountID  string    `json:"account_id" db:"account_id"`
	TerminalID string    `json:"terminal_id" db:"terminal_id"`
	Type       string    `json:"type" db:"type"`
	Mask       string    `json:"mask" db:"mask"`
	Issuer     string    `json:"issuer" db:"issuer"`
	IsDefault  bool      `json:"is_default" db:"is_default"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func ParseToEpayRequest(b Entity) epay.Request {
	return epay.Request{
		ID:            b.ID,
		IIN:           b.IIN,
		CorrelationID: b.CorrelationID,
		InvoiceID:     *b.InvoiceID,
		Amount:        b.Amount,
		Currency:      b.Currency,
		TerminalID:    b.TerminalID,
		Description:   b.Description,
		AccountID:     b.AccountID,
		Name:          b.Name,
		Email:         b.Email,
		Phone:         b.Phone,
		Language:      b.Language,
		PostLink:      b.PostLink + b.ID,
		BackLink:      b.BackLink,
		CardSave:      true,
	}
}
