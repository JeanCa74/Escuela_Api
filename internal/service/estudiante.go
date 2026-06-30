package service

import (
	"strings"

	"escuela-api/internal/models"
	"escuela-api/internal/storage"
)

// EstudianteService contiene la lógica de negocio de estudiantes.
// Depende SOLO de storage.EstudianteRepository (interfaz estrecha).
type EstudianteService struct {
	repo storage.EstudianteRepository
}

func NuevoEstudianteService(repo storage.EstudianteRepository) *EstudianteService {
	return &EstudianteService{repo: repo}
}

func (s *EstudianteService) Listar() []models.Estudiante {
	return s.repo.ListarEstudiantes()
}

func (s *EstudianteService) Obtener(id int) (models.Estudiante, error) {
	e, ok := s.repo.BuscarEstudiantePorID(id)
	if !ok {
		return models.Estudiante{}, ErrNoEncontrado
	}
	return e, nil
}

func (s *EstudianteService) Crear(e models.Estudiante) (models.Estudiante, error) {
	if err := validarEstudiante(e); err != nil {
		return models.Estudiante{}, err
	}
	return s.repo.CrearEstudiante(e), nil
}

func (s *EstudianteService) Actualizar(id int, datos models.Estudiante) (models.Estudiante, error) {
	if err := validarEstudiante(datos); err != nil {
		return models.Estudiante{}, err
	}
	actualizado, ok := s.repo.ActualizarEstudiante(id, datos)
	if !ok {
		return models.Estudiante{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *EstudianteService) Borrar(id int) error {
	if !s.repo.BorrarEstudiante(id) {
		return ErrNoEncontrado
	}
	return nil
}

// validarEstudiante centraliza las reglas de negocio de la entidad.
func validarEstudiante(e models.Estudiante) error {
	if strings.TrimSpace(e.Nombre) == "" {
		return ErrNombreVacio
	}
	if strings.TrimSpace(e.Email) == "" {
		return ErrEmailVacio
	}
	return nil
}
