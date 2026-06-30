// Command escuela-api arranca el servidor HTTP de la Escuela.
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"escuela-api/internal/handlers"
	"escuela-api/internal/middleware"
	"escuela-api/internal/service"
	"escuela-api/internal/storage"
)

func main() {
	// 1. GORM abre la base, ejecuta AutoMigrate y devuelve el almacén.
	almacen, err := storage.NuevoAlmacenGORM("escuela.db")
	if err != nil {
		log.Fatal("no se pudo abrir la base de datos: ", err)
	}

	// 2. Capa de servicio con inyección de dependencias.
	estudianteSvc := service.NuevoEstudianteService(almacen)
	cursoSvc := service.NuevoCursoService(almacen)
	inscripcionSvc := service.NuevoInscripcionService(almacen)
	authSvc := service.NuevoAuthService(almacen)

	// 3. Server con los servicios inyectados.
	servidor := handlers.NewServer(estudianteSvc, cursoSvc, inscripcionSvc, authSvc)

	// 4. Router + middleware global.
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 5. Rutas versionadas /api/v1/.
	r.Route("/api/v1", func(r chi.Router) {
		// Públicas: registro y login.
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Protegidas: exigen JWT válido en Authorization: Bearer <token>.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))

			// Módulo Jean Carlos — Estudiantes
			r.Get("/estudiantes", servidor.ListarEstudiantes)
			r.Post("/estudiantes", servidor.CrearEstudiante)
			r.Get("/estudiantes/{id}", servidor.ObtenerEstudiante)
			r.Put("/estudiantes/{id}", servidor.ActualizarEstudiante)
			r.Delete("/estudiantes/{id}", servidor.BorrarEstudiante)

			// Módulo Jhon — Cursos
			r.Get("/cursos", servidor.ListarCursos)
			r.Post("/cursos", servidor.CrearCurso)
			r.Get("/cursos/{id}", servidor.ObtenerCurso)

			// Módulo Maria José — Inscripciones
			r.Get("/inscripciones", servidor.ListarInscripciones)
			r.Post("/inscripciones", servidor.CrearInscripcion)
			r.Get("/inscripciones/{id}", servidor.ObtenerInscripcion)
		})
	})

	log.Println("Servidor escuela-api escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
