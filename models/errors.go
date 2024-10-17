package models

type FriendlyError struct {
	ErrorCode            string `json:"errorCode"`
	ErrorMessage         string `json:"errorMessage"`
	FriendlyErrorMessage string `json:"friendlyErrorMessage"`
}

func FriendlyErrorInit(errCode string) FriendlyError {
	friendlyError := FriendlyError{
		ErrorCode: errCode,
	}
	switch friendlyError.ErrorCode {
	case "400":
		friendlyError.FriendlyErrorMessage = "Looks like there was a problem with your request"
		friendlyError.ErrorMessage = "Bad Request"
	case "404":
		friendlyError.FriendlyErrorMessage = "The page you're looking for doesn't seem to exist."
		friendlyError.ErrorMessage = "Page Not Found"
	default: // 500
		friendlyError.FriendlyErrorMessage = "We're having some issues on our end. Please try again later"
		friendlyError.ErrorMessage = "Internal Server Error"
	}
	return friendlyError
}
