package service

import (
	"strings"

	"escuela-api/internal/models"
	"escuela-api/internal/storage"
)

// CursoService contiene la lógica de negocio de cursos.
type CursoService struct {
	repo storage.CursoRepository
}

func NuevoCursoService(repo storage.CursoRepository) *CursoService {
	return &CursoService{repo: repo}
}

func (s *CursoService) Listar() []models.Curso {
	return s.repo.ListarCursos()
}

func (s *CursoService) Obtener(id int) (models.Curso, error) {
	c, ok := s.repo.BuscarCursoPorID(id)
	if !ok {
		return models.Curso{}, ErrNoEncontrado
	}
	return c, nil
}

func (s *CursoService) Crear(c models.Curso) (models.Curso, error) {
	if err := validarCurso(c); err != nil {
		return models.Curso{}, err
	}
	return s.repo.CrearCurso(c), nil
}

// validarCurso aplica las reglas de negocio del dominio.
// Regla 1: el nombre no puede estar vacío.
// Regla 2: los créditos deben ser mayores a cero.
func validarCurso(c models.Curso) error {
	if strings.TrimSpace(c.Nombre) == "" {
		return ErrNombreVacio
	}
	if c.Creditos <= 0 {
		return ErrCreditosInvalidos
	}
	return nil
}
