package storage

import (
	"sync"

	"escuela-api/internal/models"
)

// Memoria es un almacén unificado en RAM; se usa como fake en los tests de handler.
type Memoria struct {
	estudiantes      []models.Estudiante
	nextEstudianteID int

	cursos      []models.Curso
	nextCursoID int

	inscripciones      []models.Inscripcion
	nextInscripcionID  int

	mu sync.Mutex
}

// NuevaMemoria crea un almacén vacío listo para usar.
func NuevaMemoria() *Memoria {
	return &Memoria{
		estudiantes:       []models.Estudiante{},
		nextEstudianteID:  1,
		cursos:            []models.Curso{},
		nextCursoID:       1,
		inscripciones:     []models.Inscripcion{},
		nextInscripcionID: 1,
	}
}

// SeedEstudiantes carga datos iniciales de prueba.
func (m *Memoria) SeedEstudiantes() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.estudiantes = []models.Estudiante{
		{ID: 1, Nombre: "Ana Torres", Email: "ana@escuela.edu", Matricula: "EST001", Activo: true},
		{ID: 2, Nombre: "Luis Pérez", Email: "luis@escuela.edu", Matricula: "EST002", Activo: true},
	}
	m.nextEstudianteID = 3
}

// =========================================================
// ESTUDIANTES
// =========================================================

func (m *Memoria) ListarEstudiantes() []models.Estudiante {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.Estudiante, len(m.estudiantes))
	copy(copia, m.estudiantes)
	return copia
}

func (m *Memoria) BuscarEstudiantePorID(id int) (models.Estudiante, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.estudiantes {
		if e.ID == id {
			return e, true
		}
	}
	return models.Estudiante{}, false
}

func (m *Memoria) CrearEstudiante(e models.Estudiante) models.Estudiante {
	m.mu.Lock()
	defer m.mu.Unlock()
	e.ID = m.nextEstudianteID
	m.nextEstudianteID++
	m.estudiantes = append(m.estudiantes, e)
	return e
}

func (m *Memoria) ActualizarEstudiante(id int, datos models.Estudiante) (models.Estudiante, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, e := range m.estudiantes {
		if e.ID == id {
			datos.ID = id
			m.estudiantes[i] = datos
			return datos, true
		}
	}
	return models.Estudiante{}, false
}

func (m *Memoria) BorrarEstudiante(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, e := range m.estudiantes {
		if e.ID == id {
			m.estudiantes = append(m.estudiantes[:i], m.estudiantes[i+1:]...)
			return true
		}
	}
	return false
}

// =========================================================
// CURSOS
// =========================================================

func (m *Memoria) ListarCursos() []models.Curso {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.Curso, len(m.cursos))
	copy(copia, m.cursos)
	return copia
}

func (m *Memoria) BuscarCursoPorID(id int) (models.Curso, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.cursos {
		if c.ID == id {
			return c, true
		}
	}
	return models.Curso{}, false
}

func (m *Memoria) CrearCurso(c models.Curso) models.Curso {
	m.mu.Lock()
	defer m.mu.Unlock()
	c.ID = m.nextCursoID
	m.nextCursoID++
	m.cursos = append(m.cursos, c)
	return c
}

// =========================================================
// INSCRIPCIONES
// =========================================================

func (m *Memoria) ListarInscripciones() []models.Inscripcion {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.Inscripcion, len(m.inscripciones))
	copy(copia, m.inscripciones)
	return copia
}

func (m *Memoria) BuscarInscripcionPorID(id int) (models.Inscripcion, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, i := range m.inscripciones {
		if i.ID == id {
			return i, true
		}
	}
	return models.Inscripcion{}, false
}

func (m *Memoria) CrearInscripcion(i models.Inscripcion) models.Inscripcion {
	m.mu.Lock()
	defer m.mu.Unlock()
	i.ID = m.nextInscripcionID
	m.nextInscripcionID++
	m.inscripciones = append(m.inscripciones, i)
	return i
}

// Chequeo en tiempo de compilación: Memoria debe cumplir Almacen.
var _ Almacen = (*Memoria)(nil)
