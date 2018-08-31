package dumpcurl

import (
	"net/http"

	"github.com/moul/http2curl"

	capis "lwebco.de/go-capis"
)

// New dump curl request middleware.
func New(out func(string)) capis.RequestMiddlewareFunc {
	return func(r *http.Request) *http.Request {
		curl, err := http2curl.GetCurlCommand(r)

		if err != nil {
			out("unable to create CURL command from request: " + err.Error())
		} else {
			out(curl.String())
		}

		return r
	}
}
