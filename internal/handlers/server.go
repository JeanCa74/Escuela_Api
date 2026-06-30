package handlers

import "escuela-api/internal/service"

// Server agrupa los servicios de los que dependen los handlers.
type Server struct {
	Estudiantes   *service.EstudianteService
	Cursos        *service.CursoService
	Inscripciones *service.InscripcionService
	Auth          *service.AuthService
}

func NewServer(
	estudiantes *service.EstudianteService,
	cursos *service.CursoService,
	inscripciones *service.InscripcionService,
	auth *service.AuthService,
) *Server {
	return &Server{
		Estudiantes:   estudiantes,
		Cursos:        cursos,
		Inscripciones: inscripciones,
		Auth:          auth,
	}
}
