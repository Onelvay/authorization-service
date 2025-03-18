package epay

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type StatusResponse struct {
	InvoiceID     string `json:"invoiceID"`
	ResultCode    string `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
	Transaction   struct {
		ID                string          `json:"id"`
		Datetime          string          `json:"-"`
		CreatedDate       time.Time       `json:"createdDate"`
		InvoiceID         string          `json:"invoiceID"`
		Amount            decimal.Decimal `json:"amount"`
		AmountBonus       decimal.Decimal `json:"amountBonus"`
		PayoutAmount      decimal.Decimal `json:"payoutAmount"`
		Currency          string          `json:"currency"`
		Terminal          string          `json:"terminal"`
		AccountID         string          `json:"accountID"`
		Description       string          `json:"description"`
		Language          string          `json:"language"`
		CardMask          string          `json:"cardMask"`
		CardType          string          `json:"cardType"`
		Issuer            string          `json:"issuer"`
		Reference         string          `json:"reference"`
		IntReference      string          `json:"intReference"`
		Secure            bool            `json:"secure"`
		Status            string          `json:"status"`
		StatusID          string          `json:"statusID"`
		StatusName        string          `json:"statusName"`
		StatusTitle       string          `json:"statusTitle"`
		StatusDescription string          `json:"statusDescription"`
		ReasonCode        string          `json:"reasonCode"`
		Reason            string          `json:"reason"`
		Name              string          `json:"name"`
		Email             string          `json:"email"`
		Phone             string          `json:"phone"`
		CardID            string          `json:"cardID"`
	} `json:"transaction"`
}

func (c *Client) CheckStatus(invoiceID, terminalID string) (res StatusResponse, err error) {
	// preparation of request params
	token, err := c.getToken(Request{TerminalID: terminalID})
	if err != nil {
		return
	}

	// setup request handler
	path := c.credential.Endpoint + "/check-status/payment/transaction/" + invoiceID
	resBytes, status, err := c.handler("GET", path, "", token.AccessToken, nil)
	if err != nil {
		return
	}

	// check response status
	switch status {
	case http.StatusOK, http.StatusBadRequest:
		if err = json.Unmarshal(resBytes, &res); err != nil {
			return
		}

		switch res.Transaction.StatusName {
		case "REFUND":
			res.Transaction.StatusTitle = "Средства возвращены"
		case "AUTH":
			res.Transaction.StatusTitle = "Средства заблокированы"
		case "CANCEL", "CANCEL_OLD":
			res.Transaction.StatusTitle = "Средства разблокированы"
		case "CHARGE":
			res.Transaction.StatusTitle = "Оплата прошла успешно"
		case "FAILED":
			res.Transaction.StatusTitle = "Транзакция неудачна"
		case "3D":
			res.Transaction.StatusTitle = "Ошибка 3D-проверки"
		case "NEW":
			res.Transaction.StatusTitle = "Транзакция в обработке"
		case "REJECT":
			res.Transaction.StatusTitle = "Платеж отклонен"
		default:
			res.Transaction.StatusTitle = "Платеж отклонен"
		}
		res.Transaction.StatusDescription = GetStatusDescription(res.Transaction.ReasonCode)

	default:
		err = errors.New(string(resBytes))
	}

	return
}

func GetStatusDescription(reasonCode string) (description string) {
	switch reasonCode {
	case "454":
		description = "Операция не удалась, проверьте не заблокирована ли сумма на карте и повторите попытку позже"
	case "455":
		description = "Проверка 3DSecure/SecureCode недоступна, либо неверно введен номер карты. Попробуйте воспользоваться другим браузером/устройством. Если ошибка повторяется, переустановите код"
	case "456":
		description = "Невозможно провести оплату по данной карте"
	case "457":
		description = "Некорректно введен срок действия карты"
	case "458":
		description = "Сервер не отвечает. Попробуйте попозже"
	case "459":
		description = "Сервер не отвечает. Попробуйте попозже"
	case "460":
		description = "Произошла ошибка, возможно сумма заблокировалась на карте, обратитесь службу поддержки"
	case "461":
		description = "Системная ошибка, попробуйте провести транзакцию позже, если ошибка повторяется обратитесь в службу поддержки"
	case "462":
		description = "Транзакция отклонена вашим банком. Для уточнения причины отказа необходимо обратиться по контактам, указанным на обратной стороне вашей карты"
	case "463":
		description = "Транзакция отклонена вашим банком. Для уточнения причины отказа необходимо обратиться по контактам, указанным на обратной стороне вашей карты"
	case "464":
		description = "Недействительный коммерсант"
	case "465":
		description = "Карта заблокирована"
	case "466":
		description = "Транзакция отклонена вашим банком. Для уточнения причины отказа необходимо обратиться по контактам, указанным на обратной стороне вашей карты"
	case "467":
		description = "Карта заблокирована"
	case "468":
		description = "Требуется дополнительная идентификация"
	case "469":
		description = "Недействительная транзакция, перепроверить введенные данные"
	case "470":
		description = "Сумма транзакции равно нулю, пожалуйста, попробуйте ещё раз"
	case "471":
		description = "Недействительный номер карточки, пожалуйста, убедитесь в корректности ввода номера карты и попробуйте ещё раз"
	case "472":
		description = "Недействительный номер карточки, пожалуйста, убедитесь в корректности ввода номера карты и попробуйте ещё раз"
	case "473":
		description = "3DSecure/SecureCode введен некорректно. Пожалуйста, убедитесь в корректности ввода, либо переустановите пароль. Если ошибка повторяется обратитесь в службу поддержки"
	case "475":
		description = "Транзакция не успешна. Пожалуйста, повторите снова"
	case "476":
		description = "Повторное проведение транзакции будет доступно не менее чем через 30 минут"
	case "477":
		description = "Ошибка, пожалуйста, воспользуйтесь другой картой. В случае её отсутствия обратитесь в службу поддержки по адресу epay@halykbank.kz"
	case "478":
		description = "Просрочен срок действия карты"
	case "479":
		description = "Карточка заблокирована"
	case "480":
		description = "Обратиться к банку - эмитенту"
	case "481":
		description = "Карта недействительна. Пожалуйста, обратитесь в Банк для выпуска новой карты"
	case "482":
		description = "Карта недействительна. Пожалуйста, обратитесь в Банк для выпуска новой карты"
	case "483":
		description = "Статус карты - украдена. Пожалуйста, обратитесь в Банк для выпуска новой карты"
	case "484":
		description = "Недостаточно средств на карте"
	case "485":
		description = "Срок действия карты истек"
	case "486":
		description = "Транзакция отклонена. На карте запрещена возможность покупок в сети интернет, либо карточные данные введены не верно"
	case "487":
		description = "Транзакция отклонена, пожалуйста, обратитесь в службу поддержки"
	case "488":
		description = "Сумма превышает допустимый лимит"
	case "489":
		description = "Карточка заблокирована"
	case "490":
		description = "Запрет на проведение транзакции по вашей карте, за дополнительной информацией обратитесь по контактам, указанным на обратной стороне вашей карты"
	case "491":
		description = "Превышен лимит частоты оплат"
	case "492":
		description = "Карта заблокирована по причине неверного ввода пин-кода, за дополнительной информацией обратитесь по контактам, указанным на обратной стороне вашей карты"
	case "493":
		description = "Недоступен банк, выпустивший карту, попробуйте повторить оплату позже"
	case "494":
		description = "Недоступен банк, выпустивший карту, попробуйте провести транзакцию позже"
	case "495":
		description = "Транзакция запрещена, воспользуйтесь другой картой"
	case "496":
		description = "Системная ошибка при оплатах"
	case "497":
		description = "Сервер не отвечает. Попробуйте попозже"
	case "498":
		description = "Оплата бонусами невозможна. Попробуйте попозже"
	case "499":
		description = "Неверно введен или не введен 3DSecure/SecureCode"
	case "500":
		description = "Сервер не отвечает. Попробуйте попозже"
	case "501":
		description = "Ошибка обслуживания карты. Проверьте правильность ввода карты. Если ошибка повторяется, обратитесь в службу поддержки"
	case "502":
		description = "Сервер не отвечает. Попробуйте попозже"
	case "503":
		description = "Сервер не отвечает. Попробуйте попозже"
	case "521":
		description = "Транзакция отклонена вашим банком. Для уточнения причины отказа необходимо обратиться по контактам, указанным на обратной стороне вашей карты"
	case "522":
		description = "Запись не найдена, проверьте карточку"
	case "523":
		description = "Транзакция отклонена вашим банком. Для уточнения причины отказа необходимо обратиться по контактам, указанным на обратной стороне вашей карты"
	case "524":
		description = "Карта недействительна. Пожалуйста, обратитесь в Банк"
	case "525":
		description = "Карта недействительна. Пожалуйста, обратитесь в Банк"
	case "526":
		description = "Системная ошибка. Пожалуйста, попробуйте позднее. Если эта проблема не исчезнет, обратитесь в службу поддержки"
	case "527":
		description = "Транзакция отклонена вашим банком. Для уточнения причины отказа необходимо обратиться по контактам, указанным на обратной стороне вашей карты"
	case "528":
		description = "Превышен суточный лимит входящих переводов на карту получателя, предоставьте другую карту для перевода"
	case "529":
		description = "Сработал суточный лимит на терминал"
	default:
		description = "Ошибка с платежного шлюза не найдена"
	}

	return description + " (" + reasonCode + ")"
}
