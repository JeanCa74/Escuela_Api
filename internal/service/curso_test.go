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

// cursoRepoMock es un doble de prueba que registra llamadas a CursoRepository.
// Permite verificar que una operación NO llamó al repositorio (AssertNotCalled).
type cursoRepoMock struct {
	mock.Mock
}

func (m *cursoRepoMock) ListarCursos() []models.Curso {
	args := m.Called()
	return args.Get(0).([]models.Curso)
}

func (m *cursoRepoMock) BuscarCursoPorID(id int) (models.Curso, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Curso), args.Bool(1)
}

func (m *cursoRepoMock) CrearCurso(c models.Curso) models.Curso {
	args := m.Called(c)
	return args.Get(0).(models.Curso)
}

// Garantía en tiempo de compilación: cursoRepoMock cumple la interfaz.
var _ storage.CursoRepository = (*cursoRepoMock)(nil)

// TestCursoService_Crear verifica que validarCurso rechaza datos inválidos
// sin llegar al repositorio, y que datos válidos sí persisten.
func TestCursoService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Curso
		errEsperado   error
		debePersistir bool
	}{
		{
			"creditos cero -> ErrCreditosInvalidos y NO llega al repo",
			models.Curso{Nombre: "Matemáticas", Codigo: "MAT101", Creditos: 0},
			service.ErrCreditosInvalidos,
			false,
		},
		{
			"creditos negativos -> ErrCreditosInvalidos y NO llega al repo",
			models.Curso{Nombre: "Física", Codigo: "FIS101", Creditos: -3},
			service.ErrCreditosInvalidos,
			false,
		},
		{
			"nombre vacio -> ErrNombreVacio y NO llega al repo",
			models.Curso{Nombre: "   ", Codigo: "PRG101", Creditos: 4},
			service.ErrNombreVacio,
			false,
		},
		{
			"datos validos -> sin error y se persiste en el repo",
			models.Curso{Nombre: "Programación", Codigo: "PRG201", Creditos: 4},
			nil,
			true,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := new(cursoRepoMock)

			if tc.debePersistir {
				// Configuramos lo que el mock devuelve cuando se llame CrearCurso.
				esperado := tc.entrada
				esperado.ID = 77
				repo.On("CrearCurso", tc.entrada).Return(esperado)
			}

			svc := service.NuevoCursoService(repo)
			creado, err := svc.Crear(tc.entrada)

			if tc.errEsperado != nil {
				require.ErrorIs(t, err, tc.errEsperado)
				// La clave del test: el repositorio NUNCA debe ser llamado.
				repo.AssertNotCalled(t, "CrearCurso")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 77, creado.ID)
				repo.AssertCalled(t, "CrearCurso", tc.entrada)
			}

			repo.AssertExpectations(t)
		})
	}
}
