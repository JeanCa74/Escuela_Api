package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"escuela-api/internal/models"
)

func (s *Server) ListarInscripciones(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.Inscripciones.Listar())
}

func (s *Server) ObtenerInscripcion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	ins, err := s.Inscripciones.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, ins)
}

func (s *Server) CrearInscripcion(w http.ResponseWriter, r *http.Request) {
	var i models.Inscripcion
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	creada, err := s.Inscripciones.Crear(i)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}
