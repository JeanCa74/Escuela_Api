package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"escuela-api/internal/models"
)

func (s *Server) ListarCursos(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.Cursos.Listar())
}

func (s *Server) ObtenerCurso(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	curso, err := s.Cursos.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, curso)
}

func (s *Server) CrearCurso(w http.ResponseWriter, r *http.Request) {
	var c models.Curso
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	creado, err := s.Cursos.Crear(c)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}
