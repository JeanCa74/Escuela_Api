package storage

import "escuela-api/internal/models"

// EstudianteRepository es el contrato de persistencia de estudiantes.
type EstudianteRepository interface {
	ListarEstudiantes() []models.Estudiante
	BuscarEstudiantePorID(id int) (models.Estudiante, bool)
	CrearEstudiante(e models.Estudiante) models.Estudiante
	ActualizarEstudiante(id int, datos models.Estudiante) (models.Estudiante, bool)
	BorrarEstudiante(id int) bool
}

// CursoRepository es el contrato de persistencia de cursos.
type CursoRepository interface {
	ListarCursos() []models.Curso
	BuscarCursoPorID(id int) (models.Curso, bool)
	CrearCurso(c models.Curso) models.Curso
}

// InscripcionRepository es el contrato de persistencia de inscripciones.
type InscripcionRepository interface {
	ListarInscripciones() []models.Inscripcion
	BuscarInscripcionPorID(id int) (models.Inscripcion, bool)
	CrearInscripcion(i models.Inscripcion) models.Inscripcion
}

// Almacen combina los tres repositorios del dominio escolar.
type Almacen interface {
	EstudianteRepository
	CursoRepository
	InscripcionRepository
}

// UserRepository es el contrato de persistencia de usuarios (para auth).
type UserRepository interface {
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
}
