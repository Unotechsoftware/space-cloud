package integration

import (
	"net/http"
)

type adminManager interface {
	GetInternalAccessToken() (string, error)
	// ValidateIntegrationSyncOperation(integrations config.Integrations) error
}

type authResponse struct {
	checkResponse bool
	err           error
	result        interface{}
	status        int

	integration, hook string
}

// CheckResponse indicates whether the integration is hijacking the authentication of the request or not.
// Its a humble way of saying that I'm the boss for this request
func (r authResponse) CheckResponse() bool {
	return r.checkResponse
}

// Error returns error generated by the module if CheckResponse() returns true.
func (r authResponse) Error() error {
	return r.err
}

// Status returns the status code of the hook
func (r authResponse) Status() int {
	if r.status == 0 {
		return http.StatusServiceUnavailable
	}

	return r.status
}

// Result returns the value received from the integration
func (r authResponse) Result() interface{} {
	return r.result
}
