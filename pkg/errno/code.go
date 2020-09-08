package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}
	ErrQuery            = &Errno{Code: 10003, Message: "Error occurred while getting request query."}

	ErrValidation   = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase     = &Errno{Code: 20002, Message: "Database error."}
	ErrTokenInvalid = &Errno{Code: 20003, Message: "The token was invalid."}

	// Ele errors
	ErrGetEle = &Errno{Code: 30001, Message: "Errors occurred while getting electricity info."}
)
