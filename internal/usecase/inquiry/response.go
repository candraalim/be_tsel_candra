package inquiry

type ReferralCodeResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ReferralCode string `json:"referralCode"`
	} `json:"data"`
}

type ReferralRewardResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TotalReferral int    `json:"totalReferral"`
		Reward        string `json:"reward"`
	} `json:"data"`
}

type ReferralHistoryResponse struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Data    ReferralHistoryData `json:"data"`
}

type ReferralHistoryData struct {
	List []ReferralHistory `json:"list"`
	Meta Meta              `json:"meta"`
}

type ReferralHistory struct {
	Msisdn       string `json:"msisdn"`
	ReferralDate string `json:"referralDate"`
	DateTime     int64  `json:"dateTime"`
}

type Meta struct {
	TotalPage   int  `json:"totalPage"`
	TotalRecord int  `json:"totalRecord"`
	Page        int  `json:"page"`
	Size        int  `json:"size"`
	Limit       int  `json:"limit"`
	FirstPage   bool `json:"firstPage"`
	LastPage    bool `json:"lastPage"`
}
