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

// estudianteRepoMock es un doble de prueba de storage.EstudianteRepository.
// Implementa solo los 5 métodos de la interfaz estrecha; si el service
// necesitara más, el compilador nos lo diría aquí.
type estudianteRepoMock struct {
	mock.Mock
}

func (m *estudianteRepoMock) ListarEstudiantes() []models.Estudiante {
	args := m.Called()
	return args.Get(0).([]models.Estudiante)
}

func (m *estudianteRepoMock) BuscarEstudiantePorID(id int) (models.Estudiante, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Estudiante), args.Bool(1)
}

func (m *estudianteRepoMock) CrearEstudiante(e models.Estudiante) models.Estudiante {
	args := m.Called(e)
	return args.Get(0).(models.Estudiante)
}

func (m *estudianteRepoMock) ActualizarEstudiante(id int, datos models.Estudiante) (models.Estudiante, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Estudiante), args.Bool(1)
}

func (m *estudianteRepoMock) BorrarEstudiante(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// Red de seguridad en tiempo de compilación: el mock DEBE cumplir el contrato.
var _ storage.EstudianteRepository = (*estudianteRepoMock)(nil)

// TestEstudianteService_Crear verifica las reglas de negocio de validarEstudiante
// de forma completamente aislada, sin base de datos ni red.
//
// Qué comprueba: que un dato inválido sea rechazado y NO llegue al repositorio.
// Qué se rompería: si se elimina la validación en service/estudiante.go, el
// mock recibiría una llamada no esperada y el test fallaría con "unexpected call".
func TestEstudianteService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Estudiante
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre:        "nombre vacio -> ErrNombreVacio y NO llega al repo",
			entrada:       models.Estudiante{Nombre: "   ", Email: "test@escuela.edu", Matricula: "EST010"},
			errEsperado:   service.ErrNombreVacio,
			debePersistir: false,
		},
		{
			nombre:        "email vacio -> ErrEmailVacio y NO llega al repo",
			entrada:       models.Estudiante{Nombre: "Carlos López", Email: "", Matricula: "EST011"},
			errEsperado:   service.ErrEmailVacio,
			debePersistir: false,
		},
		{
			nombre:        "datos validos -> sin error y se persiste",
			entrada:       models.Estudiante{Nombre: "María García", Email: "maria@escuela.edu", Matricula: "EST012", Activo: true},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := new(estudianteRepoMock)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 99
				repo.On("CrearEstudiante", c.entrada).Return(guardado)
			}
			svc := service.NuevoEstudianteService(repo)

			creado, err := svc.Crear(c.entrada)

			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearEstudiante")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 99, creado.ID, "el service debe devolver el estudiante que entregó el repo")
				repo.AssertCalled(t, "CrearEstudiante", c.entrada)
			}
		})
	}
}

// TestEstudianteService_Obtener_NoEncontrado comprueba que el service traduce
// el comma-ok=false del repositorio en ErrNoEncontrado.
func TestEstudianteService_Obtener_NoEncontrado(t *testing.T) {
	repo := new(estudianteRepoMock)
	repo.On("BuscarEstudiantePorID", 999).Return(models.Estudiante{}, false)
	svc := service.NuevoEstudianteService(repo)

	_, err := svc.Obtener(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
