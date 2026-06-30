package service

import (
	"escuela-api/internal/models"
	"escuela-api/internal/storage"
)

// InscripcionService contiene la lógica de negocio de inscripciones.
type InscripcionService struct {
	repo storage.InscripcionRepository
}

func NuevoInscripcionService(repo storage.InscripcionRepository) *InscripcionService {
	return &InscripcionService{repo: repo}
}

func (s *InscripcionService) Listar() []models.Inscripcion {
	return s.repo.ListarInscripciones()
}

func (s *InscripcionService) Obtener(id int) (models.Inscripcion, error) {
	i, ok := s.repo.BuscarInscripcionPorID(id)
	if !ok {
		return models.Inscripcion{}, ErrNoEncontrado
	}
	return i, nil
}

func (s *InscripcionService) Crear(i models.Inscripcion) (models.Inscripcion, error) {
	if err := validarInscripcion(i); err != nil {
		return models.Inscripcion{}, err
	}
	return s.repo.CrearInscripcion(i), nil
}

// validarInscripcion aplica las reglas de negocio del dominio.
// Regla: la calificación debe estar en el rango [0, 10].
func validarInscripcion(i models.Inscripcion) error {
	if i.Calificacion < 0 || i.Calificacion > 10 {
		return ErrCalificacionInvalida
	}
	return nil
}
