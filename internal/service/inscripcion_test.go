package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"escuela-api/internal/models"
	"escuela-api/internal/service"
	"escuela-api/internal/storage"
)

// inscripcionRepoMock es un doble de prueba que registra llamadas a InscripcionRepository.
type inscripcionRepoMock struct {
	mock.Mock
}

func (m *inscripcionRepoMock) ListarInscripciones() []models.Inscripcion {
	args := m.Called()
	return args.Get(0).([]models.Inscripcion)
}

func (m *inscripcionRepoMock) BuscarInscripcionPorID(id int) (models.Inscripcion, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Inscripcion), args.Bool(1)
}

func (m *inscripcionRepoMock) CrearInscripcion(i models.Inscripcion) models.Inscripcion {
	args := m.Called(i)
	return args.Get(0).(models.Inscripcion)
}

// Garantía en tiempo de compilación: inscripcionRepoMock cumple la interfaz.
var _ storage.InscripcionRepository = (*inscripcionRepoMock)(nil)

// TestInscripcionService_Crear verifica que validarInscripcion rechaza calificaciones
// fuera del rango [0,10] sin llamar al repositorio.
func TestInscripcionService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Inscripcion
		errEsperado   error
		debePersistir bool
	}{
		{
			"calificacion > 10 -> ErrCalificacionInvalida y NO llega al repo",
			models.Inscripcion{EstudianteID: 1, CursoID: 2, Calificacion: 15},
			service.ErrCalificacionInvalida,
			false,
		},
		{
			"calificacion negativa -> ErrCalificacionInvalida y NO llega al repo",
			models.Inscripcion{EstudianteID: 1, CursoID: 2, Calificacion: -1},
			service.ErrCalificacionInvalida,
			false,
		},
		{
			"calificacion valida -> sin error y se persiste en el repo",
			models.Inscripcion{EstudianteID: 1, CursoID: 2, Calificacion: 8.5},
			nil,
			true,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := new(inscripcionRepoMock)

			if tc.debePersistir {
				esperado := tc.entrada
				esperado.ID = 55
				repo.On("CrearInscripcion", tc.entrada).Return(esperado)
			}

			svc := service.NuevoInscripcionService(repo)
			creada, err := svc.Crear(tc.entrada)

			if tc.errEsperado != nil {
				require.ErrorIs(t, err, tc.errEsperado)
				// La clave del test: el repositorio NUNCA debe ser llamado.
				repo.AssertNotCalled(t, "CrearInscripcion")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 55, creada.ID)
				repo.AssertCalled(t, "CrearInscripcion", tc.entrada)
			}

			repo.AssertExpectations(t)
		})
	}
}
