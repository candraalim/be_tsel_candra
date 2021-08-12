package referral

type ReferResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ReferRequest struct {
	Code   string `json:"code" validate:"required"`
	Msisdn string `json:"msisdn" validate:"required"`
}
