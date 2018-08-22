package capis

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/satori/go.uuid"
)

func unmarshalResponse(res *http.Response, obj interface{}) error {
	rb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(rb, obj)
}

// NewUUIDString is a utility function to get a UUID without worry about
// errors being returned from the commonly use library for generating
// UUIDs in Golang.
func NewUUIDString() string {
	uu, _ := uuid.NewV4()
	return uu.String()
}
