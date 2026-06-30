package handlers

import "escuela-api/internal/service"

// Server agrupa los servicios de los que dependen los handlers.
type Server struct {
	Estudiantes *service.EstudianteService
	Auth        *service.AuthService
}

func NewServer(estudiantes *service.EstudianteService, auth *service.AuthService) *Server {
	return &Server{
		Estudiantes: estudiantes,
		Auth:        auth,
	}
}
