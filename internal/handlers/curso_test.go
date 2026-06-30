package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"escuela-api/internal/handlers"
	"escuela-api/internal/middleware"
	"escuela-api/internal/models"
	"escuela-api/internal/service"
	"escuela-api/internal/storage"
)

// construirEntornoCurso arma el router con las rutas de Cursos y middleware Auth real.
// Reutiliza nuevoUsuarioRepoFake y registrarYObtenerToken definidos en estudiante_test.go.
func construirEntornoCurso(t *testing.T) (http.Handler, string) {
	t.Helper()

	almacen := storage.NuevaMemoria()
	usuarios := nuevoUsuarioRepoFake()

	cursoSvc := service.NuevoCursoService(almacen)
	authSvc := service.NuevoAuthService(usuarios)
	srv := handlers.NewServer(nil, cursoSvc, nil, authSvc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Get("/cursos", srv.ListarCursos)
			r.Post("/cursos", srv.CrearCurso)
			r.Get("/cursos/{id}", srv.ObtenerCurso)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// TestCrearCurso_Exitoso: POST con token y créditos válidos -> 201 Created.
func TestCrearCurso_Exitoso(t *testing.T) {
	h, token := construirEntornoCurso(t)
	body := `{"nombre":"Algebra Lineal","codigo":"ALG101","creditos":4}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cursos", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var creado models.Curso
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&creado))
	assert.NotZero(t, creado.ID)
	assert.Equal(t, "Algebra Lineal", creado.Nombre)
}

// TestCrearCurso_CreditosInvalidos: créditos = 0 viola la regla de negocio -> 400.
func TestCrearCurso_CreditosInvalidos(t *testing.T) {
	h, token := construirEntornoCurso(t)
	body := `{"nombre":"Curso Sin Creditos","codigo":"SIN001","creditos":0}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cursos", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestRutaCursos_SinToken: sin header Authorization el middleware devuelve 401.
//
// Qué se rompería: si se elimina r.Use(middleware.Auth(...)), la petición llegaría
// al handler y respondería 201 en lugar de 401.
func TestRutaCursos_SinToken(t *testing.T) {
	h, _ := construirEntornoCurso(t)
	body := `{"nombre":"Redes","codigo":"RED101","creditos":3}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cursos", strings.NewReader(body))
	// A propósito: NO establecemos el header Authorization.
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
