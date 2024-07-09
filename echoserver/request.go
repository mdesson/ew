package echoserver

import (
	"fmt"
	"time"
)

type RequestDetails struct {
	Date        time.Time
	Method      string
	Path        string
	QueryParams map[string]string
	Headers     map[string]string
	Body        string
	Err         error
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

	if d.Body != "" {
		str += "\nBODY"
		str += fmt.Sprintf("\n%s", d.Body)
	}

	if d.Err != nil {
		str += "\nERROR"
		str += fmt.Sprintf("\n%s", d.Err)
	}

	return str
}
