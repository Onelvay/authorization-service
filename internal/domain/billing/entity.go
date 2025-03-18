package billing

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
	BackLink        string  `db:"back_link" json:"backLink"`
	FailureBackLink string  `db:"failure_back_link" json:"failureBackLink"`
	PostLink        string  `db:"post_link" json:"post_link"`
	FailurePostLink string  `db:"failure_post_link" json:"failurePostLink"`
	Language        string  `db:"language" json:"language"`
	CardSave        bool    `db:"card_save" json:"cardSave"`
}
