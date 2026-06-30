package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"escuela-api/internal/models"
	"escuela-api/internal/storage"
)

// TestGORM_CrearYBuscarCurso verifica que el repositorio real (GORM + SQLite en memoria)
// refleja la creación al buscar por ID.
func TestGORM_CrearYBuscarCurso(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err, "debe abrir la base en memoria sin error")

	entrada := models.Curso{
		Nombre:   "Cálculo Diferencial",
		Codigo:   "MAT201",
		Creditos: 5,
	}

	creado := almacen.CrearCurso(entrada)

	require.NotZero(t, creado.ID, "GORM debe asignar un ID autoincremental")

	encontrado, ok := almacen.BuscarCursoPorID(creado.ID)
	require.True(t, ok, "el curso recién creado debe encontrarse por su ID")
	assert.Equal(t, "Cálculo Diferencial", encontrado.Nombre)
	assert.Equal(t, "MAT201", encontrado.Codigo)
	assert.Equal(t, 5, encontrado.Creditos)
}

// TestGORM_ListarCursosReflejaCreacion verifica que ListarCursos devuelve los registros insertados.
func TestGORM_ListarCursosReflejaCreacion(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err)

	almacen.CrearCurso(models.Curso{Nombre: "Química", Codigo: "QUI101", Creditos: 3})
	almacen.CrearCurso(models.Curso{Nombre: "Biología", Codigo: "BIO101", Creditos: 4})

	lista := almacen.ListarCursos()
	assert.Len(t, lista, 2, "deben listarse los dos cursos creados")
}
