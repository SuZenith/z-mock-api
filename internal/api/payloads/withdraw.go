package payloads

type DataPayload struct {
	VerifyCode string  `json:"verifyCode" validate:"required,len=4,numeric"`
	Amount     float64 `json:"amount" validate:"required,gt=0,lte=10000"`
}

type ApplyPayload struct {
	Data        DataPayload `json:"data" validate:"required"`
	AccessToken string      `json:"accessToken" validate:"required"`
}
