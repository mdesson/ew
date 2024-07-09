package echoserver

import (
	"fmt"
	"time"
)

type RequestDetails struct {
	Date        time.Time         `json:"date"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	QueryParams map[string]string `json:"query_params"`
	Headers     map[string]string `json:"headers"`
	Body        *string           `json:"body,omitempty"`
	ErrorMsg    *string           `json:"error,omitempty"`
}

func (d RequestDetails) String() string {
	str := d.Date.String()
	str += fmt.Sprintf("\n%s %s", d.Method, d.Path)

	str += "\nQUERY PARAMS"
	for k, v := range d.QueryParams {
		str += fmt.Sprintf("\n%s=%s", k, v)
	}

	str += "\nHEADERS"
	for k, v := range d.Headers {
		str += fmt.Sprintf("\n%s=%s", k, v)
	}

	if d.Body != nil {
		str += "\nBODY"
		str += fmt.Sprintf("\n%s", *d.Body)
	}

	if d.ErrorMsg != nil {
		str += "\nERROR"
		str += fmt.Sprintf("\n%s", d.ErrorMsg)
	}

	return str
}
