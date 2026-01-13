package errors

type EndpointError struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
}

func (e *EndpointError) Error() string {
	return e.ErrorMessage
}

func NewAccountAlreadyExist(accountID string) error {
	return &EndpointError{
		ErrorMessage: "Account with ID" + accountID + " already exist",
		ErrorCode:    "68",
	}
}

func NewAccountNotFound(accountID string) error {
	return &EndpointError{
		ErrorMessage: "Account with ID" + accountID + " not found",
		ErrorCode:    "76",
	}
}
