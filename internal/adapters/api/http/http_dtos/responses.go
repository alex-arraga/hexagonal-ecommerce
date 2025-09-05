package httpdtos

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, code int, msg string, payload interface{}) {
	if payload == nil {
		payload = struct{}{}
	}

	// Response struct
	response := struct {
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		Msg:  msg,
		Data: payload,
	}

	// Convert to JSON
	data, err := json.Marshal(response)
	if err != nil {
		slog.Error("Failed to marshal JSON response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Config headers and write the response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		slog.Error("Error writing JSON response", "error", err)
		return
	}
}

func RespondError(w http.ResponseWriter, code int, msg string) {
	slog.Error("Responding with error", "code", code, "msg", msg)

	// Estructura de error con "error" en lugar de "msg"
	response := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	// Convertir a JSON
	data, err := json.Marshal(response)
	if err != nil {
		slog.Error("Failed to marshal JSON error response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Configurar headers y escribir la respuesta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		slog.Error("Error writing JSON response", "error", err)
	}
}
