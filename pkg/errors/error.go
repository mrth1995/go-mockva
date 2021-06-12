package errors

type EndpointError struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
}

func (e *EndpointError) Error() string {
	return e.ErrorMessage
}

func NewAccountAlreadyExist(accountId string) error {
	return &EndpointError{
		ErrorMessage: "Account with ID" + accountId + " already exist",
		ErrorCode:    "68",
	}
}

func NewAccountNotFound(accountId string) error {
	return &EndpointError{
		ErrorMessage: "Account with ID" + accountId + " not found",
		ErrorCode:    "76",
	}
}
