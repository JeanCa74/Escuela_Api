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

// construirEntornoInscripcion arma el router con rutas de Inscripciones y middleware Auth real.
func construirEntornoInscripcion(t *testing.T) (http.Handler, string) {
	t.Helper()

	almacen := storage.NuevaMemoria()
	usuarios := nuevoUsuarioRepoFake()

	inscripcionSvc := service.NuevoInscripcionService(almacen)
	authSvc := service.NuevoAuthService(usuarios)
	srv := handlers.NewServer(nil, nil, inscripcionSvc, authSvc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Get("/inscripciones", srv.ListarInscripciones)
			r.Post("/inscripciones", srv.CrearInscripcion)
			r.Get("/inscripciones/{id}", srv.ObtenerInscripcion)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// TestCrearInscripcion_Exitosa: POST con token y calificación válida -> 201 Created.
func TestCrearInscripcion_Exitosa(t *testing.T) {
	h, token := construirEntornoInscripcion(t)
	body := `{"estudiante_id":1,"curso_id":2,"calificacion":8.5}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/inscripciones", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var creada models.Inscripcion
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&creada))
	assert.NotZero(t, creada.ID)
	assert.Equal(t, 8.5, creada.Calificacion)
}

// TestCrearInscripcion_CalificacionInvalida: calificación > 10 -> 400 Bad Request.
func TestCrearInscripcion_CalificacionInvalida(t *testing.T) {
	h, token := construirEntornoInscripcion(t)
	body := `{"estudiante_id":1,"curso_id":2,"calificacion":15}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/inscripciones", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestRutaInscripciones_SinToken: sin Authorization el middleware devuelve 401.
//
// Qué se rompería: si se elimina r.Use(middleware.Auth(...)), la petición llega
// al handler y responde 201 en lugar de 401.
func TestRutaInscripciones_SinToken(t *testing.T) {
	h, _ := construirEntornoInscripcion(t)
	body := `{"estudiante_id":1,"curso_id":2,"calificacion":7}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/inscripciones", strings.NewReader(body))
	// A propósito: NO establecemos el header Authorization.
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
