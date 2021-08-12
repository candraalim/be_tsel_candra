package util

const (
	CodeSuccess    = "0000"
	MessageSuccess = "Success"
)

type ApplicationError struct {
	HttpStatus int    `json:"-"`
	ErrorCode  string `json:"code"`
	Message    string `json:"message"`
}

func (e *ApplicationError) Error() string {
	return e.Message
}

var (
	ErrorDataNotFound         = &ApplicationError{HttpStatus: 400, ErrorCode: "0001", Message: "data not found"}
	ErrorInvalidRequest       = &ApplicationError{HttpStatus: 400, ErrorCode: "0077", Message: "invalid request"}
	ErrorDatabase             = &ApplicationError{HttpStatus: 500, ErrorCode: "0081", Message: "unexpected error"}
	ErrorGenerateReferralCode = &ApplicationError{HttpStatus: 500, ErrorCode: "0082", Message: "unexpected error"}
	ErrorGeneral              = &ApplicationError{HttpStatus: 500, ErrorCode: "9999", Message: "system internal error"}
)
