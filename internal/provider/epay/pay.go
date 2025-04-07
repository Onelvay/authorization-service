package epay

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type Request struct {
	CreatedAt       *time.Time     `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt       *time.Time     `json:"updatedAt,omitempty" db:"updated_at"`
	ID              string         `json:"id,omitempty" db:"id"`
	IIN             string         `json:"iin" db:"iin" validate:"required,len=12"`
	CorrelationID   string         `json:"correlationId" db:"correlation_id" validate:"required,uuid4"`
	InvoiceID       string         `json:"invoiceId" db:"invoice_id" validate:"required"`
	Amount          string         `json:"amount" db:"amount" validate:"required"`
	Currency        string         `json:"currency" db:"currency" validate:"required"`
	TerminalID      string         `json:"terminalId" db:"terminal_id" validate:"required"`
	Description     string         `json:"description" db:"description" validate:"required"`
	AccountID       string         `json:"accountId" db:"account_id"`
	Name            string         `json:"name" db:"name"`
	Email           string         `json:"email" db:"email"`
	Phone           string         `json:"phone" db:"phone" validate:"required"`
	BackLink        string         `json:"backLink" db:"back_link" validate:"required"`
	FailureBackLink string         `json:"failureBackLink" db:"failure_back_link"`
	PostLink        string         `json:"postLink" db:"post_link" validate:"required"`
	FailurePostLink string         `json:"failurePostLink" db:"failure_post_link"`
	Language        string         `json:"language" db:"language"`
	Data            string         `json:"data" db:"data"`
	CardSave        bool           `json:"cardSave" db:"card_save"`
	PaymentType     string         `json:"paymentType" db:"payment_type"`
	CardID          interface{}    `json:"cardId,omitempty" db:"card_id"`
	DueDate         *time.Time     `json:"-"`
	PaymentJsLink   string         `json:"-"`
	HomebankToken   string         `json:"-"`
	Token           Token          `json:"-"`
	Status          StatusResponse `json:"-"`
	ReceiptLink     string         `json:"-" db:"-"`

	RetryLink string
	Retry     bool
}

type RequestByCardID struct {
	ID          string      `json:"id,omitempty" db:"id"`
	IIN         string      `json:"iin" db:"iin" validate:"required,len=12"`
	InvoiceID   string      `json:"invoiceId" db:"invoice_id" validate:"required"`
	Amount      float64     `json:"amount" db:"amount" validate:"required"`
	Currency    string      `json:"currency" db:"currency" validate:"required"`
	TerminalID  string      `json:"terminalId" db:"terminal_id" validate:"required"`
	Description string      `json:"description" db:"description" validate:"required"`
	PaymentType string      `json:"paymentType"`
	AccountID   string      `json:"accountId" db:"account_id"`
	Name        string      `json:"name" db:"name"`
	Email       string      `json:"email" db:"email"`
	Phone       string      `json:"phone" db:"phone" validate:"required"`
	CardID      interface{} `json:"cardId,omitempty" db:"card_id"`
	PostLink    string      `json:"postLink"`
}

func parseToRequestByCardID(data Request) (res RequestByCardID, err error) {
	amount, err := decimal.NewFromString(data.Amount)
	if err != nil {
		return
	}

	return RequestByCardID{
		ID:          data.ID,
		IIN:         data.IIN,
		InvoiceID:   data.InvoiceID,
		Amount:      amount.InexactFloat64(),
		Currency:    data.Currency,
		TerminalID:  data.TerminalID,
		Description: data.Description,
		AccountID:   data.AccountID,
		Name:        data.Name,
		CardID:      data.CardID,
		PostLink:    data.PostLink,
		PaymentType: "cardId",
	}, nil
}

type PaymentCardID struct {
	ID interface{} `json:"id"`
}

type Response struct {
	ID           string      `json:"id,omitempty"`
	AccountID    string      `json:"accountId,omitempty"`
	Amount       float32     `json:"amount,omitempty"`
	AmountBonus  int         `json:"amountBonus,omitempty"`
	Currency     string      `json:"currency,omitempty"`
	Description  string      `json:"description,omitempty"`
	Email        string      `json:"email,omitempty"`
	InvoiceID    string      `json:"invoiceId,omitempty"`
	CardIssuer   string      `json:"issuer"`
	Language     string      `json:"language,omitempty"`
	Phone        string      `json:"phone,omitempty"`
	Reference    string      `json:"reference,omitempty"`
	IntReference string      `json:"intReference,omitempty"`
	Secure3D     interface{} `json:"secure3D,omitempty"`
	CardID       string      `json:"cardId,omitempty"`
	CardMask     string      `json:"cardMask"`
	CardType     string      `json:"cardType"`
	TerminalID   string      `json:"terminal"`
	PaymentLink  string      `json:"paymentLink,omitempty"`
	Status       string      `json:"status"`
	Code         string      `json:"code"`
	ApprovalCode string      `json:"approvalCode"`
}

type ResponseCardID struct {
	ID           string      `json:"id,omitempty"`
	AccountID    string      `json:"accountId,omitempty"`
	Amount       int         `json:"amount,omitempty"`
	AmountBonus  int         `json:"amountBonus,omitempty"`
	Currency     string      `json:"currency,omitempty"`
	Description  string      `json:"description,omitempty"`
	Email        string      `json:"email,omitempty"`
	InvoiceID    string      `json:"invoiceId,omitempty"`
	Issuer       string      `json:"issuer"`
	Language     string      `json:"language,omitempty"`
	Phone        string      `json:"phone,omitempty"`
	Reference    string      `json:"reference,omitempty"`
	IntReference string      `json:"intReference,omitempty"`
	Secure3D     interface{} `json:"secure3D,omitempty"`
	CardID       string      `json:"cardId,omitempty"`
	CardMask     string      `json:"cardMask"`
	Terminal     string      `json:"terminal"`
	PaymentLink  string      `json:"paymentLink,omitempty"`
	Status       string      `json:"status"`
	Code         int         `json:"code"`
	ApprovalCode string      `json:"approvalCode"`
	Message      string      `json:"message"`
}

func (c *Client) PayByTemplate(w http.ResponseWriter, requestSrc Request) (err error) {

	filenames := ""
	switch requestSrc.Status.Transaction.StatusName {
	case "NEW", "AUTH", "EXPIRED":
		requestSrc.Status.Transaction.Status = "pending"
		filenames = "success.html"
	case "CHARGE":
		filenames = "success.html"
		requestSrc.Status.Transaction.Status = "success"
	case "CANCEL", "CANCEL_OLD", "REFUND":
		filenames = "success.html"
		requestSrc.Status.Transaction.Status = "cancel"
	case "REJECT", "FAILED", "3D":
		filenames = "success.html"
		requestSrc.Status.Transaction.Status = "failed"
	case "":
		filenames = "payment.html"
		token, err := c.getToken(requestSrc)
		if err != nil {
			return err
		}

		requestSrc.Token = token
		requestSrc.PaymentJsLink = c.credential.JS

		requestSrc.BackLink = c.configs.APP.Host + "/invoices/" + requestSrc.ID + "/pay"

	default:
		requestSrc.Status.Transaction.Status = "failed"
		filenames = "status.html"
	}

	requestSrc.Status.Transaction.Datetime = requestSrc.Status.Transaction.CreatedDate.Format(time.DateTime)

	// setup request handler
	tmpl, err := template.ParseFiles(filenames)
	if err != nil {
		return
	}

	return tmpl.Execute(w, requestSrc)
}

func (c *Client) PayByCard(requestSrc Request) (responseSrc *ResponseCardID, err error) {
	//preparation of request params
	token, err := c.getToken(requestSrc)
	if err != nil {
		return
	}

	req, err := parseToRequestByCardID(requestSrc)
	if err != nil {
		return
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return
	}
	reqBytes := bytes.NewReader(reqBody)

	// setup request handler
	path := c.credential.Endpoint + "/payments/cards/auth"
	resBytes, status, err := c.handler("POST", path, "", token.AccessToken, reqBytes)
	if err != nil {
		return
	}

	// check response status
	switch status {
	case http.StatusOK:
		err = json.Unmarshal(resBytes, &responseSrc)
	default:
		err = errors.New(string(resBytes))
	}

	return
}
