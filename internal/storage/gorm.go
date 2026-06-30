package storage

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"escuela-api/internal/models"
)

// AlmacenGORM implementa Almacen usando GORM sobre SQLite.
// En tests se abre con DSN ":memory:" para no tocar el disco.
type AlmacenGORM struct {
	db *gorm.DB
}

// NuevoAlmacenGORM abre la base de datos, ejecuta AutoMigrate y devuelve el almacén.
func NuevoAlmacenGORM(dsn string) (*AlmacenGORM, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Estudiante{}, &models.Curso{}, &models.Usuario{}); err != nil {
		return nil, err
	}
	return &AlmacenGORM{db: db}, nil
}

// =========================================================
// ESTUDIANTES
// =========================================================

func (a *AlmacenGORM) ListarEstudiantes() []models.Estudiante {
	var estudiantes []models.Estudiante
	a.db.Find(&estudiantes)
	return estudiantes
}

func (a *AlmacenGORM) BuscarEstudiantePorID(id int) (models.Estudiante, bool) {
	var e models.Estudiante
	if err := a.db.First(&e, id).Error; err != nil {
		return models.Estudiante{}, false
	}
	return e, true
}

func (a *AlmacenGORM) CrearEstudiante(e models.Estudiante) models.Estudiante {
	a.db.Create(&e)
	return e
}

func (a *AlmacenGORM) ActualizarEstudiante(id int, datos models.Estudiante) (models.Estudiante, bool) {
	var existente models.Estudiante
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Estudiante{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenGORM) BorrarEstudiante(id int) bool {
	res := a.db.Delete(&models.Estudiante{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// CURSOS
// =========================================================

func (a *AlmacenGORM) ListarCursos() []models.Curso {
	var cursos []models.Curso
	a.db.Find(&cursos)
	return cursos
}

func (a *AlmacenGORM) BuscarCursoPorID(id int) (models.Curso, bool) {
	var c models.Curso
	if err := a.db.First(&c, id).Error; err != nil {
		return models.Curso{}, false
	}
	return c, true
}

func (a *AlmacenGORM) CrearCurso(c models.Curso) models.Curso {
	a.db.Create(&c)
	return c
}

// =========================================================
// USUARIOS (para AuthService)
// =========================================================

func (a *AlmacenGORM) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	if err := a.db.Create(&u).Error; err != nil {
		return models.Usuario{}, err
	}
	return u, nil
}

func (a *AlmacenGORM) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario
	if err := a.db.Where("email = ?", email).First(&u).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

// Chequeos en tiempo de compilación.
var _ Almacen = (*AlmacenGORM)(nil)
var _ UserRepository = (*AlmacenGORM)(nil)
