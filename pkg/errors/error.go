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

func NewInvalidAmount() error {
	return &EndpointError{
		ErrorMessage: "Invalid amount",
		ErrorCode:    "68",
	}
}

func NewInvalidAccount() error {
	return &EndpointError{
		ErrorMessage: "Cannot transfer with same account",
		ErrorCode:    "68",
	}
}

func NewInsufficientAmount() error {
	return &EndpointError{
		ErrorMessage: "Insufficient amount",
		ErrorCode:    "51",
	}
}
