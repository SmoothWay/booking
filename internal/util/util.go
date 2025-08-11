package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func ReadJSON(r *http.Request, v any) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if len(bodyBytes) == 0 {
		return nil
	}

	return json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(v)
}
