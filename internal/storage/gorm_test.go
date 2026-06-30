package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"escuela-api/internal/models"
	"escuela-api/internal/storage"
)

// TestGORM_CrearYBuscarEstudiante verifica que el repositorio real (GORM + SQLite
// en memoria) refleja una creacion al buscar por ID.
// Patron: Crear → BuscarPorID devuelve el mismo registro.
func TestGORM_CrearYBuscarEstudiante(t *testing.T) {
	// ":memory:" abre una base SQLite que vive solo mientras dure el test.
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err, "debe abrir la base en memoria sin error")

	entrada := models.Estudiante{
		Nombre:    "Ana Torres",
		Email:     "ana@escuela.edu",
		Matricula: "EST001",
		Activo:    true,
	}

	creado := almacen.CrearEstudiante(entrada)

	// GORM debe haber asignado un ID mayor que cero.
	require.NotZero(t, creado.ID, "GORM debe asignar un ID autoincremental")

	encontrado, ok := almacen.BuscarEstudiantePorID(creado.ID)
	require.True(t, ok, "el estudiante recien creado debe encontrarse por su ID")
	assert.Equal(t, "Ana Torres", encontrado.Nombre)
	assert.Equal(t, "ana@escuela.edu", encontrado.Email)
	assert.Equal(t, "EST001", encontrado.Matricula)
}

// TestGORM_ListarReflejaCreacion verifica que ListarEstudiantes devuelve los
// registros insertados con CrearEstudiante.
func TestGORM_ListarReflejaCreacion(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err)

	almacen.CrearEstudiante(models.Estudiante{Nombre: "Mario Ruiz", Email: "mario@escuela.edu", Matricula: "EST002", Activo: true})
	almacen.CrearEstudiante(models.Estudiante{Nombre: "Sofia Vega", Email: "sofia@escuela.edu", Matricula: "EST003", Activo: true})

	lista := almacen.ListarEstudiantes()
	assert.Len(t, lista, 2, "deben listarse los dos estudiantes creados")
}
