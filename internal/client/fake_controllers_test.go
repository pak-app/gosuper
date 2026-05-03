package client

import (
	"fmt"
	"net/http"
)

const SimpleResponse string = `{"message":"%s"}`
const SupervisorStatusResponse string = `
{
	"%s": {
		"name": "%s",
		"services": {}
	}
}
`

func successMessageResponse(message string, statusCode int) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(SimpleResponse, message)))
	}
}

func successServiceStatusResponse_1Sup(name string) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(SupervisorStatusResponse, name, name)))
	}
}
