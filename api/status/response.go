package status

import "net/http"

var (
	StatusOk = Status{
		Code:    http.StatusOK,
		Message: "Ok",
	}

	StatusInternal = Status{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}

	StatusBadRequest = Status{
		Code:    http.StatusBadRequest,
		Message: "Bad request",
	}

	StatusNotFound = Status{
		Code:    http.StatusNotFound,
		Message: "Not found",
	}

	StatusUnauthorized = Status{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
	}
	StatusForbidden = Status{
		Code:    http.StatusForbidden,
		Message: "Forbidden",
	}
)
