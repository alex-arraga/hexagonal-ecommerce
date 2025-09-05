package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseRequestBody[T any](r *http.Request) (T, error) {
	var params T

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	defer r.Body.Close()

	if err := decoder.Decode(&params); err != nil {
		return params, fmt.Errorf("error decoding parameters: %w", err)
	}
	return params, nil
}
