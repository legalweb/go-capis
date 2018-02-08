package capis

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"lwebco.de/comparisonapis/pkg/utils/uuid"
)

func unmarshalResponse(res *http.Response, obj interface{}) error {
	rb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(rb, obj)
}

func NewUUIDString() string {
	return uuid.New().String()
}
