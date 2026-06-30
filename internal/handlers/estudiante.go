package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"escuela-api/internal/models"
)

// ListarEstudiantes atiende GET /api/v1/estudiantes.
func (s *Server) ListarEstudiantes(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Estudiantes.Listar())
}

// ObtenerEstudiante atiende GET /api/v1/estudiantes/{id}.
func (s *Server) ObtenerEstudiante(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	estudiante, err := s.Estudiantes.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, estudiante)
}

// CrearEstudiante atiende POST /api/v1/estudiantes.
func (s *Server) CrearEstudiante(w http.ResponseWriter, r *http.Request) {
	var nuevo models.Estudiante
	if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	creado, err := s.Estudiantes.Crear(nuevo)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

// ActualizarEstudiante atiende PUT /api/v1/estudiantes/{id}.
func (s *Server) ActualizarEstudiante(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	var datos models.Estudiante
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	actualizado, err := s.Estudiantes.Actualizar(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

// BorrarEstudiante atiende DELETE /api/v1/estudiantes/{id}.
func (s *Server) BorrarEstudiante(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	if err := s.Estudiantes.Borrar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
