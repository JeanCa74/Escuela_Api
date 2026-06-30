package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"escuela-api/internal/models"
	"escuela-api/internal/storage"
)

// TestGORM_CrearYBuscarInscripcion verifica que el repositorio real (GORM + SQLite en memoria)
// refleja la creación al buscar por ID.
func TestGORM_CrearYBuscarInscripcion(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err, "debe abrir la base en memoria sin error")

	entrada := models.Inscripcion{
		EstudianteID: 1,
		CursoID:      2,
		Calificacion: 9.0,
	}

	creada := almacen.CrearInscripcion(entrada)

	require.NotZero(t, creada.ID, "GORM debe asignar un ID autoincremental")

	encontrada, ok := almacen.BuscarInscripcionPorID(creada.ID)
	require.True(t, ok, "la inscripción recién creada debe encontrarse por su ID")
	assert.Equal(t, 1, encontrada.EstudianteID)
	assert.Equal(t, 2, encontrada.CursoID)
	assert.Equal(t, 9.0, encontrada.Calificacion)
}

// TestGORM_ListarInscripcionesReflejaCreacion verifica que ListarInscripciones devuelve los registros insertados.
func TestGORM_ListarInscripcionesReflejaCreacion(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err)

	almacen.CrearInscripcion(models.Inscripcion{EstudianteID: 1, CursoID: 1, Calificacion: 7.5})
	almacen.CrearInscripcion(models.Inscripcion{EstudianteID: 2, CursoID: 1, Calificacion: 8.0})

	lista := almacen.ListarInscripciones()
	assert.Len(t, lista, 2, "deben listarse las dos inscripciones creadas")
}
