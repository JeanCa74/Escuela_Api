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

// usuarioRepoFake: repositorio de usuarios en memoria para los tests de handler.
// Es un fake (guarda de verdad), no un mock (no verifica llamadas).
type usuarioRepoFake struct {
	porEmail map[string]models.Usuario
	nextID   int
}

func nuevoUsuarioRepoFake() *usuarioRepoFake {
	return &usuarioRepoFake{porEmail: map[string]models.Usuario{}, nextID: 1}
}

func (f *usuarioRepoFake) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.ID = f.nextID
	f.nextID++
	f.porEmail[u.Email] = u
	return u, nil
}

func (f *usuarioRepoFake) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	u, ok := f.porEmail[email]
	return u, ok
}

// construirEntorno arma el MISMO router que main.go (mismas rutas, mismo
// middleware.Auth real) pero con almacén en memoria y repo de usuarios fake.
// Devuelve el handler listo para httptest y un token válido ya emitido.
//
// Clave pedagógica: probamos a través del middleware REAL, no de uno simplificado.
// Si el wiring de la ruta protegida se rompe, este test lo detecta.
func construirEntorno(t *testing.T) (http.Handler, string) {
	t.Helper()

	almacen := storage.NuevaMemoria()
	almacen.SeedEstudiantes()
	usuarios := nuevoUsuarioRepoFake()

	estudianteSvc := service.NuevoEstudianteService(almacen)
	authSvc := service.NuevoAuthService(usuarios)
	srv := handlers.NewServer(estudianteSvc, authSvc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc)) // <- middleware real de autenticación
			r.Get("/estudiantes", srv.ListarEstudiantes)
			r.Post("/estudiantes", srv.CrearEstudiante)
			r.Get("/estudiantes/{id}", srv.ObtenerEstudiante)
			r.Put("/estudiantes/{id}", srv.ActualizarEstudiante)
			r.Delete("/estudiantes/{id}", srv.BorrarEstudiante)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// registrarYObtenerToken hace register + login contra el propio router para
// conseguir un JWT válido, igual que lo haría un cliente real.
func registrarYObtenerToken(t *testing.T, h http.Handler) string {
	t.Helper()
	cred := `{"email":"docente@escuela.edu","password":"clave123"}`

	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(cred))
	h.ServeHTTP(httptest.NewRecorder(), reqReg)

	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(cred))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, reqLogin)
	require.Equal(t, http.StatusOK, rec.Code, "el login deberia devolver 200")

	var resp struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.NotEmpty(t, resp.Token)
	return resp.Token
}

// TestCrearEstudiante_Exitoso: POST con token y cuerpo válido -> 201 Created.
//
// Qué comprueba: que el handler acepta un estudiante válido y responde 201.
// Qué se rompería: si el handler omite RespondJSON(201, ...) o el middleware
// rechaza el token por error en el wiring de la ruta.
func TestCrearEstudiante_Exitoso(t *testing.T) {
	h, token := construirEntorno(t)
	body := `{"nombre":"Pedro Almeida","email":"pedro@escuela.edu","matricula":"EST099","activo":true}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/estudiantes", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var creado models.Estudiante
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&creado))
	assert.NotZero(t, creado.ID)
	assert.Equal(t, "Pedro Almeida", creado.Nombre)
}

// TestCrearEstudiante_Invalido: nombre vacío viola la regla de negocio -> 400.
func TestCrearEstudiante_Invalido(t *testing.T) {
	h, token := construirEntorno(t)
	body := `{"nombre":"   ","email":"test@escuela.edu","matricula":"EST098"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/estudiantes", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestRutaProtegida_SinToken: sin header Authorization, el middleware rechaza
// antes de llegar al handler -> 401 Unauthorized.
//
// Qué comprueba: que la ruta /estudiantes está efectivamente protegida.
// Qué se rompería: si se eliminara r.Use(middleware.Auth(...)) del grupo protegido,
// la petición llegaría al handler y respondería 201 en lugar de 401.
func TestRutaProtegida_SinToken(t *testing.T) {
	h, _ := construirEntorno(t)
	body := `{"nombre":"Sin Token","email":"sin@escuela.edu","matricula":"EST000"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/estudiantes", strings.NewReader(body))
	// A propósito: NO establecemos el header Authorization.
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
