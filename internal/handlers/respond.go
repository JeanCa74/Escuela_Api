package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"escuela-api/internal/service"
)

// RespondJSON escribe data como JSON con el status HTTP indicado.
func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error codificando JSON: %v", err)
	}
}

// RespondError escribe un error en formato JSON: {"error": "..."}.
func RespondError(w http.ResponseWriter, status int, mensaje string) {
	RespondJSON(w, status, map[string]string{"error": mensaje})
}

// statusDeError traduce errores de dominio a código HTTP.
func statusDeError(err error) int {
	switch {
	case errors.Is(err, service.ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, service.ErrEmailEnUso):
		return http.StatusConflict
	case errors.Is(err, service.ErrCredencialesInvalidas):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrNombreVacio),
		errors.Is(err, service.ErrEmailVacio),
		errors.Is(err, service.ErrCreditosInvalidos),
		errors.Is(err, service.ErrCalificacionInvalida):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
