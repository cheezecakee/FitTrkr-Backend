package helper

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github/cheezecakee/fitrkr/pkg/logger"
)

// Errors
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	logger.Log.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ClientError(w http.ResponseWriter, status int) {
	logger.Log.InfoLog.Printf("Client error: %d", status)
	http.Error(w, http.StatusText(status), status)
}

func NotFound(w http.ResponseWriter) {
	ClientError(w, http.StatusNotFound)
}

func Clamp(value, min, max int) int {
	switch {
	case value < min:
		return min
	case value > max:
		return max
	default:
		return value
	}
}
